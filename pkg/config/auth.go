package config

import (
	"encoding/json"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
)

type AuthConfig struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

func ReadAuthConfig() (*AuthConfig, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	cfg := new(AuthConfig)

	configPath := path.Join(home, ".fing")
	authPath := path.Join(configPath, "auth.json")

	bs, err := os.ReadFile(authPath)
	if err != nil {
		return cfg, err
	}

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

	configPath := path.Join(home, ".fing")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		os.MkdirAll(configPath, os.ModePerm)
	}
	authPath := path.Join(configPath, "auth.json")

	return os.WriteFile(authPath, bs, 0644)
}
