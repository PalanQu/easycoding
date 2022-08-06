package config

type DatabaseConfig struct {
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	DBName         string `mapstructure:"db_name"`
	CreateDatabase bool   `mapstructure:"create_database"`
}
