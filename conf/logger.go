package conf

import (
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() {

	logConfig := zap.NewProductionConfig()
	logConfig.Sampling = nil

	// Log Level
	var logLevel zapcore.Level
	if err := logLevel.Set(C.String("logger.level")); err != nil {
		zap.S().Fatalw("Could not determine logger.level", "error", err)
	}
	logConfig.Level.SetLevel(logLevel)

	// Handle different logger encodings
	loggerEncoding := C.String("logger.encoding")
	switch loggerEncoding {
	case "stackdriver":
		logConfig.Encoding = "json"
		logConfig.EncoderConfig = zapdriver.NewDevelopmentEncoderConfig()
	default:
		logConfig.Encoding = loggerEncoding
		// Enable Color
		if C.Bool("logger.color") {
			logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}
		logConfig.DisableStacktrace = C.Bool("logger.disable_stacktrace")
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
	logConfig.Development = C.Bool("logger.dev_mode")
	logConfig.DisableCaller = C.Bool("logger.disable_caller")

	// Build the logger
	globalLogger, _ := logConfig.Build()
	zap.ReplaceGlobals(globalLogger)

}
