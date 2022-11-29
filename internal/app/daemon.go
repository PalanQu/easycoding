package app

import (
	"easycoding/internal/daemon"
)

func (k *Kernel) StartDaemons() {
	m := daemon.NewManager(
		k.context.ctx,
		k.wg,
		k.log,
	)

	go m.Start(daemon.NewExampleDaemon(k.config.Daemon.ExampleDaemon.DurationSeconds))
}
