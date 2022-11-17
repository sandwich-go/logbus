package logbus

// NewFluentdCore creates a Core that writes logs to a WriteSyncer.
/*func NewFluentdCore(cfg zapcore.EncoderConfig, tags []string, enab zapcore.LevelEnabler) zapcore.Core {
	enc := NewFluentdEncoder(cfg)
	return &fluentdCore{
		LevelEnabler: enab,
		enc:          enc,
		tags:         tags,
	}
}

type fluentdCore struct {
	zapcore.LevelEnabler
	enc  *basics.fluentdEncoder
	tags []string
}

func (c *fluentdCore) With(fields []zapcore.Field) zapcore.Core {
	clone := c.clone()
	addFields(clone.enc.ObjectEncoder, fields)
	return clone
}

func addFields(enc *zapcore.MapObjectEncoder, fields []zapcore.Field) {
	for i := range fields {
		fields[i].AddTo(enc)
	}
}

func (c *fluentdCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	return ce.AddCore(ent, c)
}

func (c *fluentdCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	mapFields, err := c.enc.GetAllEntry(ent, fields)
	if err != nil {
		return err
	}
	for _, tag := range c.tags {
		err := fluentd.GetClient().Post(tag, mapFields)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *fluentdCore) Sync() error {
	return nil
}

func (c *fluentdCore) clone() *fluentdCore {
	return &fluentdCore{
		LevelEnabler: c.LevelEnabler,
		enc:          c.enc.Clone(),
		tags:         c.tags,
	}
}*/
