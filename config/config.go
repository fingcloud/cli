package config

type Config struct {
	App      string `mapstructure:"app" json:"app"`
	Port     string `mapstructure:"port" json:"port"`
	Platform string `mapstructure:"platform" json:"platform"`
}
