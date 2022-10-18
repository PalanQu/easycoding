package config

var defaultConfig = Config{
	Server: ServerConfig{
		GatewayPort: "10000",
		GrpcPort:    "10001",
		SwaggerPort: "10002",
	},
	Database: DatabaseConfig{
		Host:           "localhost",
		Port:           "3306",
		User:           "root",
		Password:       "123456",
		DBName:         "test",
		CreateDatabase: true,
	},
	Log: LogConfig{
		Level: "INFO",
		Dir:   "",
	},
	Daemon: DaemonConfig{
		DurationSeconds: 10,
	},
}
