package app

import (
	"easycoding/internal/daemon"
	"time"
)

func (k *Kernel) StartDaemons() {
	m := daemon.NewManager(
		k.context.ctx,
		k.wg,
		k.Log,
		time.Duration(k.Config.Daemon.DurationSeconds)*time.Second,
	)

	go m.Start(daemon.NewExampleDaemon())
}
