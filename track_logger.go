package logbus

import (
	"github.com/sandwich-go/boost/xerror"
	"go.uber.org/zap/zapcore"

	"github.com/sandwich-go/logbus/bigquery"
	"github.com/sandwich-go/logbus/thinkingdata"
)

var ErrIgnore = xerror.NewText("ignore track log")
var ErrTagNotImplement = xerror.NewText("tag not implement")

// Tracker 获取ITracker来打印thinkingData和bigQuery日志
func Tracker(tags ...string) ITracker {
	return &trackLogger{
		StdLogger: gStdLogger,
		tags:      tags,
	}
}

type trackLogger struct {
	*StdLogger
	tags []string
}

func (t *trackLogger) Track(fields ...Field) error {
	if ce := t.StdLogger.z.Check(TrackLevel, ""); ce == nil {
		// 检查逻辑前置，不做无用功
		return nil
	}
	for _, tag := range t.tags {
		switch tag {
		case THINKINGDATA:
			memoryEncoder := zapcore.NewMapObjectEncoder()
			for _, v := range fields {
				v.AddTo(memoryEncoder)
			}
			data, err := thinkingdata.ExtractEncoder(memoryEncoder)
			if err != nil {
				return err
			}
			t.StdLogger.PrintThingkingData(data)
		case BIGQUERY:
			tableName, bigFields, err := bigquery.ExtractEncoder(fields)
			if err != nil {
				return err
			}
			t.StdLogger.PrintBigQuery(tableName, bigFields...)
		default:
			return xerror.Wrap(ErrTagNotImplement, "tag of %s", tag)
		}
	}
	return nil
}

func (t *trackLogger) TrackWithTGAData(d thinkingdata.Data) error {
	if ce := t.StdLogger.z.Check(TrackLevel, ""); ce == nil {
		// 检查逻辑前置，不做无用功
		return nil
	}
	t.StdLogger.PrintThingkingData(d)
	return nil
}
