package logbus

// ITracker thinkingData和bigQuery日志输出
type ITracker interface {
	Track(...Field) error
}
