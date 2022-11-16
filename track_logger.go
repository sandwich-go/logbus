package logbus

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sandwich-go/logbus/bigquery"
	"github.com/sandwich-go/logbus/thinkingdata"
)

// Tracker 获取ITracker来打印thinkingData和bigQuery日志
func Tracker(tags ...string) ITracker {
	//return getStdLogger(tags...)
	return &TrackLogger{
		StdLogger: gStdLogger,
		tags:      tags,
	}
}

type TrackLogger struct {
	*StdLogger
	tags []string
}

func (t *TrackLogger) Track(fields ...zap.Field) error {
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
			return fmt.Errorf("tag %s not implement", tag)
		}
	}
	return nil
}
