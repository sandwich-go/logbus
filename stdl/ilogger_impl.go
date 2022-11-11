package stdl

import (
	"bitbucket.org/funplus/sandwich/pkg/logbus/basics"
	"go.uber.org/zap"
)

func (s *StdLogger) Debug(fields ...zap.Field) {
	s.z.Debug(basics.Setting.DefaultChannel, fields...)
}

func (s *StdLogger) DebugWithChannel(c string, fields ...zap.Field) {
	s.z.Debug(c, fields...)
}

func (s *StdLogger) Info(fields ...zap.Field) {
	s.z.Info(basics.Setting.DefaultChannel, fields...)
}

func (s *StdLogger) InfoWithChannel(c string, fields ...zap.Field) {
	s.z.Info(c, fields...)
}

func (s *StdLogger) Warn(fields ...zap.Field) {
	s.z.Warn(basics.Setting.DefaultChannel, fields...)
}

func (s *StdLogger) WarnWithChannel(c string, fields ...zap.Field) {
	s.z.Warn(c, fields...)
}

func (s *StdLogger) Error(fields ...zap.Field) {
	s.z.Error(basics.Setting.DefaultChannel, fields...)
}

func (s *StdLogger) ErrorWithChannel(c string, fields ...zap.Field) {
	s.z.Error(c, fields...)
}

func (s *StdLogger) DPanic(fields ...zap.Field) {
	s.z.DPanic(basics.Setting.DefaultChannel, fields...)
}

func (s *StdLogger) DPanicWithChannel(c string, fields ...zap.Field) {
	s.z.DPanic(c, fields...)
}

func (s *StdLogger) Panic(fields ...zap.Field) {
	s.z.Panic(basics.Setting.DefaultChannel, fields...)
}

func (s *StdLogger) PanicWithChannel(c string, fields ...zap.Field) {
	s.z.Panic(c, fields...)
}

func (s *StdLogger) Fatal(fields ...zap.Field) {
	s.z.Fatal(basics.Setting.DefaultChannel, fields...)
}

func (s *StdLogger) FatalWithChannel(c string, fields ...zap.Field) {
	s.z.Fatal(c, fields...)
}
