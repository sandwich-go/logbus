package logbus

import (
	"github.com/sandwich-go/logbus/bi"
	"github.com/sandwich-go/logbus/thinkingdata"
)

// ITracker thinkingData、bigQuery、BI 日志输出
type ITracker interface {
	Track(...Field) error
	TrackWithTGAData(thinkingdata.Data) error
	TrackWithBIData(bi.Data) error
}
