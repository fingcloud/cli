package cmd

import (
	"fmt"
	"time"

	"github.com/fingcloud/cli/api"
	"github.com/fingcloud/cli/internal/cli"
	"github.com/fingcloud/cli/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewLogsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "show app logs",
		Run: func(cmd *cobra.Command, args []string) {
			cli := cli.New(cmd, args, token, devMode)

			app := viper.GetString("app")
			since := viper.GetDuration("since")
			timestamps := viper.GetBool("timestamps")
			follow := viper.GetBool("follow")

			runLogs(cli, app, since, timestamps, follow)
		},
	}

	cmd.Flags().StringP("app", "a", "", "app name")
	viper.BindPFlag("app", cmd.Flags().Lookup("app"))

	cmd.Flags().BoolP("follow", "f", false, "follow log output")
	viper.BindPFlag("follow", cmd.Flags().Lookup("follow"))

	cmd.Flags().Duration("since", time.Minute, "show logs since relative time (e.g. 10m for 10 minutes)")
	viper.BindPFlag("since", cmd.Flags().Lookup("since"))

	cmd.Flags().Bool("timestamps", false, "show timestamps")
	viper.BindPFlag("timestamps", cmd.Flags().Lookup("timestamps"))

	return cmd
}

func runLogs(cli *cli.FingCli, app string, since time.Duration, timestamps, follow bool) {
	fmt.Println(ui.Info(fmt.Sprintf("Reading %s logs...", app)))

	from := time.Now().Add(-since).Unix()
	for {
		logs, err := cli.Client.AppLogs(app, &api.AppLogsOptions{
			Since: from,
		})
		checkError(err)

		for _, log := range logs {
			fmt.Println(ui.Gray(log.Message))
		}

		if !follow {
			break
		}

		if len(logs) > 0 {
			from = logs[len(logs)-1].Timestamp + 1
		}

		time.Sleep(50 * time.Millisecond)
	}
}
