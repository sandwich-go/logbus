package basics

/*type glsEncoder struct {
	zapcore.Encoder
}

func newGlsEncoder(encoder zapcore.Encoder) zapcore.Encoder {
	return &glsEncoder{Encoder: encoder}
}

func (g *glsEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	var head []zapcore.Field
	head = append(head, zap.String(config.LogId, xid.New().String())) // 日志规范要求必须有xid
	if glsFields := globalfields.GetGlobalFields(); glsFields != nil {
		head = append(head, glsFields...)
	}
	fields = append(head, fields...)
	return g.Encoder.EncodeEntry(entry, fields)
}

func (g *glsEncoder) Clone() zapcore.Encoder {
	encoderClone := g.Encoder.Clone()
	return &glsEncoder{Encoder: encoderClone}
}

func NewConsoleEncoder(config zapcore.EncoderConfig) (encoder zapcore.Encoder) {
	return newGlsEncoder(zapcore.NewConsoleEncoder(config))
}

func NewJSONEncoder(config zapcore.EncoderConfig) (encoder zapcore.Encoder) {
	return newGlsEncoder(zapcore.NewJSONEncoder(config))
}*/
