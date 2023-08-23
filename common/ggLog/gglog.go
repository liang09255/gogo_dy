package ggLog

import (
	"go.uber.org/zap"
)

var defaultLogger *zap.Logger
var defaultSugarLogger *zap.SugaredLogger

func init() {
	var err error
	defaultLogger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defaultSugarLogger = defaultLogger.Sugar()
}

func Info(args ...interface{}) {
	defaultSugarLogger.Info(args...)
}

func Debug(args ...interface{}) {
	defaultSugarLogger.Debug(args...)
}

func Error(args ...interface{}) {
	defaultSugarLogger.Error(args...)
}

func Warn(args ...interface{}) {
	defaultSugarLogger.Warn(args...)
}

func Fatal(args ...interface{}) {
	defaultSugarLogger.Fatal(args...)
}

func Panic(args ...interface{}) {
	defaultSugarLogger.Panic(args...)
}

func FastInfo(msg string, fields ...zap.Field) {
	defaultLogger.Info(msg, fields...)
}

func FastDebug(msg string, fields ...zap.Field) {
	defaultLogger.Debug(msg, fields...)
}

func FastError(msg string, fields ...zap.Field) {
	defaultLogger.Error(msg, fields...)
}

func FastWarn(msg string, fields ...zap.Field) {
	defaultLogger.Warn(msg, fields...)
}

func FastFatal(msg string, fields ...zap.Field) {
	defaultLogger.Fatal(msg, fields...)
}

func FastPanic(msg string, fields ...zap.Field) {
	defaultLogger.Panic(msg, fields...)
}

func Infof(template string, args ...interface{}) {
	defaultSugarLogger.Infof(template, args...)
}

func Debugf(template string, args ...interface{}) {
	defaultSugarLogger.Debugf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	defaultSugarLogger.Errorf(template, args...)
}

func Warnf(template string, args ...interface{}) {
	defaultSugarLogger.Warnf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	defaultSugarLogger.Fatalf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	defaultSugarLogger.Panicf(template, args...)
}
