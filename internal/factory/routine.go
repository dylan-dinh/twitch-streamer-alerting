package factory

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Routines struct {
	Name    string
	Routine func(ctx context.Context, wg *sync.WaitGroup, errChan chan error)
}

type RoutinesFactory struct {
	RoutinesFunc []Routines
	ctx          context.Context
	wg           *sync.WaitGroup
	ErrChan      chan error
	cancelFunc   context.CancelFunc
	signalChan   chan os.Signal
	logger       *slog.Logger
}

func NewRoutinesFactory(routines []Routines) *RoutinesFactory {
	ctx, cancel := context.WithCancel(context.Background())
	return &RoutinesFactory{
		RoutinesFunc: routines,
		ctx:          ctx,
		wg:           &sync.WaitGroup{},
		ErrChan:      make(chan error, 1),
		signalChan:   make(chan os.Signal, 1),
		cancelFunc:   cancel,
		logger:       slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

func (f *RoutinesFactory) StartRoutinesFactory() {
	for _, routine := range f.RoutinesFunc {
		f.wg.Add(1)
		go routine.Routine(context.Background(), f.wg, f.ErrChan)
		f.logger.Info("started routine", "routine", routine.Name)
	}
}

// StopRoutinesFactory will stop all the routines if any exiting signal is caught
func (f *RoutinesFactory) StopRoutinesFactory() {
	// Catch OS signals like Ctrl+C to shut down the server and routines gracefully
	signal.Notify(f.signalChan, os.Interrupt, syscall.SIGTERM)
	select {
	case sig := <-f.signalChan:
		f.logger.Log(f.ctx, slog.LevelError, "received signal, exiting server... ", "signal", sig)
		f.cancelFunc()
	case err := <-f.ErrChan:
		f.logger.Log(f.ctx, slog.LevelError, "errors in routines ", "error", err)
	}
	f.wg.Wait()
}
