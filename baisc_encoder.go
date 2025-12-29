package logbus

import (
	"github.com/rs/xid"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type glsEncoder struct {
	zapcore.Encoder
}

func (g *glsEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	glsFields := GetGlobalFields()
	var head = make([]zapcore.Field, 0, len(glsFields)+1+len(fields)) // 预分配空间，避免多次扩容
	head = append(head, zap.String(LogId, xid.New().String()))        // 日志规范要求必须有xid
	if len(glsFields) > 0 {
		head = append(head, glsFields...)
	}
	// 添加动态全局字段
	if dynamicFields := GetDynamicGlobalFields(); dynamicFields != nil {
		head = append(head, dynamicFields...)
	}
	fields = append(head, fields...)

	return g.Encoder.EncodeEntry(entry, fields)
}

func (g *glsEncoder) Clone() zapcore.Encoder {
	encoderClone := g.Encoder.Clone()
	return &glsEncoder{Encoder: encoderClone}
}

func newConsoleEncoder(config zapcore.EncoderConfig) (encoder zapcore.Encoder) {
	return &glsEncoder{Encoder: zapcore.NewConsoleEncoder(config)}
}

func newJSONEncoder(config zapcore.EncoderConfig) (encoder zapcore.Encoder) {
	return &glsEncoder{Encoder: zapcore.NewJSONEncoder(config)}
}
