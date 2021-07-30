package config

import (
	"encoding/json"
	"os"
	"path"
	"time"

	"github.com/mitchellh/go-homedir"
)

type FingConfig struct {
	LastCheckedAt time.Time `json:"last_checked_at"`
}

func ReadFingConfig() (*FingConfig, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	cfg := new(FingConfig)

	configPath := path.Join(home, ".fing")
	fingConfigPath := path.Join(configPath, "config.json")

	bs, err := os.ReadFile(fingConfigPath)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(bs, cfg)
	return cfg, err
}

func WriteFingConfig(cfg *FingConfig) error {
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
	fingConfigPath := path.Join(configPath, "config.json")

	return os.WriteFile(fingConfigPath, bs, 0644)
}
