package cmd

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
	"github.com/fingcloud/cli/api"
	"github.com/fingcloud/cli/internal/cli"
	"github.com/fingcloud/cli/internal/helpers"
	"github.com/fingcloud/cli/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thoas/go-funk"
)

func NewDeployCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy application",
		Long:  `you can deploy your application using this command`,
		Run: func(cmd *cobra.Command, args []string) {
			cli := cli.New(cmd, args, token, devMode)
			runDeploy(cli)
		},
	}

	cmd.Flags().StringP("app", "a", "", "your app name")
	viper.BindPFlag("app", cmd.Flags().Lookup("app"))

	cmd.Flags().String("path", ".", "path to your app")
	viper.BindPFlag("path", cmd.Flags().Lookup("path"))

	return cmd
}

func runDeploy(cli *cli.FingCli) {
	app := viper.GetString("app")
	path := viper.GetString("path")

	if app == "" {
		apps, err := cli.Client.AppsList(&api.ListAppsOptions{})
		checkError(err)

		appOptions := funk.Map(apps, func(app *api.App) string { return app.Name }).([]string)

		err = ui.PromptSelect("Choose your app", appOptions, &app)
		checkError(err)
	}

	printAppInfo()

	fmt.Println(ui.Info("Getting files..."))
	files, err := helpers.GetFiles(path)
	checkError(err)

	var deployment *api.Deployment
	err = retry.Do(func() error {
		d, changes, err := cli.Client.DeployemntCreate(app, &api.CreateDeploymentOptions{
			Files:  files,
			Config: cli.Config,
		})
		if err != nil {
			return retry.Unrecoverable(err)
		}

		if changes != nil {
			err = deployUploadChanges(cli, path, app, changes)
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
	checkError(err)

	err = readBuildLogs(cli, app, deployment.ID)
	checkError(err)
}

func deployUploadChanges(cli *cli.FingCli, projectPath, app string, files []*api.FileInfo) error {
	fmt.Println(ui.Info("Getting changed files..."))
	tarBuf := new(bytes.Buffer)
	err := helpers.Compress(projectPath, files, tarBuf)
	if err != nil {
		return err
	}
	fmt.Println(ui.Info("Uploading..."))

	bar := ui.NewProgress(0)
	updateProgress := func(n int64, max int64) {
		bar.ChangeMax64(max)
		bar.Set64(n)
	}

	return cli.Client.AppsUploadFiles(app, tarBuf, updateProgress)
}

func readBuildLogs(cli *cli.FingCli, app string, deploymentId int64) error {
	fmt.Println(ui.Info("Building..."))

	interruptCh := make(chan os.Signal)
	builtCh := make(chan bool, 1)
	ticker := time.NewTicker(100 * time.Millisecond)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	defer func() {
		close(interruptCh)
		close(builtCh)
		ticker.Stop()
	}()

	var from int64
	var starting bool
	for {
		select {
		case <-builtCh:
			return nil
		case <-interruptCh:
			return cli.Client.DeploymentCancel(app, deploymentId)
		case <-ticker.C:
			buildLogs, err := cli.Client.DeploymentListBuildLogs(app, deploymentId, &api.ListLogsOptions{From: from})
			if err != nil {
				return err
			}

			go func() {
				for _, log := range buildLogs.Logs {
					if log.Message == "" {
						continue
					}
					fmt.Println(ui.Gray(log.Message))
					if strings.HasPrefix(log.Message, "Successfully tagged") {
						fmt.Println(ui.Info("Build completed"))
						builtCh <- true
					}
				}
			}()

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
				fmt.Println(ui.Info("Deployment created :)"))
				return nil
			}

			if len(buildLogs.Logs) > 0 {
				from = buildLogs.Logs[len(buildLogs.Logs)-1].ID
			}
		}
	}
}

func printAppInfo() {
	path, _ := filepath.Abs(viper.GetString("path"))
	fmt.Printf("%s %s\n", ui.Gray("path:"), ui.Green(path))
	fmt.Printf("%s %s\n", ui.Gray("app:"), ui.Green(viper.GetString("app")))
	fmt.Printf("%s %s\n", ui.Gray("platform:"), ui.Green(viper.GetString("platform")))
	fmt.Printf("%s %s\n", ui.Gray("port:"), ui.Green(viper.GetString("port")))
}
