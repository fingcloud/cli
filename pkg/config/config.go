package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/fingcloud/cli/pkg/api"
	"github.com/fingcloud/cli/pkg/config/encoding"
)

var (
	ErrConfigFileNotFound = errors.New("fing config file not found")

	configFiles = []string{
		"fing.yaml",
		"fing.yml",
		"fing.json",
		"fing.toml",
	}
)

func findConfigFile(path string) (string, error) {
	for _, configFile := range configFiles {
		configPath := filepath.Join(path, configFile)
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			continue
		}
		return configPath, nil
	}
	return "", ErrConfigFileNotFound
}

// ReadAppConfig reads config file from path
func ReadAppConfig(path string) (*api.AppConfig, error) {
	configPath, err := findConfigFile(path)
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cfg := new(api.AppConfig)
	err = encoding.Decode(filepath.Ext(configPath)[1:], bytes, cfg)

	return cfg, err
}

// WriteAppConfig writes config file to path
func WriteAppConfig(path, filename string, cfg *api.AppConfig) error {
	bytes, err := encoding.Encode(filepath.Ext(filename)[1:], cfg)
	if err != nil {
		return err
	}

	configPath := filepath.Join(path, filename)
	return os.WriteFile(configPath, bytes, 0644)
}
