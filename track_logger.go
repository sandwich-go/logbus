package logbus

import (
	"github.com/sandwich-go/boost/xerror"
	"github.com/sandwich-go/boost/xpanic"
	"github.com/sandwich-go/boost/xslice"
	"go.uber.org/zap/zapcore"

	"github.com/sandwich-go/logbus/bi"
	"github.com/sandwich-go/logbus/bigquery"
	"github.com/sandwich-go/logbus/thinkingdata"
)

var ErrIgnore = xerror.NewText("ignore track log")
var ErrTagNotImplement = xerror.NewText("tag not implement")
var ErrBiAppIDEmpty = xerror.NewText("BiAppID can not be empty when using BI tag")

// Tracker 获取ITracker来打印thinkingData、bigQuery、BI日志
//
// Deprecated: 推荐使用 NewTracker
func Tracker(tags ...string) ITracker {
	return NewTracker(WithTags(tags...))
}

func NewTracker(opts ...TrackLoggerConfOption) ITracker {
	cc := NewTrackLoggerConf(opts...)
	xpanic.WhenTrue(len(cc.tags) == 0, "tags can not be empty")
	xpanic.WhenTrue(xslice.StringsContain(cc.tags, BI) && cc.BiAppID == "", "BiAppID can not be empty when using BI tag")
	return &trackLogger{
		cc: cc,
	}
}

type trackLogger struct {
	cc *TrackLoggerConf
}

func (t *trackLogger) Track(fields ...Field) error {
	logger := gStdLogger.Load()
	if ce := logger.z.Check(TrackLevel, ""); ce == nil {
		// 检查逻辑前置，不做无用功
		return nil
	}
	for _, tag := range t.cc.tags {
		switch tag {
		case THINKINGDATA:
			data, err := thinkingdata.ExtractFields(fields)
			if err != nil {
				return err
			}
			logger.PrintThingkingData(data)
		case BIGQUERY:
			tableName, bigFields, err := bigquery.ExtractEncoder(fields)
			if err != nil {
				return err
			}
			logger.PrintBigQuery(tableName, bigFields...)
		case BI:
			memoryEncoder := zapcore.NewMapObjectEncoder()
			for _, v := range fields {
				v.AddTo(memoryEncoder)
			}
			data, err := bi.ExtractEncoder(memoryEncoder)
			if err != nil {
				return err
			}
			_ = t.TrackWithBIData(data)
		default:
			return xerror.Wrap(ErrTagNotImplement, "tag of %s", tag)
		}
	}
	return nil
}

func (t *trackLogger) TrackWithTGAData(d thinkingdata.Data) error {
	logger := gStdLogger.Load()
	if ce := logger.z.Check(TrackLevel, ""); ce == nil {
		// 检查逻辑前置，不做无用功
		return nil
	}
	logger.PrintThingkingData(d)
	return nil
}

func (t *trackLogger) TrackWithBIData(d bi.Data) error {
	logger := gStdLogger.Load()
	if ce := logger.z.Check(TrackLevel, ""); ce == nil {
		return nil
	}
	if d.AppID == "" {
		d.AppID = t.cc.BiAppID
	}
	if d.AppID == "" {
		return ErrBiAppIDEmpty
	}
	logger.PrintBIData(d)
	return nil
}
