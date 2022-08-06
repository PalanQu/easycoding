package config

import (
	"log"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"serveer"`
	Database DatabaseConfig `mapstructure:"database"`
	Log      LogConfig      `mapstructure:"log"`
}

func LoadConfig(configPath string) *Config {
	SetDefaults()
	if err := loadConfig(configPath); err != nil {
		log.Println("warning, config file not found, use default config", err)
	}
	return &Config{
		Server: ServerConfig{
			GatewayPort:    viper.GetString("server.gateway_port"),
			GrpcPort:       viper.GetString("server.grpc_port"),
			SwaggerPort:    viper.GetString("server.swagger_port"),
			RestartOnError: viper.GetBool("server.restart_on_error"),
		},
		Database: DatabaseConfig{
			Host:           viper.GetString("database.host"),
			Port:           viper.GetString("database.port"),
			User:           viper.GetString("database.user"),
			Password:       viper.GetString("database.password"),
			DBName:         viper.GetString("database.db_name"),
			CreateDatabase: viper.GetBool("database.create_database"),
		},
		Log: LogConfig{
			Level: viper.GetString("log.level"),
			Dir:   viper.GetString("log.dir"),
		},
	}
}

func loadConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("EASYCODING")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	c := Config{}
	bindEnvs(c)
	viper.Unmarshal(&c)

	viper.WatchConfig()
	return nil
}

func SetDefaults() {
	viper.SetDefault("server.gateway_port", defaultConfig.Server.GatewayPort)
	viper.SetDefault("server.grpc_port", defaultConfig.Server.GrpcPort)
	viper.SetDefault("server.swagger_port", defaultConfig.Server.SwaggerPort)
	viper.SetDefault("server.restart_on_error", defaultConfig.Server.RestartOnError)

	viper.SetDefault("database.host", defaultConfig.Database.Host)
	viper.SetDefault("database.port", defaultConfig.Database.Port)
	viper.SetDefault("database.user", defaultConfig.Database.User)
	viper.SetDefault("database.password", defaultConfig.Database.Password)
	viper.SetDefault("database.db_name", defaultConfig.Database.DBName)
	viper.SetDefault("database.create_database", defaultConfig.Database.CreateDatabase)

	viper.SetDefault("log.level", defaultConfig.Log.Level)
	viper.SetDefault("log.dir", defaultConfig.Log.Dir)
}

func bindEnvs(iface interface{}, parts ...string) {
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
			bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}
