package config

import (
	"log"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
	Daemon   DaemonConfig   `mapstructure:"daemon"`
}

func LoadConfig(configPath string) *Config {
	c, err := loadConfig(configPath)
	if err != nil {
		log.Println("warning, config file not found, use default config", err)
	}
	return c
}

func loadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("EASYCODING")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	c := Config{}
	bindValues(c)
	viper.Unmarshal(&c)

	viper.WatchConfig()
	return &c, nil
}

func bindValues(iface interface{}, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			bindValues(v.Interface(), append(parts, tv)...)
		default:
			key := strings.Join(append(parts, tv), ".")
			viper.BindEnv(key)
			dv, ok := t.Tag.Lookup("defaultvalue")
			if !ok {
				continue
			}
			viper.SetDefault(key, dv)
		}
	}
}
