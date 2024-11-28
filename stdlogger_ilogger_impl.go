package logbus

func (s *StdLogger) fields(fields []Field) []Field {
	f := s.fetch
	if f == nil {
		f = Setting.FetchLogContext
	}
	if f == nil {
		return fields
	}
	fRet := f()
	if fRet == nil {
		return fields
	}
	return append(fields, fRet...)
}

func (s *StdLogger) DebugWithChannel(c string, fields ...Field) {
	s.z.Debug(c, s.fields(fields)...)
}

func (s *StdLogger) InfoWithChannel(c string, fields ...Field) {
	s.z.Info(c, s.fields(fields)...)
}

func (s *StdLogger) WarnWithChannel(c string, fields ...Field) {
	s.z.Warn(c, s.fields(fields)...)
}

func (s *StdLogger) ErrorWithChannel(c string, fields ...Field) {
	s.z.Error(c, s.fields(fields)...)
}

func (s *StdLogger) DPanicWithChannel(c string, fields ...Field) {
	s.z.DPanic(c, s.fields(fields)...)
}

func (s *StdLogger) PanicWithChannel(c string, fields ...Field) {
	s.z.Panic(c, s.fields(fields)...)
}

func (s *StdLogger) FatalWithChannel(c string, fields ...Field) {
	s.z.Fatal(c, s.fields(fields)...)
}

func (s *StdLogger) TrackWithChannel(c string, fields ...Field) {
	s.z.Log(TrackLevel, c, s.fields(fields)...)
}
