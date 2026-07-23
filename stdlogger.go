package logbus

import (
	"go.uber.org/zap"
)

type StdLogger struct {
	fetch FetchLogContext
	z     *zap.Logger
	cfg   *configSnapshot
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

func newStdLogger(z *zap.Logger, config *configSnapshot, fetch FetchLogContext) *StdLogger {
	if config == nil {
		config = currentConfig()
	}
	if fetch == nil {
		fetch = config.FetchLogContext
	}
	s := &StdLogger{
		fetch: fetch,
		z:     z,
		cfg:   config,
	}
	return s
}

func (s *StdLogger) config() *configSnapshot {
	if s != nil && s.cfg != nil {
		return s.cfg
	}
	return currentConfig()
}
