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
		// encoded 引用 Zap 对象池中的 buffer，回调结束后会被 WithZapJSON 归还。
		// Zap2Json 返回的字节需要在函数返回后仍可使用，因此这里必须复制。
		bytes = append(bytes, encoded...)
	})
	return
}

// WithZapJSON 使用 Zap 对象池中的 buffer 调用 consume。encoded 仅在回调执行期间有效；
// consume 必须同步编码或复制，不能持有该切片，也不能交给异步消费者。
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
