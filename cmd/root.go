package cmd

import (
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	cli "github.com/spf13/cobra"

	"github.com/snowzach/gorestapi/pkg/conf"
	"github.com/snowzach/gorestapi/pkg/log"
	"github.com/snowzach/gorestapi/pkg/version"
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")
}

var (

	// Config and global logger
	pidFile string
	cfgFile string

	// The Root Cli Handler
	rootCmd = &cli.Command{
		Version: version.GitVersion,
		Use:     version.Executable,
		PersistentPreRunE: func(cmd *cli.Command, args []string) error {

			// Parse defaults, config file and environment.
			if err := conf.C.Parse(
				conf.WithMap(defaults()),
				conf.WithFile(cfgFile),
				conf.WithEnv(),
			); err != nil {
				fmt.Printf("could not load config: %v", err)
				os.Exit(1)
			}

			var loggerConfig log.LoggerConfig
			if err := conf.C.Unmarshal(&loggerConfig, conf.UnmarshalConf{Path: "logger"}); err != nil {
				fmt.Printf("could not parse logger config: %v", err)
				os.Exit(1)
			}
			if err := log.InitLogger(&loggerConfig); err != nil {
				fmt.Printf("could not configure logger: %v", err)
				os.Exit(1)
			}

			// Load the metrics server
			if conf.C.Bool("metrics.enabled") {
				hostPort := net.JoinHostPort(conf.C.String("metrics.host"), conf.C.String("metrics.port"))
				r := http.NewServeMux()
				r.Handle("/metrics", promhttp.Handler())
				if conf.C.Bool("profiler.enabled") {
					r.HandleFunc("/debug/pprof/", pprof.Index)
					r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
					r.HandleFunc("/debug/pprof/profile", pprof.Profile)
					r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
					r.HandleFunc("/debug/pprof/trace", pprof.Trace)
					log.Infow("Profiler enabled", "profiler_path", fmt.Sprintf("http://%s/debug/pprof/", hostPort))
				}
				go func() {
					if err := http.ListenAndServe(hostPort, r); err != nil {
						log.Errorf("Metrics server error: %v", err)
					}
				}()
				log.Infow("Metrics enabled", "address", hostPort)
			}

			// Create Pid File
			pidFile = conf.C.String("pidfile")
			if pidFile != "" {
				file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
				if err != nil {
					return fmt.Errorf("could not create pid file: %s error:%v", pidFile, err)
				}
				defer file.Close()
				_, err = fmt.Fprintf(file, "%d\n", os.Getpid())
				if err != nil {
					return fmt.Errorf("could not create pid file: %s error:%v", pidFile, err)
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
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
}
