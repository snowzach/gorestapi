package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Wrapper implements various useful log functions with the provided logger
// to fulfill interfaces.
type Wrapper struct {
	logger *zap.Logger
	level  zapcore.Level
}

func NewWrapper(logger *zap.Logger, level zapcore.Level) *Wrapper {
	return &Wrapper{
		logger: logger.WithOptions(zap.AddCallerSkip(1)),
		level:  level,
	}
}

func (w *Wrapper) Printf(template string, args ...interface{}) {
	w.logger.Check(w.level, fmt.Sprintf(template, args...)).Write()
}
