package application

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Runner struct {
	components []Component
	defers     []func() error
	cancel     context.CancelFunc
}

func (a *Runner) Defer(fn func() error) {
	a.defers = append(a.defers, fn)
}

func NewAppRunner() *Runner {
	return &Runner{}
}

func (a *Runner) Register(c Component) {
	a.components = append(a.components, c)
}

func (a *Runner) StartBlocking() error {
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	for _, c := range a.components {
		log.Printf("Starting %s...", c.Name())
		if err := c.Start(ctx); err != nil {
			return fmt.Errorf("failed to start %s: %w", c.Name(), err)
		}
	}

	log.Println("All components started. Waiting for shutdown signal...")

	<-sigCh
	log.Println("Shutdown signal received. Stopping components...")
	a.StopAll(context.Background())
	return nil
}

func (a *Runner) StopAll(ctx context.Context) {
	for _, c := range a.components {
		log.Printf("Stopping %s...", c.Name())
		if err := c.Stop(ctx); err != nil {
			log.Printf("Error stopping %s: %v", c.Name(), err)
		}
	}

	for i := len(a.defers) - 1; i >= 0; i-- {
		err := a.defers[i]()
		if err != nil {
			log.Printf("Error running defer: %v", err)
		}
	}

	a.cancel()
}
