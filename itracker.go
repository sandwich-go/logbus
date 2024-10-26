package logbus

import "github.com/sandwich-go/logbus/thinkingdata"

// ITracker thinkingData和bigQuery日志输出
type ITracker interface {
	Track(...Field) error
	TrackWithTGAData(thinkingdata.Data) error
}
