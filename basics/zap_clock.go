package basics

import (
	"time"

	"go.uber.org/zap/zapcore"

	"bitbucket.org/funplus/sandwich/pkg/z"
)

type localClock struct{}

func (c localClock) Now() time.Time { return z.Now() }
func (c localClock) NewTicker(d time.Duration) *time.Ticker {
	return zapcore.DefaultClock.NewTicker(d)
}
