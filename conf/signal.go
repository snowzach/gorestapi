package conf

import (
	"os"
	"os/signal"
	"sync"

	"go.uber.org/zap"
)

type stop struct {
	// c is a channel that is closed when we are stopping
	c chan struct{}
	// WaitGroup is a embedded WaitGroup that will wait before exiting cleanly to allow for cleanup
	sync.WaitGroup
}

var (
	// Global Stop instance
	Stop = &stop{
		c: make(chan struct{}),
	}
	// Handle signals
	signalChannel = make(chan os.Signal, 1)
)

// Handles all incoming signals
func init() {

	// Stop flag will indicate if Ctrl-C/Interrupt has been sent to the process
	signal.Notify(signalChannel, os.Interrupt)

	// Handke signals
	go func() {
		for {
			for sig := range signalChannel {
				switch sig {
				case os.Interrupt:
					zap.S().Info("Received Interrupt...")
					close(Stop.c)
					return
				}
			}
		}
	}()

}

// Chan returns a read only channel that is closed when the program should exit
func (s *stop) Chan() <-chan struct{} {
	return s.c
}

// Bool returns t/f if we should stop
func (s *stop) Bool() bool {
	select {
	case <-s.c:
		return true
	default:
		return false
	}
}
