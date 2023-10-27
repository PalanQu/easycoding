package daemon

import (
	"context"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	stopTimeout      = 1 * time.Minute
	recorverInterval = 5 * time.Second
)

// Daemon defines all required methods of a background daemon process.
type Daemon interface {
	Name() string
	Run(ctx context.Context, wg *sync.WaitGroup, logger *logrus.Logger) error
}

// Manager takes care of starting and orchestrating daemon processes.
type Manager struct {
	log *logrus.Logger
	// wg is used to keep track of running daemons.
	wg *sync.WaitGroup
	// ctx is used to signal cancellation to running daemons.
	ctx context.Context
}

// NewManager returns a new manager instance.
func NewManager(
	ctx context.Context,
	wg *sync.WaitGroup,
	l *logrus.Logger,
) *Manager {
	return &Manager{
		ctx: ctx,
		wg:  wg,
		log: l,
	}
}

// Start starts a daemon.
func (m *Manager) Start(d Daemon) {
	m.wg.Add(1)
	defer m.wg.Done()

	var wg sync.WaitGroup
	var try int

	logger := m.log
	go m.recoverPanic(d, m.log)

	for {
		select {
		// If the daemon should stop, wait for the wg to be done, then
		// stop the restart ticket and exit the method.
		case <-m.ctx.Done():
			wg.Wait()

			shutdownComplete := make(chan struct{})
			go func() {
				defer close(shutdownComplete)
				wg.Wait()
			}()

			select {
			case <-shutdownComplete:
				logger.Infof(`daemon "%s" shutdown complete`, d.Name())
			case <-time.After(stopTimeout):
				logger.Warnf(`daemon "%s" shutdown timed out`, d.Name())
			}
			return
		default:
			try++
			logger.Tracef(`starting daemon "%s" (%d. try)...`, d.Name(), try)
			// Run the daemon. If it crashes, continue to the next ticker iteration.
			if err := d.Run(m.ctx, &wg, logger); err != nil {
				logger.Warnf(`daemon "%s" crashed: %s`, d.Name(), err)
				continue
			}
			// If the Run method returned without an error, reset the try counter
			// and restart the daemon again. All daemons are run forever, even
			// if the return without an error.
			logger.Debugf(`daemon "%s" exited without errors`, d.Name())
			try = 0
		}
	}
}

// recoverPanic recovers a crashed daemon and restarts it.
func (m *Manager) recoverPanic(d Daemon, logger *logrus.Logger) {
	if err := recover(); err != nil {
		logger.Errorf("daemon exited with panic (restarting in 5 seconds): %s", err)
		time.Sleep(recorverInterval)
		m.Start(d)
	}
}
