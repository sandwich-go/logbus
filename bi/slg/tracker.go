package slg

import (
	"github.com/sandwich-go/logbus"
	"github.com/sandwich-go/logbus/bi"
	"github.com/sandwich-go/logbus/thinkingdata"
)

var (
	tracker logbus.ITracker = &noop{}
)

func MustInitialize(t logbus.ITracker) {
	if t == nil {
		panic("tracker is nil")
	}
	tracker = t
}

type noop struct{}

func (n noop) Track(field ...logbus.Field) error {
	return nil
}

func (n noop) TrackWithTGAData(data thinkingdata.Data) error {
	return nil
}

func (n noop) TrackWithBIData(data bi.Data) error {
	return nil
}

func track(roleID, event, fpid string, properties map[string]interface{}) error {
	data, err := bi.Track("", roleID, event, fpid, properties)
	if err != nil {
		return err
	}
	return tracker.TrackWithBIData(data)
}
