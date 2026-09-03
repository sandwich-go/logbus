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
	err = WithZapJSON(data, func(encoded []byte) {
		bytes = append(bytes, encoded...)
	})
	return
}

// WithZapJSON calls consume with JSON backed by a Zap buffer. encoded is valid
// only for the duration of consume.
func WithZapJSON(data []zap.Field, consume func(encoded []byte)) error {
	buffer, err := jsonEncoder.EncodeEntry(entry, data)
	if err != nil {
		if buffer != nil {
			buffer.Free()
		}
		return err
	}
	defer buffer.Free()

	encoded := buffer.Bytes()
	encoded = encoded[:len(encoded)-len(zapcore.DefaultLineEnding)]
	consume(encoded)
	return nil
}
