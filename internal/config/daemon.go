package config

type DaemonConfig struct {
	ExampleDaemon ExampleDaemon `mapstructure:"example_daemon"`
}

type ExampleDaemon struct {
	DurationSeconds int `mapstructure:"example_daemon" defaultvalue:"5"`
}
