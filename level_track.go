package logbus

import "go.uber.org/zap/zapcore"

// TrackLevel track level支持自定义的数据级日志，如bi、thinking data
const TrackLevel = zapcore.DebugLevel - 1

// TraceLevelName track 级使用的字符串名称
var TraceLevelName = "track"

// NewTraceLevelEnabler 构造level过滤器，默认TrackLevel不受level过滤限制，允许覆盖实现自定义逻辑
var NewTraceLevelEnabler = func(level zapcore.Level) zapcore.LevelEnabler {
	return &traceLevelEnabler{Level: level}
}

type traceLevelEnabler struct {
	zapcore.Level
}

func (t *traceLevelEnabler) Enabled(level zapcore.Level) bool {
	if level == TrackLevel {
		return true
	}
	return t.Level.Enabled(level)
}

func traceLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if level == TrackLevel {
		enc.AppendString(TraceLevelName)
		return
	}
	zapcore.LowercaseLevelEncoder(level, enc)
}
