package cmd

import (
	"fmt"
	"net"
	"os"

	"net/http"
	_ "net/http/pprof" // Import for pprof

	cli "github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/snowzach/gorestapi/conf"
)

var (

	// Config and global logger
	configFile string
	pidFile    string
	logger     *zap.SugaredLogger

	// The Root Cli Handler
	rootCmd = &cli.Command{
		Version: conf.GitVersion,
		Use:     conf.Executable,
		PersistentPreRunE: func(cmd *cli.Command, args []string) error {
			// Create Pid File
			pidFile = conf.C.String("pidfile")
			if pidFile != "" {
				file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
				if err != nil {
					return fmt.Errorf("Could not create pid file: %s Error:%v", pidFile, err)
				}
				defer file.Close()
				_, err = fmt.Fprintf(file, "%d\n", os.Getpid())
				if err != nil {
					return fmt.Errorf("Could not create pid file: %s Error:%v", pidFile, err)
				}
			}
			return nil
		},
		PersistentPostRun: func(cmd *cli.Command, args []string) {
			// Remove Pid file
			if pidFile != "" {
				os.Remove(pidFile)
			}
		},
	}
)

// Execute starts the program
func Execute() {
	// Run the program
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

// This is the main initializer handling cli, config and log
func init() {
	// Initialize configuration
	cli.OnInitialize(initLog, initProfiler)
}

func initLog() {

	conf.InitLogger()
	logger = zap.S().With("package", "cmd")

}

// Profliter can explicitly listen on address/port
func initProfiler() {
	if conf.C.Bool("profiler.enabled") {
		hostPort := net.JoinHostPort(conf.C.String("profiler.host"), conf.C.String("profiler.port"))
		go http.ListenAndServe(hostPort, nil)
		logger.Infof("Profiler enabled on http://%s", hostPort)
	}
}
