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

	"github.com/dustin/go-humanize"
	"github.com/fingcloud/cli/api"
	"github.com/fingcloud/cli/internal/cli"
	"github.com/fingcloud/cli/internal/helpers"
	"github.com/fingcloud/cli/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	printAppInfo()

	fmt.Println(ui.Info("Getting files..."))
	files, err := helpers.GetFiles(path)
	checkError(err)

	retry, maxRetry := 0, 3
DEPLOY:
	retry++
	deployment, changes, err := cli.Client.DeployemntCreate(app, &api.CreateDeploymentOptions{
		Files:  files,
		Config: cli.Config,
	})
	checkError(err)

	if changes != nil {
		deployUploadChanges(cli, path, app, changes)
		if retry == maxRetry {
			checkError(errors.New("can't create deployment"))
		}
		goto DEPLOY
	}

	err = readBuildLogs(cli, app, deployment.ID)
	checkError(err)

	fmt.Println(deployment)
}

func deployUploadChanges(cli *cli.FingCli, projectPath, app string, files []*api.FileInfo) {
	fmt.Println(ui.Info("Getting files..."))
	tarBuf := new(bytes.Buffer)
	err := helpers.Compress(projectPath, files, tarBuf)
	checkError(err)

	fmt.Println(ui.Info(fmt.Sprintf("Upload size %s, (%d files)", humanize.Bytes(uint64(tarBuf.Len())), len(files))))
	fmt.Println(ui.Info("Uploading..."))

	err = cli.Client.AppsUploadFiles(app, tarBuf)
	checkError(err)
}

func readBuildLogs(cli *cli.FingCli, app string, deploymentId int64) error {
	fmt.Println(ui.Info("Building..."))

	interrupt := make(chan os.Signal)
	built := make(chan bool)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var canceled bool
	go func() {
		defer func() {
			close(interrupt)
			close(built)
		}()

		for {
			select {
			case <-interrupt:
				err := cli.Client.DeploymentCancel(app, deploymentId)
				if err != nil {
					checkError(err)
				}
				fmt.Println("")
				return
			case <-built:
				return
			}
		}
	}()

	var from int64
	var starting bool
	for {
		select {
		case <-interrupt:
			err := cli.Client.DeploymentCancel(app, deploymentId)
			if err != nil {
				checkError(err)
			}
		}
		buildLogs, err := cli.Client.DeploymentListBuildLogs(app, deploymentId, &api.ListLogsOptions{From: from})
		if err != nil {
			return err
		}

		for _, log := range buildLogs.Logs {
			fmt.Println(ui.Gray(log.Message))
			if strings.HasPrefix(log.Message, "Successfully tagged") {
				fmt.Println(ui.Info("Build completed"))
				built <- true
			}
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusFailed {
			return errors.New("Build failed")
		}

		if buildLogs.Deployment.Status == api.DeploymentStatusCancel || canceled {
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

func printAppInfo() {
	path, _ := filepath.Abs(viper.GetString("path"))
	fmt.Printf("%s %s\n", ui.Gray("path:"), ui.Green(path))
	fmt.Printf("%s %s\n", ui.Gray("app:"), ui.Green(viper.GetString("app")))
	fmt.Printf("%s %s\n", ui.Gray("platform:"), ui.Green(viper.GetString("platform")))
	fmt.Printf("%s %s\n", ui.Gray("port:"), ui.Green(viper.GetString("port")))
}
