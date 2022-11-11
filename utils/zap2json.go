package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var entry = zapcore.Entry{}

var jsonEncoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
})

func Zap2Json(data []zap.Field) (bytes []byte, err error) {
	buffer, err := jsonEncoder.EncodeEntry(entry, data)
	if err != nil {
		return
	}
	bytes = buffer.Bytes()
	bytes = bytes[:len(bytes)-len(zapcore.DefaultLineEnding)]
	return
}
