package cmd

import (
	"fmt"
	"net"
	"os"
	"strings"

	"net/http"
	_ "net/http/pprof" // Import for pprof

	cli "github.com/spf13/cobra"
	config "github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

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
			pidFile = config.GetString("pidfile")
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
	cli.OnInitialize(initConfig, initLog, initProfiler)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Config file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Sets up the config file, environment etc
	config.SetTypeByDefaultValue(true)                      // If a default value is []string{"a"} an environment variable of "a b" will end up []string{"a","b"}
	config.AutomaticEnv()                                   // Automatically use environment variables where available
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Environement variables use underscores instead of periods

	// If a config file is found, read it in.
	if configFile != "" {
		config.SetConfigFile(configFile)
		err := config.ReadInConfig()
		if err != nil {
			zap.S().Fatalw("Could not read config file",
				"file", configFile,
				"error", err,
			)
		}
	}
}

func initLog() {

	logConfig := zap.NewProductionConfig()

	// Log Level
	var logLevel zapcore.Level
	if err := logLevel.Set(config.GetString("logger.level")); err != nil {
		zap.S().Fatalw("Could not determine logger.level", "error", err)
	}
	logConfig.Level.SetLevel(logLevel)

	// Settings
	logConfig.Encoding = config.GetString("logger.encoding")
	logConfig.Development = config.GetBool("logger.dev_mode")
	logConfig.DisableCaller = config.GetBool("logger.disable_caller")
	logConfig.DisableStacktrace = config.GetBool("logger.disable_stacktrace")

	// Enable Color
	if config.GetBool("logger.color") {
		logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Use sane timestamp when logging to console
	if logConfig.Encoding == "console" {
		logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// JSON Fields
	logConfig.EncoderConfig.MessageKey = "msg"
	logConfig.EncoderConfig.LevelKey = "level"
	logConfig.EncoderConfig.CallerKey = "caller"

	// Build the logger
	globalLogger, _ := logConfig.Build()
	zap.ReplaceGlobals(globalLogger)
	logger = globalLogger.Sugar().With("package", "cmd")

}

// Profliter can explicitly listen on address/port
func initProfiler() {
	if config.GetBool("profiler.enabled") {
		hostPort := net.JoinHostPort(config.GetString("profiler.host"), config.GetString("profiler.port"))
		go http.ListenAndServe(hostPort, nil)
		logger.Infof("Profiler enabled on http://%s", hostPort)
	}
}
