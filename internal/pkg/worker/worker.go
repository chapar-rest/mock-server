package worker

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// Agent is a unit of periodic background work.
type Agent interface {
	Run(ctx context.Context) error
}

// Worker runs every registered agent once per tick. It runs as a goroutine
// inside the server process: all state is in memory, so there is no
// separate worker deployment.
type Worker struct {
	logger         *zap.Logger
	tickerInterval time.Duration

	agents map[string]Agent
}

func NewWorker(logger *zap.Logger, tickerInterval time.Duration) *Worker {
	return &Worker{
		logger:         logger,
		tickerInterval: tickerInterval,
		agents:         make(map[string]Agent),
	}
}

// Start blocks, running the agents on every tick until ctx is done.
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.tickerInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Shutting down worker")
			return
		case <-ticker.C:
			for name, agent := range w.agents {
				if err := agent.Run(ctx); err != nil {
					w.logger.Error("Error running agent", zap.String("agent", name), zap.Error(err))
				}
			}
		}
	}
}

// AddAgent registers an agent. It must be called before Start.
func (w *Worker) AddAgent(name string, agent Agent) error {
	if _, ok := w.agents[name]; ok {
		return fmt.Errorf("agent %s already exists", name)
	}
	w.agents[name] = agent
	return nil
}
