package daemon

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Example daemon
type Example struct {
	periodSeconds int
}

// NewExample returns a new example daemon.
func NewExampleDaemon(periodSeconds int) *Example {
	return &Example{
		periodSeconds: periodSeconds,
	}
}

func (g *Example) Name() string { return "example" }

// Example daemon logs to the console every few seconds.
func (g *Example) Run(ctx context.Context, wg *sync.WaitGroup, logger *logrus.Logger) error {
	wg.Add(1)
	defer wg.Done()
	go func() {
		for {
			select {
			case <-time.After(time.Duration(g.periodSeconds) * time.Second):
				fmt.Println("running example daemon")
			case <-ctx.Done():
				logger.Debug("example daemon is shutting down...")
				return
			}
		}
	}()
	<-ctx.Done()
	return nil
}
