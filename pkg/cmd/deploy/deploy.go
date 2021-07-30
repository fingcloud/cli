package deploy

import (
	"bytes"
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
	"github.com/fingcloud/cli/pkg/cmd/logs"
	"github.com/fingcloud/cli/pkg/cmd/util"
	"github.com/fingcloud/cli/pkg/config"
	"github.com/fingcloud/cli/pkg/fileutils"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/r6m/spinner"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"go.uber.org/atomic"
)

type DeployOptions struct {
	client *api.Client

	config *api.AppConfig
	logs   *logs.LogsOptions

	Path     string
	Quite    bool
	Dispatch bool
}

func NewCmdDeploy(ctx *cli.Context) *cobra.Command {
	o := NewOptions()

	cmd := &cobra.Command{
		Use:     "deploy",
		Short:   "deploy your application",
		Long:    "deploy your application",
		Aliases: []string{"up"},
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()

			util.CheckErr(o.Init(ctx, args))
			util.CheckErr(o.Validate())
			util.CheckErr(o.Run(ctx))
		},
	}

	cfg, err := config.ReadAppConfig(o.Path)
	if err == nil {
		o.config = cfg
	}

	cmd.Flags().StringVarP(&o.config.App, "app", "a", o.config.App, "app name")
	cmd.Flags().StringVar(&o.config.Platform, "platform", o.config.Platform, "your app platform")
	cmd.Flags().StringVar(&o.Path, "path", ".", "app path")
	cmd.Flags().BoolVarP(&o.Quite, "quite", "q", o.Quite, "quite output")
	cmd.Flags().BoolVarP(&o.Dispatch, "dispatch", "d", o.Quite, "dispatch logs")

	return cmd
}

func NewOptions() *DeployOptions {
	o := new(DeployOptions)
	o.config = &api.AppConfig{}
	return o
}

func (o *DeployOptions) Init(ctx *cli.Context, args []string) error {

	if len(args) == 1 {
		o.config.App = args[0]
	}

	if o.config.App == "" {
		apps, err := ctx.Client.AppsList(&api.ListAppsOptions{})
		util.CheckErr(err)

		appOptions := funk.Map(apps, func(app *api.App) string { return app.Name }).([]string)
		if len(appOptions) == 0 {
			fmt.Println("you don't have any apps on fing")
			fmt.Println("go to fing dashboard and create one:")
			fmt.Printf("\t%s\n\n", ui.Green("https://dashboard.fing.ir/apps"))
			return fmt.Errorf("empty apps")
		}

		err = ui.PromptSelect("Choose your app", appOptions, &o.config.App)
		util.CheckErr(err)
	}

	o.logs = &logs.LogsOptions{
		App:    o.config.App,
		Since:  time.Second,
		Follow: true,
	}

	return nil
}

func (o *DeployOptions) Validate() error {
	if o.config.App == "" {
		return fmt.Errorf("app can't be empty")
	}

	return nil
}

var s = spinner.NewSpinner().WithOptions(spinner.WithExitOnAbort(false))

func (o *DeployOptions) Run(ctx *cli.Context) error {
	o.printAppInfo()
	s.Start("Getting files...")

	files, err := fileutils.GetFiles(o.Path)
	if err != nil {
		return err
	}

	s.Success()
	s.Start("Creating Deployment...")

	var deployment *api.Deployment
	err = retry.Do(func() error {
		d, changes, err := ctx.Client.DeployemntCreate(o.config.App, &api.CreateDeploymentOptions{
			Files:  files,
			Config: o.config,
		})
		if err != nil {
			return retry.Unrecoverable(err)
		}

		if changes != nil {
			err = uploadChanges(ctx, o.Path, o.config.App, changes)
			if err != nil {
				return retry.Unrecoverable(err)
			}
			return errors.New("retry deployment")
		}

		deployment = d
		return nil
	},
		retry.Attempts(3),
	)
	if err != nil {
		s.Error(err.Error())
		return err
	}

	s.Success()
	s.Start("Analyzing...")
	if deployment.Platform != "" {
		fmt.Printf("%s %s\n", ui.Gray("platform:"), ui.Green(deployment.Platform))
	}
	s.Success()

	err = readBuildLogs(ctx, o.config.App, deployment.ID)
	if err != nil {
		return err
	}

	if !o.Dispatch {
		return o.logs.Run(ctx)
	}

	return nil
}

func uploadChanges(ctx *cli.Context, projectPath, app string, files []*api.FileInfo) error {
	s.UpdateMessage("Getting changed files...")
	fmt.Println(ui.Details(fmt.Sprintf("%d Files changed", len(files))))
	tarBuf := new(bytes.Buffer)
	err := fileutils.Compress(projectPath, files, tarBuf)
	if err != nil {
		s.Error()
		return err
	}

	s.UpdateMessage("Uploading changed files...")

	fmt.Println(ui.Details(ui.KeyValue("Upload size", humanize.Bytes(uint64(tarBuf.Len())))))
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

	go func() {
		defer func() {
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

		if !canceled.Load() {
			for _, log := range buildLogs.Logs {
				if canceled.Load() {
					break
				}
				if log.Message == "" {
					continue
				}

				s.ClearCurrentLine()
				fmt.Println(ui.Gray(log.Message))

				if strings.HasPrefix(log.Message, "Successfully tagged") {
					s.Success()
					stopCh <- true
				}
			}
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusFailed {
			s.Error("Build failed")
			return errors.New("Build failed")
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusCancel {
			s.Error("Build canceled")
			return errors.New("Build canceled")
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusStarting && !starting {
			s.Start("Starting...")
			starting = true
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusRunning {
			s.Success()
			fmt.Println(ui.Info("App started succesfully :)"))
			fmt.Println()
			fmt.Println(fmt.Sprintf("\topen the following url in your browser:"))
			fmt.Println(fmt.Sprintf("\t%s", ui.Green(buildLogs.Deployment.URL)))
			fmt.Println()
			return nil
		}

		if len(buildLogs.Logs) > 0 {
			from = buildLogs.Logs[len(buildLogs.Logs)-1].ID
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func (o *DeployOptions) printAppInfo() {
	fmt.Printf("%s %s\n", ui.Gray("app:"), ui.Green(o.config.App))
}
