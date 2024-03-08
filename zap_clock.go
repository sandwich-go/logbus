package logbus

import (
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/sandwich-go/boost/z"
)

type localClock struct{}

func (c localClock) Now() time.Time { return z.Now() }
func (c localClock) NewTicker(d time.Duration) *time.Ticker {
	return zapcore.DefaultClock.NewTicker(d)
}
