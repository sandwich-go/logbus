package logbus

// FIXME 用什么优化掉sync map
/*var loggerMap sync.Map

func resetLoggerMap() {
	loggerMap.Range(func(key, value interface{}) bool {
		loggerMap.Delete(key)
		return true
	})
}

func syncLoggerMap() {
	loggerMap.Range(func(key, value interface{}) bool {
		_ = value.(*StdLogger).Sync()
		return true
	})
}*/
