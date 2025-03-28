package logbus

import (
	"go.uber.org/zap"
)

type StdLogger struct {
	fetch FetchLogContext
	z     *zap.Logger
}

func (s *StdLogger) WithOptions(opts ...zap.Option) *zap.Logger {
	return s.z.WithOptions(opts...)
}

func (s *StdLogger) SetZLogger(zl *zap.Logger) {
	if zl != nil {
		s.z = zl
	}
}

func (s *StdLogger) Sync() error {
	return s.z.Sync()
}

func (s *StdLogger) getZapLogger() *zap.Logger {
	return s.z
}

func newStdLogger(z *zap.Logger, fetch FetchLogContext) *StdLogger {
	s := &StdLogger{
		fetch: fetch,
		z:     z,
	}
	return s
}
