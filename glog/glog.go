package glog

import (
	"log"
	"os"
)

// default logger provided
var globalLogger IGLogger = NewDefaultLogger(log.New(os.Stdout, "", log.LstdFlags|log.Llongfile))

func SetGlobalGLogger(l IGLogger) {
	globalLogger = l
}

type Logger IGLogger // alias for generated code

type IGLogger interface {
	IGDepthLogger

	GDebug(msg string, v ...Field)
	GInfo(msg string, v ...Field)
	GWarn(msg string, v ...Field)
	GError(msg string, v ...Field)
	GPanic(msg string, v ...Field)
	GFatal(msg string, v ...Field)
}

func Debug(msg string, v ...Field) {
	globalLogger.GDebug(msg, v...)
}

func Info(msg string, v ...Field) {
	globalLogger.GInfo(msg, v...)
}

func Warn(msg string, v ...Field) {
	globalLogger.GWarn(msg, v...)
}

func Error(msg string, v ...Field) {
	globalLogger.GError(msg, v...)
}

func Panic(msg string, v ...Field) {
	globalLogger.GPanic(msg, v...)
}

func Fatal(msg string, v ...Field) {
	globalLogger.GFatal(msg, v...)
}
