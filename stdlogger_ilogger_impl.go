package logbus

func (s *StdLogger) DebugWithChannel(c string, fields ...Field) {
	if s.fetch != nil {
		fields = append(fields, s.fetch()...)
	} else if Setting.FetchLogContext != nil {
		fields = append(fields, Setting.FetchLogContext()...)
	}
	s.z.Debug(c, fields...)
}

func (s *StdLogger) InfoWithChannel(c string, fields ...Field) {
	if s.fetch != nil {
		fields = append(fields, s.fetch()...)
	} else if Setting.FetchLogContext != nil {
		fields = append(fields, Setting.FetchLogContext()...)
	}
	s.z.Info(c, fields...)
}

func (s *StdLogger) WarnWithChannel(c string, fields ...Field) {
	if s.fetch != nil {
		fields = append(fields, s.fetch()...)
	} else if Setting.FetchLogContext != nil {
		fields = append(fields, Setting.FetchLogContext()...)
	}
	s.z.Warn(c, fields...)
}

func (s *StdLogger) ErrorWithChannel(c string, fields ...Field) {
	if s.fetch != nil {
		fields = append(fields, s.fetch()...)
	} else if Setting.FetchLogContext != nil {
		fields = append(fields, Setting.FetchLogContext()...)
	}
	s.z.Error(c, fields...)
}

func (s *StdLogger) DPanicWithChannel(c string, fields ...Field) {
	if s.fetch != nil {
		fields = append(fields, s.fetch()...)
	} else if Setting.FetchLogContext != nil {
		fields = append(fields, Setting.FetchLogContext()...)
	}
	s.z.DPanic(c, fields...)
}

func (s *StdLogger) PanicWithChannel(c string, fields ...Field) {
	if s.fetch != nil {
		fields = append(fields, s.fetch()...)
	} else if Setting.FetchLogContext != nil {
		fields = append(fields, Setting.FetchLogContext()...)
	}
	s.z.Panic(c, fields...)
}

func (s *StdLogger) FatalWithChannel(c string, fields ...Field) {
	if s.fetch != nil {
		fields = append(fields, s.fetch()...)
	} else if Setting.FetchLogContext != nil {
		fields = append(fields, Setting.FetchLogContext()...)
	}
	s.z.Fatal(c, fields...)
}
