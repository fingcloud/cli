package config

import (
	"os"
	"path/filepath"

	"github.com/fingcloud/cli/pkg/api"
	"gopkg.in/yaml.v3"
)

func ReadAppConfig(path string) (*api.AppConfig, error) {
	configPath := filepath.Join(path, "fing.yaml")
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	cfg := new(api.AppConfig)
	err = yaml.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func WriteAppConfig(path string, cfg *api.AppConfig) error {

	bs, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	configPath := filepath.Join(path, ".fing.yaml")
	return os.WriteFile(configPath, bs, 0644)
}
