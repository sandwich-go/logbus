package logbus

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

const (
	TrackLevelColorGreen uint8 = 32
)

// TrackLevel track level支持自定义的数据级日志，如bi、thinking data
const TrackLevel = zapcore.DebugLevel - 1

// trackLevelName track 级使用的字符串名称
var trackLevelName = "track"
var trackLevelNameColored = trackLevelName

// setTrackLevelName 更新track level使用的字符串名称
func setTrackLevelName(name string) {
	trackLevelName = name
	trackLevelNameColored = fmt.Sprintf("\x1b[%dm%s\x1b[0m", TrackLevelColorGreen, trackLevelName)
}

func init() {
	setTrackLevelName(trackLevelName)
}

// newTrackLevelEnabler 构造level过滤器，默认TrackLevel不受level过滤限制，允许覆盖实现自定义逻辑
var newTrackLevelEnabler = func(level zapcore.Level, enableTrack bool) zapcore.LevelEnabler {
	return &trackLevelEnabler{Level: level, enableTrack: enableTrack}
}

type trackLevelEnabler struct {
	zapcore.Level
	enableTrack bool
}

// Enabled track级日志不允许屏蔽
func (t *trackLevelEnabler) Enabled(level zapcore.Level) bool {
	if t.enableTrack && level == TrackLevel {
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
