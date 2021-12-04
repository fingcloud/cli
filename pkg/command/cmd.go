package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Delta456/box-cli-maker/v2"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/command/app"
	"github.com/fingcloud/cli/pkg/command/auth"
	"github.com/fingcloud/cli/pkg/command/deploy"
	"github.com/fingcloud/cli/pkg/command/logs"
	"github.com/fingcloud/cli/pkg/command/version"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/fingcloud/cli/pkg/update"
	"github.com/spf13/cobra"
)

func NewFingCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fing",
		Short: "fing deploys and manages your apps on fing service",
	}

	flags := cmd.PersistentFlags()

	ctx := cli.NewContext()
	ctx.AddFlags(flags)

	cmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	cmd.AddCommand(version.NewCmdVersion(ctx))
	cmd.AddCommand(auth.NewCmdLogin(ctx))
	cmd.AddCommand(auth.NewCmdLogout(ctx))
	cmd.AddCommand(auth.NewCmdSetSession(ctx))
	cmd.AddCommand(logs.NewCmdLogs(ctx))
	cmd.AddCommand(deploy.NewCmdDeploy(ctx))
	cmd.AddCommand(app.NewAppsCmd(ctx))

	return cmd
}

func Execute() {
	updateChan := make(chan *update.Release)

	go func() {
		if !update.ShouldCheckUpdate() {
			updateChan <- nil
			return
		}

		release, err := update.CheckForUpdate(context.Background(), cli.Version)
		if err != nil {
			fmt.Println(ui.Warning("could not check for update"), err.Error())
		}
		updateChan <- release
	}()

	release := <-updateChan
	if release != nil {
		b := box.New(box.Config{Px: 2, Py: 1, Type: "Single", Color: "Yellow"})
		b.Println(fmt.Sprintf("Update Availabe %s -> %s", cli.Version, release.Version), update.UpdateCommand())
	}

	fingCmd := NewFingCommand(os.Stdin, os.Stdout, os.Stderr)
	cobra.CheckErr(fingCmd.Execute())
}
