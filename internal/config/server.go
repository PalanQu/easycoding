package config

type ServerConfig struct {
	GatewayPort    string `mapstructure:"gateway_port" defaultvalue:"10000"`
	GrpcPort       string `mapstructure:"grpc_port" defaultvalue:"10001"`
	SwaggerPort    string `mapstructure:"swagger_port" defaultvalue:"10002"`
	RestartOnError bool   `mapstructure:"restart_on_error" defaultvalue:"false"`
}
