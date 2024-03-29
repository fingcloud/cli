package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/avast/retry-go"
	"github.com/dustin/go-humanize"
	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/cmd/app"
	"github.com/fingcloud/cli/pkg/config"
	"github.com/fingcloud/cli/pkg/config/session"
	"github.com/fingcloud/cli/pkg/fileutils"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/fingcloud/cli/pkg/util"
	"github.com/gosuri/uilive"
	"github.com/r6m/spinner"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"go.uber.org/atomic"
)

type DeployOptions struct {
	client *api.Client

	config *api.AppConfig
	logs   *app.LogsOptions

	Path     string
	Quite    bool
	Dispatch bool
}

func NewCmdDeploy(ctx *cli.Context) *cobra.Command {
	opts := newOptions()

	cmd := &cobra.Command{
		Use:     "deploy",
		Short:   "deploy your application",
		Aliases: []string{"up"},
		Run: func(cmd *cobra.Command, args []string) {

			util.CheckErr(RunDeploy(ctx, opts))
		},
	}

	cfg, err := config.ReadAppConfig(opts.Path)
	if err == nil {
		opts.config = cfg
	}

	cmd.Flags().StringVarP(&opts.config.App, "app", "a", opts.config.App, "app name")
	cmd.Flags().StringVar(&opts.config.Platform, "platform", opts.config.Platform, "your app platform")
	cmd.Flags().StringVar(&opts.Path, "path", ".", "app path")
	cmd.Flags().BoolVarP(&opts.Quite, "quite", "q", opts.Quite, "quite output")
	cmd.Flags().BoolVarP(&opts.Dispatch, "dispatch", "d", opts.Quite, "dispatch logs")

	return cmd
}

var (
	s                = spinner.New().WithOptions(spinner.WithExitOnAbort(false), spinner.WithNotifySignals(false))
	ErrBuildCanceled = errors.New("Build canceled")
)

func newOptions() *DeployOptions {
	opts := new(DeployOptions)
	opts.config = &api.AppConfig{}
	return opts
}

func (opts *DeployOptions) validate() error {
	if opts.config.App == "" {
		return fmt.Errorf("app can't be empty")
	}

	return nil
}

// RunDeploy starts deploy process
func RunDeploy(ctx *cli.Context, opts *DeployOptions) error {
	sess, _ := session.CurrentSession()
	if sess != nil {
		fmt.Printf("Using session: %s\n", sess.Email)
	}

	if opts.config.App == "" {
		apps, err := ctx.Client.AppsList(&api.ListAppsOptions{})
		util.CheckErr(err)

		appOptions := funk.Map(apps, func(app *api.App) string { return app.Name }).([]string)
		if len(appOptions) == 0 {
			helpCreateApp()
			return fmt.Errorf("empty apps")
		}

		err = ui.PromptSelect("Choose your app", appOptions, &opts.config.App)
		util.CheckErr(err)
	}

	opts.logs = &app.LogsOptions{
		App:    opts.config.App,
		Since:  time.Second,
		Follow: true,
	}

	err := opts.validate()
	util.CheckErr(err)

	opts.printAppInfo()
	s.Start("Getting files...")

	files, err := fileutils.GetFiles(opts.Path)
	util.CheckErr(err)

	s.Success("Getting files OK")
	s.Start("Creating Deployment...")

	var deployment *api.Deployment
	deployFn := func() error {
		d, changes, err := ctx.Client.DeployemntCreate(opts.config.App, &api.CreateDeploymentOptions{
			Files:  files,
			Config: opts.config,
		})
		if err != nil {
			return retry.Unrecoverable(err)
		}

		if changes != nil {
			err = uploadChanges(ctx, opts.Path, opts.config.App, changes)
			if err != nil {
				return retry.Unrecoverable(err)
			}
			return errors.New("retry deployment")
		}

		deployment = d
		return nil
	}
	err = retry.Do(deployFn, retry.Attempts(3))
	util.CheckErr(err)

	s.Success("Creating Deployment OK")
	if deployment.Platform != "" {
		fmt.Printf("platform: %s\n", ui.Green(deployment.Platform))
	}

	err = readBuildLogs(ctx, opts.config.App, deployment.ID)
	util.CheckErr(err)

	if !opts.Dispatch {
		return app.RunLogs(ctx, opts.logs)
	}

	return nil
}

