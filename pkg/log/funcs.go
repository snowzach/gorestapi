package log

import "go.uber.org/zap"

func With(args ...interface{}) *zap.SugaredLogger {
	return Logger.With(args...)
}

func Debug(args ...interface{}) {
	funcLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	funcLogger.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	funcLogger.Debugw(msg, keysAndValues...)
}

func Info(args ...interface{}) {
	funcLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	funcLogger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	funcLogger.Infow(msg, keysAndValues...)
}

func Warn(args ...interface{}) {
	funcLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	funcLogger.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	funcLogger.Warnw(msg, keysAndValues...)
}

func Error(args ...interface{}) {
	funcLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	funcLogger.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	funcLogger.Errorw(msg, keysAndValues...)
}

func Fatal(args ...interface{}) {
	funcLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	funcLogger.Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	funcLogger.Fatalw(msg, keysAndValues...)
}

func Panic(args ...interface{}) {
	funcLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	funcLogger.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	funcLogger.Panicw(msg, keysAndValues...)
}

func DPanic(args ...interface{}) {
	funcLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	funcLogger.DPanicf(template, args...)
}

func DPanicw(msg string, keysAndValues ...interface{}) {
	funcLogger.DPanicw(msg, keysAndValues...)
}

func Println(msg string) {
	funcLogger.Info(msg)
}

func Printf(template string, args ...interface{}) {
	funcLogger.Infof(template, args...)
}

func Flush() {
	_ = Base.Sync()
}
