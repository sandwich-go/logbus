package stdl

/*type StdLogger struct {
	z    *zap.Logger
	tags []string
}

func NewDefaultStdLogger(l *zap.Logger, t []string) *StdLogger {
	return &StdLogger{
		z:    l,
		tags: t,
	}
}

func GetStdLogger(tags ...string) *StdLogger {
	return getLogger(tags...)
}

func (s *StdLogger) CloneLogger(opts ...zap.Option) *StdLogger {
	cloneLogger := &StdLogger{
		z: s.WithOptions(opts...),
	}
	cloneLogger.tags = make([]string, len(s.tags))
	copy(cloneLogger.tags, s.tags)

	return cloneLogger
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
}*/

// L returns the zap Logger, // delete
/*func (s *StdLogger) L() *zap.Logger {
	return s.z
}*/
