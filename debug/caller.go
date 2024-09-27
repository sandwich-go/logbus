package debug

import (
	"strings"

	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var bufferPool = buffer.NewPool()

// callerRelativePath2Mod 根据给定的call path全路径获取相对mod的路径，便于ide跳转
func callerRelativePath2Mod(ec zapcore.EntryCaller) string {
	if strings.HasPrefix(ec.File, moduleFilePathWithSeparator) {
		return strings.TrimPrefix(ec.File, moduleFilePathWithSeparator)
	}
	idx := strings.LastIndexByte(ec.File, '/')
	if idx == -1 {
		// 部分情况下fullPath会为package path / 相对路径
		if strings.HasPrefix(ec.File, packagePathWithSeparator) {
			return strings.TrimPrefix(ec.File, packagePathWithSeparator)
		}
	}
	// Find the penultimate separator.
	idx = strings.LastIndexByte(ec.File[:idx], '/')
	if idx == -1 {
		return ec.FullPath()
	}
	buf := bufferPool.Get()
	// Keep everything after the penultimate separator.
	buf.AppendString(ec.File[idx+1:])
	buf.AppendByte(':')
	buf.AppendInt(int64(ec.Line))
	caller := buf.String()
	buf.Free()
	return caller
}

// RelativePathCallerEncoder caller相对路径encoder
func RelativePathCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(func(ec zapcore.EntryCaller) string {
		buf := bufferPool.Get()
		buf.AppendString(callerRelativePath2Mod(ec))
		buf.AppendByte(':')
		buf.AppendInt(int64(ec.Line))
		caller := buf.String()
		buf.Free()
		return caller
	}(caller))
}
