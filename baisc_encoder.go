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
	var head []zapcore.Field
	head = append(head, zap.String(LogId, xid.New().String())) // 日志规范要求必须有xid
	if glsFields := GetGlobalFields(); glsFields != nil {
		head = append(head, glsFields...)
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
