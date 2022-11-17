package logbus

import (
	"go.uber.org/zap/zapcore"
	"time"

	"github.com/sandwich-go/boost/z"
)

type localClock struct{}

func (c localClock) Now() time.Time { return z.Now() }
func (c localClock) NewTicker(d time.Duration) *time.Ticker {
	return zapcore.DefaultClock.NewTicker(d)
}
