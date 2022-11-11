package basics

import (
	"errors"
	"github.com/sandwich-go/logbus/config"
	"github.com/sandwich-go/logbus/globalfields"
	"time"

	"github.com/sandwich-go/logbus/thinkingdata"

	"github.com/rs/xid"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

func NewFluentdEncoder(cfg zapcore.EncoderConfig) *fluentdEncoder {
	return &fluentdEncoder{
		config:        &cfg,
		ObjectEncoder: zapcore.NewMapObjectEncoder(),
	}
}

type fluentdEncoder struct {
	config        *zapcore.EncoderConfig
	ObjectEncoder *zapcore.MapObjectEncoder
}

func (fe *fluentdEncoder) Clone() *fluentdEncoder {
	return NewFluentdEncoder(*fe.config)
}

func (fe *fluentdEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	return nil, errors.New("FluentdEncoder EncodeEntry not implement")
}

func (fe *fluentdEncoder) GetAllEntry(ent zapcore.Entry, fields []zapcore.Field) (map[string]interface{}, error) {
	c := fe.config
	var head []zapcore.Field
	head = append(head, zap.Time("date", ent.Time)) // 日志规范中要求必须包含"date"
	if c.LevelKey != "" {
		head = append(head, zap.String(c.LevelKey, ent.Level.String()))
	}
	if ent.LoggerName != "" && c.NameKey != "" {
		head = append(head, zap.String(c.NameKey, ent.LoggerName))
	}
	if ent.Caller.Defined {
		if c.CallerKey != "" {
			head = append(head, zap.String(c.CallerKey, ent.Caller.String()))
		}
		if c.FunctionKey != "" {
			head = append(head, zap.String(c.CallerKey, ent.Caller.Function))
		}
	}
	if c.MessageKey != "" {
		head = append(head, zap.String(c.MessageKey, ent.Message))
	}
	if ent.Stack != "" && c.StacktraceKey != "" {
		head = append(head, zap.String(c.StacktraceKey, ent.Stack))
	}
	if glsFields := globalfields.GetGlobalFields(); glsFields != nil {
		head = append(head, glsFields...)
	}
	head = append(head, zap.String(config.LogId, xid.New().String())) // 日志规范要求必须有xid
	objEncoder := zapcore.NewMapObjectEncoder()
	for _, v := range head {
		if v.Type == zapcore.TimeType {
			if v.Interface != nil {
				zap.String(v.Key, time.Unix(0, v.Integer).In(v.Interface.(*time.Location)).Format(thinkingdata.DATE_FORMAT)).AddTo(objEncoder)
			} else {
				zap.String(v.Key, time.Unix(0, v.Integer).Format(thinkingdata.DATE_FORMAT)).AddTo(objEncoder)
			}
			continue
		}
		v.AddTo(objEncoder)
	}
	for _, v := range fields {
		if v.Type == zapcore.TimeType {
			if v.Interface != nil {
				zap.String(v.Key, time.Unix(0, v.Integer).In(v.Interface.(*time.Location)).Format(thinkingdata.DATE_FORMAT)).AddTo(objEncoder)
			} else {
				zap.String(v.Key, time.Unix(0, v.Integer).Format(thinkingdata.DATE_FORMAT)).AddTo(objEncoder)
			}
			continue
		}
		v.AddTo(objEncoder)
	}
	return objEncoder.Fields, nil
}
