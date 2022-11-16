package glog

// default logger provided
/*var GlobalGLogger IGLogger = NewDefaultLogger(log.New(os.Stdout, "", log.LstdFlags|log.Llongfile))

func SetGlobalGLogger(l IGLogger) {
	GlobalGLogger = l
}
*/
/*type Logger IGLogger // alias for generated code

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
	GlobalGLogger.GDebug(msg, v...)
}

func Info(msg string, v ...Field) {
	GlobalGLogger.GInfo(msg, v...)
}

func Warn(msg string, v ...Field) {
	GlobalGLogger.GWarn(msg, v...)
}

func Error(msg string, v ...Field) {
	GlobalGLogger.GError(msg, v...)
}

func Panic(msg string, v ...Field) {
	GlobalGLogger.GPanic(msg, v...)
}

func Fatal(msg string, v ...Field) {
	GlobalGLogger.GFatal(msg, v...)
}*/
