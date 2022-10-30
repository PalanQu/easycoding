package config

type DatabaseConfig struct {
	Host           string `mapstructure:"host" defaultvalue:"localhost"`
	Port           string `mapstructure:"port" defaultvalue:"3306"`
	User           string `mapstructure:"user" defaultvalue:"root"`
	Password       string `mapstructure:"password" defaultvalue:"123456"`
	DBName         string `mapstructure:"db_name" defaultvalue:"test"`
	CreateDatabase bool   `mapstructure:"create_database" defaultvalue:"true"`
}
