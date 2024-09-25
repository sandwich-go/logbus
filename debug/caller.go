package debug

import (
	"strings"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var bufferPool = buffer.NewPool()

// callerRelativePath2Mod 根据给定的call path全路径获取相对mod的路径，便于ide跳转
func callerRelativePath2Mod(fullPath string) string {
	return strings.TrimPrefix(fullPath, moduleFilePath)
}

// RelativePathCallerEncoder caller相对路径encoder
func RelativePathCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(func(ec zapcore.EntryCaller) string {
		buf := bufferPool.Get()
		buf.AppendString(callerRelativePath2Mod(ec.File))
		buf.AppendByte(':')
		buf.AppendInt(int64(ec.Line))
		caller := buf.String()
		buf.Free()
		return caller
	}(caller))
}
