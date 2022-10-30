package config

type LogConfig struct {
	Level string `mapstructure:"level" defaultvalue:"INFO"`
	Dir   string `mapstructure:"dir" defaultvalue:""`
}
