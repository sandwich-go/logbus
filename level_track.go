package logbus

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

const (
	TrackLevelColorGreen = 32
)

// TrackLevel track level支持自定义的数据级日志，如bi、thinking data
const TrackLevel = zapcore.DebugLevel - 1

// trackLevelName track 级使用的字符串名称
var trackLevelName = "track"
var trackLevelNameColored = trackLevelName

// SetTrackLevelName 更新track level使用的字符串名称
func SetTrackLevelName(name string) {
	trackLevelName = name
	trackLevelNameColored = fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(TrackLevelColorGreen), trackLevelName)
}

func init() {
	SetTrackLevelName(trackLevelName)
}

// NewTrackLevelEnabler 构造level过滤器，默认TrackLevel不受level过滤限制，允许覆盖实现自定义逻辑
var NewTrackLevelEnabler = func(level zapcore.Level) zapcore.LevelEnabler {
	return &trackLevelEnabler{Level: level}
}

type trackLevelEnabler struct {
	zapcore.Level
}

// Enabled track级日志不允许屏蔽
func (t *trackLevelEnabler) Enabled(level zapcore.Level) bool {
	if level == TrackLevel {
		return true
	}
	return t.Level.Enabled(level)
}
func trackLevelEncoderWithColor(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if level == TrackLevel {
		enc.AppendString(trackLevelNameColored)
		return
	}
	zapcore.LowercaseColorLevelEncoder(level, enc)
}
func trackLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	if level == TrackLevel {
		enc.AppendString(trackLevelName)
		return
	}
	zapcore.LowercaseLevelEncoder(level, enc)
}
