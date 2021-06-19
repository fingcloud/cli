package config

import (
	"encoding/json"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

type Config struct {
	App      string `mapstructure:"app" json:"app"`
	Port     string `mapstructure:"port" json:"port"`
	Platform string `mapstructure:"platform" json:"platform"`
}

type AuthConfig struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

func ReadAuthConfig() (*AuthConfig, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	path := path.Join(home, ".fing")
	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := new(AuthConfig)
	err = json.Unmarshal(bs, cfg)

	return cfg, err
}

func WriteAuthConfig(cfg *AuthConfig) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	bs, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	path := path.Join(home, ".fing")
	return os.WriteFile(path, bs, 0644)
}
