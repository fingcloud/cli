package app

import (
	"fmt"
	"time"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/fingcloud/cli/pkg/util"
	"github.com/spf13/cobra"
)

type LogsOptions struct {
	App        string
	Since      time.Duration
	Follow     bool
	Timestamps bool
}

func NewCmdLogs(ctx *cli.Context) *cobra.Command {
	opts := new(LogsOptions)

	cmd := &cobra.Command{
		Use:   "logs [app]",
		Short: "show app logs",
		Args:  cli.Exact(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()

			opts.App = args[0]

			util.CheckErr(RunLogs(ctx, opts))
		},
	}

	cmd.Flags().BoolVarP(&opts.Follow, "follow", "f", false, "follow log output")
	cmd.Flags().DurationVar(&opts.Since, "since", time.Minute, "show logs since relative time (e.g. 10m for 10 minutes)")
	cmd.Flags().BoolVar(&opts.Timestamps, "timestamps", false, "show timestamps")

	return cmd
}

func RunLogs(ctx *cli.Context, opts *LogsOptions) error {
	fmt.Println(ui.Info(fmt.Sprintf("Reading %s logs...", opts.App)))

	from := time.Now().Add(-opts.Since).Unix()
	for {
		logs, err := ctx.Client.AppLogs(opts.App, &api.AppLogsOptions{
			Since: from,
		})
		util.CheckErr(err)

		for _, log := range logs {
			var timestamp string
			if opts.Timestamps {
				timestamp = time.Unix(log.Timestamp, 0).Format("2006-01-02 15:04:05 ")
			}
			fmt.Printf("%s%s\n", timestamp, log.Message)
		}

		if !opts.Follow {
			break
		}

		if len(logs) > 0 {
			from = logs[len(logs)-1].Timestamp + 1
		}

		time.Sleep(200 * time.Millisecond)
	}
	return nil
}
