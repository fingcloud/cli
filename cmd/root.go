package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fingcloud/cli/api"
	"github.com/fingcloud/cli/internal/config"
	"github.com/fingcloud/cli/internal/ui"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	cfgFile string
	token   string
	devMode bool
)

var rootCmd = &cobra.Command{
	Use:   "fing [COMMAND] [OPTIONS]",
	Short: "fing is a command line interface (CLI) for the Fing API",
}

var cfg = &config.Config{}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	initConfig()

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "fing", "fing config file")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "fing api token")
	rootCmd.PersistentFlags().BoolVar(&devMode, "dev", false, "development mode")

	cobra.OnInitialize(initConfig)

	opts := make([]api.Option, 0)
	if devMode {
		opts = append(opts, api.WithDevMode(true))
	}

	rootCmd.AddCommand(
		NewLoginCommand(),
		NewDeployCommand(),
		NewAppsCommand(),
	)
}

func initConfig() {
	viper.SetEnvPrefix("FING")
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.SetConfigName("fing")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.ReadInConfig()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, ui.Alert(err.Error()))
		os.Exit(1)
	}
}
