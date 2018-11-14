package conf

import (
	"os"
	"os/signal"
	"sync"

	"go.uber.org/zap"
)

var (
	// StopFlag is a global boolean for if we are stopping
	StopFlag bool
	// StopChan is a global channel that is closed when we are stopping
	StopChan = make(chan struct{})
	// StopWaitGroup is a global WaitGroup that will wait before exiting cleanly to allow for cleanup
	StopWaitGroup sync.WaitGroup

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
					StopFlag = true
					close(StopChan)
					return
				}
			}
		}
	}()

}
