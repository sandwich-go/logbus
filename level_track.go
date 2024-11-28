package logbus

import "go.uber.org/zap/zapcore"

// TrackLevel track level支持自定义的数据级日志，如bi、thinking data
const TrackLevel = zapcore.DebugLevel - 1

// TrackLevelName track 级使用的字符串名称
var TrackLevelName = "track"

// NewTrackLevelEnabler 构造level过滤器，默认TrackLevel不受level过滤限制，允许覆盖实现自定义逻辑
var NewTrackLevelEnabler = func(level zapcore.Level) zapcore.LevelEnabler {
	return &trackLevelEnabler{Level: level}
}

type trackLevelEnabler struct {
	zapcore.Level
}

func (t *trackLevelEnabler) Enabled(level zapcore.Level) bool {
	if level == TrackLevel {
		return true
	}
	return t.Level.Enabled(level)
}

func trackLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if level == TrackLevel {
		enc.AppendString(TrackLevelName)
		return
	}
	zapcore.LowercaseLevelEncoder(level, enc)
}
