package glog

type IGDepthLogger interface {
	// GDebugDepth DepthLoggerV2
	GDebugDepth(depth int, msg string, v ...Field)
	GInfoDepth(depth int, msg string, v ...Field)
	GWarnDepth(depth int, msg string, v ...Field)
	GErrorDepth(depth int, msg string, v ...Field)
	GFatalDepth(depth int, msg string, v ...Field)
}

// DebugDepth 用于glog被再封装
func DebugDepth(depth int, msg string, v ...Field) {
	globalLogger.GDebugDepth(depth, msg, v...)
}

// InfoDepth 用于glog被再封装
func InfoDepth(depth int, msg string, v ...Field) {
	globalLogger.GInfoDepth(depth, msg, v...)
}

// WarnDepth 用于glog被再封装
func WarnDepth(depth int, msg string, v ...Field) {
	globalLogger.GWarnDepth(depth, msg, v...)
}

// ErrorDepth 用于glog被再封装
func ErrorDepth(depth int, msg string, v ...Field) {
	globalLogger.GErrorDepth(depth, msg, v...)
}

// FatalDepth 用于glog被再封装
func FatalDepth(depth int, msg string, v ...Field) {
	globalLogger.GFatalDepth(depth, msg, v...)
}
