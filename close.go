package logbus

import "sync"

// Close 程序结束时打印缓存中的所有日志 并清理资源
func Close() {
	mutexLogger.RLock()
	defer mutexLogger.RUnlock()
	for _, v := range loggerExist {
		_ = v.Sync()
		v.syncDepthLogger()
	}
}

var mutexLogger sync.RWMutex
var loggerExist []*GLogger

// PushGLogger 推入一个GLogger,缓存由 logbus 创建出去的所有logger便于统一close
func PushGLogger(logger *GLogger) {
	mutexLogger.Lock()
	defer mutexLogger.Unlock()
	loggerExist = append(loggerExist, logger)
}
