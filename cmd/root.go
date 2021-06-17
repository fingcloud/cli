package cmd

import (
	"strings"

	"github.com/fingcloud/fing-cli/api"
	"github.com/fingcloud/fing-cli/internal/config"
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
