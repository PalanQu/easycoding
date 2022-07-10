package config

type LogConfig struct {
	Level string `mapstructure:"level"`
	Dir   string `mapstructure:"dir"`
}
