package cmd

import (
	"fmt"
	"net"
	"os"

	"net/http"
	"net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	cli "github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/snowzach/gorestapi/conf"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
}

var (

	// Config and global logger
	pidFile string
	cfgFile string
	logger  *zap.SugaredLogger

	// The Root Cli Handler
	rootCmd = &cli.Command{
		Version: conf.GitVersion,
		Use:     conf.Executable,
		PersistentPreRunE: func(cmd *cli.Command, args []string) error {

			// Load configuration
			_ = conf.Defaults(conf.C)
			if cfgFile != "" {
				if err := conf.File(conf.C, cfgFile); err != nil {
					return fmt.Errorf("could not load config file %s: %v", cfgFile, err)
				}
			}
			_ = conf.Env(conf.C)

			conf.InitLogger(conf.C)

			logger = zap.S().With("package", "cmd")

			if conf.C.Bool("metrics.enabled") {

				hostPort := net.JoinHostPort(conf.C.String("metrics.host"), conf.C.String("metrics.port"))
				logger.Infow("Metrics enabled", "address", hostPort)

				r := http.NewServeMux()

				r.Handle("/metrics", promhttp.Handler())

				if conf.C.Bool("profiler.enabled") {
					r.HandleFunc("/debug/pprof/", pprof.Index)
					r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
					r.HandleFunc("/debug/pprof/profile", pprof.Profile)
					r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
					r.HandleFunc("/debug/pprof/trace", pprof.Trace)
					logger.Infow("Profiler enabled", "profiler_path", fmt.Sprintf("http://%s/debug/pprof/", hostPort))
				}

				go func() {
					if err := http.ListenAndServe(hostPort, r); err != nil {
						logger.Errorf("metrics server error: %v", err)
					}
				}()
			}

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
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
