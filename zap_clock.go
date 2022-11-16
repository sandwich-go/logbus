package logbus

import (
	"time"

	"github.com/sandwich-go/boost/xz"
)

type localClock struct{}

func (c localClock) Now() time.Time { return xz.Now() }
func (c localClock) NewTicker(d time.Duration) *time.Ticker {
	return xz.NewTicker(d)
}
