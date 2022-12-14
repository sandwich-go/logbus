package logbus

func (s *StdLogger) DebugWithChannel(c string, fields ...Field) {
	s.z.Debug(c, fields...)
}

func (s *StdLogger) InfoWithChannel(c string, fields ...Field) {
	s.z.Info(c, fields...)
}

func (s *StdLogger) WarnWithChannel(c string, fields ...Field) {
	s.z.Warn(c, fields...)
}

func (s *StdLogger) ErrorWithChannel(c string, fields ...Field) {
	s.z.Error(c, fields...)
}

func (s *StdLogger) DPanicWithChannel(c string, fields ...Field) {
	s.z.DPanic(c, fields...)
}

func (s *StdLogger) PanicWithChannel(c string, fields ...Field) {
	s.z.Panic(c, fields...)
}

func (s *StdLogger) FatalWithChannel(c string, fields ...Field) {
	s.z.Fatal(c, fields...)
}
