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

// Stop is the global stop instance
var Stop = func() *stop {
	ctx, cancel := context.WithCancel(context.Background())
	return &stop{
		ctx:    ctx,
		cancel: cancel,
	}
}()

// Stop manually triggers stop
func (s *stop) Stop() {
	s.cancel()
}

// Chan returns a read only channel that is closed when the program should exit
func (s *stop) Chan() <-chan struct{} {
	return s.ctx.Done()
}

// Context returns a context tied to the stop handler
func (s *stop) Context() context.Context {
	return s.ctx
}

// Bool returns t/f if the stop handler has triggered
func (s *stop) Bool() bool {
	return s.ctx.Err() != nil
}

// InitInterrupt sets up stop handler to trigger from interrupt signal
func (s *stop) InitInterrupt() {

	// Handle signals
	signalChannel := make(chan os.Signal, 1)

	// Stop flag will indicate if Ctrl-C/Interrupt has been sent to the process
	signal.Notify(signalChannel, os.Interrupt)

	// Handle signals
	go func() {
		for {
			for sig := range signalChannel {
				switch sig {
				case os.Interrupt:
					zap.S().Info("Received Interrupt...")
					s.cancel()
					return
				}
			}
		}
	}()
}
