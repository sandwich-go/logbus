package logbus

import "go.uber.org/zap"

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
	newFields := make([]Field, 0, len(fields)+len(fRet))
	newFields = append(newFields, fRet...)
	return append(newFields, fields...)
}

func (s *StdLogger) DebugWithChannel(c string, fields ...Field) {
	// 级别关闭时直接返回：s.fields 会调 FetchLogContext 钩子并分配新切片，
	// 而它在 s.z.Debug 之前就求值了，zap 内部的级别判定拦不住这部分开销。
	if !s.Enabled(zap.DebugLevel) {
		return
	}
	s.z.Debug(c, s.fields(fields)...)
}

func (s *StdLogger) InfoWithChannel(c string, fields ...Field) {
	if !s.Enabled(zap.InfoLevel) {
		return
	}
	s.z.Info(c, s.fields(fields)...)
}

func (s *StdLogger) WarnWithChannel(c string, fields ...Field) {
	if !s.Enabled(zap.WarnLevel) {
		return
	}
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
	if !s.Enabled(TrackLevel) {
		return
	}
	s.z.Log(TrackLevel, c, s.fields(fields)...)
}