func uploadChanges(ctx *cli.Context, projectPath, app string, files []*api.FileInfo) error {
	s.UpdateMessage("Getting changed files...")
	tarBuf := new(bytes.Buffer)
	err := fileutils.Compress(projectPath, files, tarBuf)
	if err != nil {
		s.Error()
		return err
	}

	s.Successf("%d Files changed [%s]", len(files), humanize.Bytes(uint64(tarBuf.Len())))

	bar := ui.NewProgress(tarBuf.Len(), "Uploading")

	reporter := &api.ProgressReader{
		SetMax: func(max int64) { bar.ChangeMax64(max) },
		Add:    func(n int64) { bar.Add64(n) },
	}

	return ctx.Client.AppsUploadFiles(app, tarBuf, reporter)
}

func readBuildLogs(ctx *cli.Context, app string, deploymentID int64) error {
	s.Start("Building...")

	interruptCh := make(chan os.Signal, 1)
	stopCh := make(chan bool, 1)

	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	canceled := atomic.NewBool(false)

	bars := make([]*pullInfo, 0)
	writer := uilive.New()
	writer.Start()

	go func() {
		defer func() {
			writer.Stop()
			signal.Stop(interruptCh)
			close(interruptCh)
			close(stopCh)
		}()

		for {
			select {
			case <-stopCh:
				return
			case <-interruptCh:
				canceled.Store(true)
				ctx.Client.DeploymentCancel(app, deploymentID)
				return
			}
		}
	}()

	var from int64
	var starting bool
	for {
		buildLogs, err := ctx.Client.DeploymentBuildLogs(app, deploymentID, &api.LogsOptions{From: from})
		if err != nil {
			return err
		}

		for _, log := range buildLogs.Logs {

			if log.Message == "" {
				continue
			}

			// image pull
			if strings.HasPrefix(log.Message, "{'status': ") {
				log.Message = strings.ReplaceAll(log.Message, "'", "\"")

				pb := new(pullInfo)
				if err := json.Unmarshal([]byte(log.Message), pb); err != nil {
					continue
				}

				if pb.ID == "" {
					continue
				}

				var found bool
				for _, b := range bars {
					if b.ID == pb.ID {
						found = true
						b.Status = pb.Status
						b.Progress = pb.Progress
						break
					}
				}

				if !found {
					bars = append(bars, pb)
				}

				for _, b := range bars {
					fmt.Fprintln(writer, b)
				}
				writer.Flush()

				continue
			}

			s.ClearCurrentLine()
			fmt.Println(ui.Gray(log.Message))

			if strings.HasPrefix(log.Message, "Successfully tagged") {
				s.Success("Building OK")
				if !canceled.Load() {
					stopCh <- true
				}
				s.Start("Pushing...")
			}
		}

		if canceled.Load() {
			return ErrBuildCanceled
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusFailed {
			s.Error("Build failed")
			return errors.New("Build failed")
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusCancel {
			s.Error("Build canceled")
			return ErrBuildCanceled
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusStarting && !starting {
			s.Success("Pushing OK")
			s.Start("Starting...")
			starting = true
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusFinished {
			s.Success("Starting OK")
			helpSuccessulDeploy(buildLogs.Deployment.URL)
			return nil
		}

		if len(buildLogs.Logs) > 0 {
			from = buildLogs.Logs[len(buildLogs.Logs)-1].ID
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func (opts *DeployOptions) printAppInfo() {
	fmt.Printf("app: %s\n", ui.Green(opts.config.App))
}

func helpCreateApp() {
	fmt.Println("you don't have any apps on fing")
	fmt.Println("go to fing dashboard and create one:")
	fmt.Printf("\t%s\n\n", ui.Green("https://cloud.fing.ir/apps"))
}

func helpSuccessulDeploy(url string) {
	fmt.Println(ui.Info("App started successfully :)"))
	if url != "" {
		fmt.Println()
		fmt.Printf("\topen the following url in your browser:\n")
		fmt.Printf("\t%s", ui.Green(url))
		fmt.Println()
	}
}

type pullInfo struct {
	ID       string `json:"id"`
	Status   string `json:"status"`
	Progress string `json:"progress"`
}

func (p *pullInfo) String() string {
	if p.Status == "Downloading" {
		return fmt.Sprintf("%s %s", p.ID, p.Progress)
	}
	return fmt.Sprintf("%s %s", p.ID, p.Status)
}
