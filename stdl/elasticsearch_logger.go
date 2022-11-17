package stdl

/*func toFields(msg string, uid uint64, fields []zap.Field) []zap.Field {
	data, _ := utils.Zap2Json(fields)
	res := []zap.Field{
		zap.String("msg", msg),
		zap.Uint64("uid", uid),
		zap.ByteString("content", data),
	}
	return res
}

// Msg 用于输出到不支持flattend类型的elasticserch(例如aws) https://upsource.diandian.info:3003/document/module/log/logbus_faq#tip8
func (s *StdLogger) Msg(level zapcore.Level, msg string, uid uint64, fields ...zap.Field) {
	fields = toFields(msg, uid, fields)
	switch level {
	case zapcore.DebugLevel:
		s.z.Debug(basics.Setting.DefaultChannel, fields...)
	case zapcore.InfoLevel:
		s.z.Info(basics.Setting.DefaultChannel, fields...)
	case zapcore.WarnLevel:
		s.z.Warn(basics.Setting.DefaultChannel, fields...)
	case zapcore.ErrorLevel:
		s.z.Error(basics.Setting.DefaultChannel, fields...)
	case zapcore.DPanicLevel:
		s.z.DPanic(basics.Setting.DefaultChannel, fields...)
	case zapcore.PanicLevel:
		s.z.Panic(basics.Setting.DefaultChannel, fields...)
	case zapcore.FatalLevel:
		s.z.Fatal(basics.Setting.DefaultChannel, fields...)
	}
}*/
