package logs

import (
	"go.uber.org/zap"
)

type Field = zap.Field

var defaultLogger *zap.Logger

func Init() *zap.Logger {
	var err error
	defaultLogger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return defaultLogger
}

func Debug(msg string, fields ...Field) {
	defaultLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	defaultLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	defaultLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	defaultLogger.Error(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	defaultLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	defaultLogger.Fatal(msg, fields...)
}

func String(key string, value string) Field {
	return zap.String(key, value)
}

func Int(key string, value int) Field {
	return zap.Int(key, value)
}
