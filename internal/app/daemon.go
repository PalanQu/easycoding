package app

import (
	"easycoding/internal/daemon"
)

func (k *Kernel) StartDaemons() {
	m := daemon.NewManager(
		k.context.ctx,
		k.wg,
		k.Log,
	)

	go m.Start(daemon.NewExampleDaemon(k.Config.Daemon.ExampleDaemon.DurationSeconds))
}
