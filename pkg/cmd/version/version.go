package version

import (
	"fmt"

	"github.com/fingcloud/cli/pkg/cli"
	"github.com/spf13/cobra"
)

type VersionOptions struct{}

func NewCmdVersion(ctx *cli.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "print cli version",
		Long:  "print cli version",

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Version:", cli.Version)
			fmt.Println("BuildDate:", cli.BuildDate)
			fmt.Println("Commit:", cli.Commit)
		},
	}

	return cmd
}
