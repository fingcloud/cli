package deploy

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/avast/retry-go"
	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/cmd/logs"
	"github.com/fingcloud/cli/pkg/cmd/util"
	"github.com/fingcloud/cli/pkg/fileutils"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"go.uber.org/atomic"
	"gopkg.in/yaml.v3"
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

	cmd.Flags().StringVarP(&o.config.App, "app", "a", o.config.App, "app name")
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

	configPath := filepath.Join(o.Path, "fing.yaml")
	f, err := os.Open(configPath)
	if err == nil {
		err = yaml.NewDecoder(f).Decode(o.config)
		if err != nil {
			return err
		}
	}

	if len(args) == 1 {
		o.config.App = args[0]
	}

	if o.config.App == "" {
		apps, err := ctx.Client.AppsList(&api.ListAppsOptions{})
		util.CheckErr(err)

		appOptions := funk.Map(apps, func(app *api.App) string { return app.Name }).([]string)

		err = ui.PromptSelect("Choose your app", appOptions, o.config.App)
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

	if o.config.Platform == "" {
		return fmt.Errorf("platform not specified")
	}

	if o.config.Port == 0 {
		return fmt.Errorf("port not specified")
	}

	return nil
}

func (o *DeployOptions) Run(ctx *cli.Context) error {

	o.printAppInfo()
	fmt.Println(ui.Info("Getting files..."))
	files, err := fileutils.GetFiles(o.Path)
	if err != nil {
		return err
	}

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
		return err
	}

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
	fmt.Println(ui.Info("Getting changed files..."))
	tarBuf := new(bytes.Buffer)
	err := fileutils.Compress(projectPath, files, tarBuf)
	if err != nil {
		return err
	}

	fmt.Println(ui.Info("Uploading..."))

	bar := ui.NewProgress(0)
	updateProgress := func(n int64, max int64) {
		bar.ChangeMax64(max)
		bar.Set64(n)
	}

	return ctx.Client.AppsUploadFiles(app, tarBuf, updateProgress)
}

func readBuildLogs(ctx *cli.Context, app string, deploymentId int64) error {
	fmt.Println(ui.Info("Building..."))

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
				fmt.Println(ui.Warning("Cancelling..."))
				canceled.Store(true)
				ctx.Client.DeploymentCancel(app, deploymentId)
				return
			}
		}
	}()

	var from int64
	var starting bool
	for {
		buildLogs, err := ctx.Client.DeploymentBuildLogs(app, deploymentId, &api.LogsOptions{From: from})
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
				fmt.Println(ui.Gray(log.Message))
				if strings.HasPrefix(log.Message, "Successfully tagged") {
					fmt.Println(ui.Info("Build completed"))
					stopCh <- true
				}
			}
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusFailed {
			return errors.New("Build failed")
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusCancel {
			return errors.New("Build canceled")
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusStarting && !starting {
			fmt.Println(ui.Info("Starting..."))
			starting = true
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusRunning {
			fmt.Println(ui.Info("App Started successfuly"))
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
	fmt.Printf("%s %s\n", ui.Gray("platform:"), ui.Green(o.config.Platform))
	fmt.Printf("%s %d\n", ui.Gray("port:"), ui.Green(o.config.Port))
}
