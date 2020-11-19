package conf

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"go.uber.org/zap"
)

type stop struct {
	// Used to signal when we are done
	ctx    context.Context
	cancel context.CancelFunc

	// WaitGroup is a embedded WaitGroup that will wait before exiting cleanly to allow for cleanup
	sync.WaitGroup
}

var (
	// Global Stop instance
	Stop = newStop()
)

func newStop() *stop {
	ctx, cancel := context.WithCancel(context.Background())
	return &stop{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *stop) SetupInterrupt() {

	// Handle signals
	signalChannel := make(chan os.Signal, 1)

	// Stop flag will indicate if Ctrl-C/Interrupt has been sent to the process
	signal.Notify(signalChannel, os.Interrupt)

	// Handke signals
	go func() {
		for {
			for sig := range signalChannel {
				switch sig {
				case os.Interrupt:
					zap.S().Info("Received Interrupt...")
					Stop.cancel()
					return
				}
			}
		}
	}()
}

// Chan returns a read only channel that is closed when the program should exit
func (s *stop) Chan() <-chan struct{} {
	return s.ctx.Done()
}

func (s *stop) Context() context.Context {
	return s.ctx
}

// Bool returns t/f if we should stop
func (s *stop) Bool() bool {
	return s.ctx.Err() != nil
}
