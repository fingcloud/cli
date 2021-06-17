package cli

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fingcloud/cli/api"
	"github.com/fingcloud/cli/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	contextFilename = ".fing"
)

type FingCli struct {
	Cmd    *cobra.Command
	Args   []string
	Client *api.Client
	Config *config.Config
}

func New(cmd *cobra.Command, args []string, token string, devMode bool) *FingCli {
	client := api.NewClient(token, api.WithDevMode(devMode))

	config := new(config.Config)
	if err := viper.Unmarshal(config); err != nil {
		cobra.CheckErr(err)
	}
	return &FingCli{
		Cmd:    cmd,
		Args:   args,
		Client: client,
		Config: config,
	}
}

func (cli *FingCli) GetAccessToken() string {
	token := viper.GetString("token")
	if token != "" {
		return token
	}

	home, err := homedir.Dir()
	cobra.CheckErr(err)

	contextPath := filepath.Join(home, contextFilename)
	bs, err := os.ReadFile(contextPath)
	if err != nil {
		if os.IsExist(err) {
			log.Printf("can't read %s file: %v", contextPath, err)
		}
		return ""
	}

	return strings.Trim(string(bs), " \n")
}

func (cli *FingCli) SetAccessToken(token string) error {
	home, err := homedir.Dir()
	cobra.CheckErr(err)

	contextPath := filepath.Join(home, contextFilename)

	return ioutil.WriteFile(contextPath, []byte(token), 0644)
}
