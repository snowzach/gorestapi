package log

import (
	"fmt"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Defaults
var (
	Base       = zap.NewNop()
	Logger     = Base.Sugar()
	funcLogger = Base.WithOptions(zap.AddCallerSkip(1)).Sugar()
)

type LoggerConfig struct {
	Level             string `conf:"level"`
	Encoding          string `conf:"encoding"`
	Color             bool   `conf:"color"`
	DevMode           bool   `conf:"dev_mode"`
	DisableCaller     bool   `conf:"disable_caller"`
	DisableStacktrace bool   `conf:"disable_stacktrace"`
}

// InitLogger loads a global logger based on a configuration
func InitLogger(c *LoggerConfig) error {

	logConfig := zap.NewProductionConfig()
	logConfig.Sampling = nil

	// Log Level
	var logLevel zapcore.Level
	if err := logLevel.Set(c.Level); err != nil {
		return fmt.Errorf("could not determine log level: %w", err)
	}
	logConfig.Level.SetLevel(logLevel)

	// Handle different logger encodings
	switch c.Encoding {
	case "stackdriver":
		logConfig.Encoding = "json"
		logConfig.EncoderConfig = zapdriver.NewDevelopmentEncoderConfig()
	default:
		logConfig.Encoding = c.Encoding
		// Enable Color
		if c.Color {
			logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		logConfig.DisableStacktrace = c.DisableStacktrace
		// Use sane timestamp when logging to console
		if logConfig.Encoding == "console" {
			logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		}

		// JSON Fields
		logConfig.EncoderConfig.MessageKey = "msg"
		logConfig.EncoderConfig.LevelKey = "level"
		logConfig.EncoderConfig.CallerKey = "caller"
	}

	// Settings
	logConfig.Development = c.DevMode
	logConfig.DisableCaller = c.DisableCaller

	// Build the logger
	globalLogger, err := logConfig.Build()
	if err != nil {
		return fmt.Errorf("could not build log config: %w", err)
	}
	zap.ReplaceGlobals(globalLogger)

	Base = zap.L()
	Logger = Base.Sugar()
	funcLogger = Base.WithOptions(zap.AddCallerSkip(1)).Sugar()

	return nil

}
