package cmd

import (
	"flag"
	"io"
	"os"

	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/cmd/auth"
	"github.com/fingcloud/cli/pkg/cmd/deploy"
	"github.com/fingcloud/cli/pkg/cmd/logs"
	"github.com/fingcloud/cli/pkg/cmd/version"
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
	cmd.AddCommand(logs.NewCmdLogs(ctx))
	cmd.AddCommand(deploy.NewCmdDeploy(ctx))

	return cmd
}

func Execute() {
	fingCmd := NewFingCommand(os.Stdin, os.Stdout, os.Stderr)
	cobra.CheckErr(fingCmd.Execute())
}
