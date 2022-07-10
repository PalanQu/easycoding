package config

type ServerConfig struct {
	GatewayPort    string `mapstructure:"gateway_port"`
	GrpcPort       string `mapstructure:"grpc_port"`
	SwaggerPort    string `mapstructure:"swagger_port"`
	RestartOnError bool   `mapstructure:"restart_on_error"`
}
