package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dustin/go-humanize"
	"github.com/fingcloud/cli/api"
	"github.com/fingcloud/cli/cli"
	"github.com/fingcloud/cli/internal/helpers"
	"github.com/fingcloud/cli/internal/spinner"
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

	spinner.Start("Getting files...")
	files, err := helpers.GetFiles(path)
	checkError(err)
	spinner.Stop()

	spinner.Start("Creating deployment...")

DEPLOY:
	deployment, changes, err := cli.Client.CreateDeployment(app, &api.CreateDeploymentOptions{
		Files:  files,
		Config: cli.Config,
	})
	spinner.Stop()
	checkError(err)

	if changes != nil {
		deployUploadChanges(cli, app, changes)
		goto DEPLOY
	}

	fmt.Println(deployment)

}

func deployUploadChanges(cli *cli.FingCli, app string, files []*api.FileInfo) {
	spinner.Start("Compressing...")
	tarBuf := new(bytes.Buffer)
	err := helpers.Compress(files, tarBuf)
	spinner.Stop()
	checkError(err)

	fmt.Printf("📦 Upload %d files (size %s)\n", len(files), humanize.Bytes(uint64(tarBuf.Len())))

	spinner.Start("Uploading...")
	err = cli.Client.AppsUploadFiles(app, tarBuf)
	checkError(err)
	spinner.Stop()
}

func printAppInfo() {
	path, _ := filepath.Abs(viper.GetString("path"))
	fmt.Printf("%s %s\n", ui.Gray("path:"), ui.Green(path))
	fmt.Printf("%s %s\n", ui.Gray("app:"), ui.Green(viper.GetString("app")))
	fmt.Printf("%s %s\n", ui.Gray("platform:"), ui.Green(viper.GetString("platform")))
	fmt.Printf("%s %s\n", ui.Gray("port:"), ui.Green(viper.GetString("port")))
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, ui.Alert(err.Error()))
		os.Exit(1)
	}
}
