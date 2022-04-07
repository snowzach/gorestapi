package conf

import (
	"fmt"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// DefaultDecodeHooks are the default decoding hooks used to unmarshal into a struct.
// This includes hooks for parsing string to time.Duration, time.Time(RFC3339 format),
// net.IP and net.IPNet. You can use this function to grab the defaults plus add your
// own with the extras option.
func DefaultDecodeHooks(extras ...mapstructure.DecodeHookFunc) []mapstructure.DecodeHookFunc {
	return append([]mapstructure.DecodeHookFunc{
		mapstructure.RecursiveStructToMapHookFunc(),
		mapstructure.StringToTimeDurationHookFunc(),
		mapstructure.StringToTimeHookFunc(time.RFC3339),
		mapstructure.StringToIPHookFunc(),
		mapstructure.StringToIPNetHookFunc(),
		ZapLogLevelHookFunc(),
	}, extras...)
}

// ZapLogLevelHookFunc()
func ZapLogLevelHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, v interface{}) (interface{}, error) {

		if t != reflect.TypeOf(zapcore.Level(1)) && t != reflect.TypeOf(zap.AtomicLevel{}) {
			return v, nil
		}
		levelString, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("could not parse log level of type %T", v)
		}

		var level zapcore.Level
		if err := level.Set(levelString); err != nil {
			return nil, fmt.Errorf("could not parse log level: %s", levelString)
		}

		// If it's a regular level return that
		if t == reflect.TypeOf(zapcore.Level(1)) {
			return level, nil
		}

		// Otherwise return an atomic level
		var atomicLevel zap.AtomicLevel
		atomicLevel.SetLevel(level)
		return atomicLevel, nil

	}
}
