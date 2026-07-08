package logbus

import "sync"

// Close 程序结束时打印缓存中的所有日志 并清理资源
func Close() {
	closeDynamicLogLevel()
	mutexLogger.RLock()
	for _, v := range loggerExist {
		_ = v.Sync()
		v.syncDepthLogger()
	}
	mutexLogger.RUnlock()

	// 关闭全局 track file writer（若启用），等待 worker 排空队列、fsync、释放文件句柄。
	// Close 之后写入会被静默丢弃（TrackFileWriteSyncer.Write 内部已处理 closed 状态），
	// 保持与过往 Close 行为一致：进程收尾阶段不再新增写入。
	if gTrackFileWriteSyncer != nil {
		_ = gTrackFileWriteSyncer.Close()
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
