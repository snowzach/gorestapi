package signal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

//  DefaultStopSignals is the SIGINT and SIGTERM signals
var DefaultStopSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}

type stop struct {
	// Used to signal when we are done
	ctx    context.Context
	cancel context.CancelFunc

	// WaitGroup is a embedded WaitGroup that will wait before exiting cleanly to allow for cleanup
	sync.WaitGroup
}

// Stop is the global stop instance if you wish to use.
var Stop = NewStop()

// NewStop creates a new stop instance
func NewStop() *stop {
	ctx, cancel := context.WithCancel(context.Background())
	return &stop{
		ctx:    ctx,
		cancel: cancel,
	}
}

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

// OnSignal sets up stop handler to trigger from the specified signals.
// If signals is not specific/nil it defaults to syscall.SIGINT and syscall.SIGTERM
func (s *stop) OnSignal(signals ...os.Signal) {

	if len(signals) == 0 {
		return
	}

	// Handle signals
	signalChannel := make(chan os.Signal, 1)

	// Get notified on signals.
	signal.Notify(signalChannel, signals...)

	// Handle signals
	go func() {
		<-signalChannel
		s.cancel()
	}()
}
