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
	pidFile string
	logger  *zap.SugaredLogger

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

	// Load configuration
	_ = conf.Defaults(conf.C)
	if configFile := rootCmd.PersistentFlags().StringP("config", "c", "", "config file"); configFile != nil && *configFile != "" {
		_ = conf.File(conf.C, *configFile)
	}
	_ = conf.Env(conf.C)

	conf.InitLogger(conf.C)

	logger = zap.S().With("package", "cmd")

	if conf.C.Bool("profiler.enabled") {
		hostPort := net.JoinHostPort(conf.C.String("profiler.host"), conf.C.String("profiler.port"))
		go func() {
			if err := http.ListenAndServe(hostPort, nil); err != nil {
				logger.Errorf("profiler server error: %v", err)
			}
		}()
		logger.Infof("profiler enabled on http://%s", hostPort)
	}

	// Run the program
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
