package logs

import (
	"fmt"
	"time"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/cli"
	"github.com/fingcloud/cli/pkg/command/util"
	"github.com/fingcloud/cli/pkg/config"
	"github.com/fingcloud/cli/pkg/ui"
	"github.com/spf13/cobra"
)

type LogsOptions struct {
	App        string
	Since      time.Duration
	Follow     bool
	Timestamps bool
}

func NewCmdLogs(ctx *cli.Context) *cobra.Command {
	o := new(LogsOptions)

	cmd := &cobra.Command{
		Use:   "logs [app]",
		Short: "show application logs",
		Long:  "show application logs",
		Run: func(cmd *cobra.Command, args []string) {
			ctx.SetupClient()

			util.CheckErr(o.Init(ctx, args))
			util.CheckErr(o.Validate())
			util.CheckErr(o.Run(ctx))
		},
	}

	cmd.Flags().BoolVarP(&o.Follow, "follow", "f", false, "follow log output")
	cmd.Flags().DurationVar(&o.Since, "since", time.Minute, "show logs since relative time (e.g. 10m for 10 minutes)")
	cmd.Flags().BoolVar(&o.Timestamps, "timestamps", false, "show timestamps")

	return cmd
}

func (o *LogsOptions) Init(ctx *cli.Context, args []string) error {
	appConfig, err := config.ReadAppConfig(*ctx.Path)
	if err == nil {
		o.App = appConfig.App
	}

	if len(args) == 1 {
		o.App = args[0]
	}

	return nil
}
func (o *LogsOptions) Validate() error {

	return nil
}
func (o *LogsOptions) Run(ctx *cli.Context) error {
	fmt.Println(ui.Info(fmt.Sprintf("Reading %s logs...", o.App)))

	from := time.Now().Add(-o.Since).Unix()
	for {
		logs, err := ctx.Client.AppLogs(o.App, &api.AppLogsOptions{
			Since: from,
		})
		if err != nil {
			return err
		}

		for _, log := range logs {
			var timestamp string
			if o.Timestamps {
				timestamp = time.Unix(log.Timestamp, 0).Format("2006-01-02 15:04:05 ")
			}
			fmt.Printf("%s%s\n", timestamp, log.Message)
		}

		if !o.Follow {
			break
		}

		if len(logs) > 0 {
			from = logs[len(logs)-1].Timestamp + 1
		}

		time.Sleep(50 * time.Millisecond)
	}
	return nil
}
