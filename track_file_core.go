package logbus

import (
	"encoding/base64"
	"fmt"

	"github.com/sandwich-go/logbus/bigquery"
	"go.uber.org/zap/zapcore"
)

type trackFileCore struct {
	enabler zapcore.LevelEnabler
	writer  *trackFileWriteSyncerProxy
	fields  []zapcore.Field
}

func newTrackFileCore(writer *trackFileWriteSyncerProxy, enabler zapcore.LevelEnabler) zapcore.Core {
	return &trackFileCore{
		enabler: enabler,
		writer:  writer,
	}
}

func (c *trackFileCore) Enabled(level zapcore.Level) bool {
	return c.enabler.Enabled(level)
}

func (c *trackFileCore) With(fields []zapcore.Field) zapcore.Core {
	clone := *c
	clone.fields = make([]zapcore.Field, 0, len(c.fields)+len(fields))
	clone.fields = append(clone.fields, c.fields...)
	clone.fields = append(clone.fields, fields...)
	return &clone
}

func (c *trackFileCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *trackFileCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	return c.writer.writeTrack(ent.Message, c.fields, fields)
}

func (c *trackFileCore) Sync() error {
	return c.writer.Sync()
}

func extractTrackChannelAndMsg(channel string, contextFields, fields []zapcore.Field) (string, []byte, bool) {
	if channel == "" {
		return "", nil, false
	}

	var (
		msg       []byte
		hasMsg    bool
		tableName string
	)
	scanTrackFields(contextFields, &msg, &hasMsg, &tableName)
	scanTrackFields(fields, &msg, &hasMsg, &tableName)
	if !hasMsg {
		return "", nil, false
	}
	if channel == BIGQUERY && tableName != "" {
		channel += "_" + tableName
	}
	return channel, msg, true
}

func scanTrackFields(fields []zapcore.Field, msg *[]byte, hasMsg *bool, tableName *string) {
	for i := range fields {
		field := fields[i]
		switch field.Key {
		case MsgBody:
			if data, ok := trackFieldBytes(field); ok {
				*msg = data
				*hasMsg = true
			}
		case bigquery.TableNameKey:
			if v, ok := trackFieldString(field); ok && v != "" {
				*tableName = v
			}
		}
	}
}

func trackFieldBytes(field zapcore.Field) ([]byte, bool) {
	switch field.Type {
	case zapcore.StringType:
		return []byte(field.String), true
	case zapcore.ByteStringType:
		if b, ok := field.Interface.([]byte); ok {
			return b, true
		}
	case zapcore.BinaryType:
		if b, ok := field.Interface.([]byte); ok {
			return []byte(base64.StdEncoding.EncodeToString(b)), true
		}
	case zapcore.StringerType:
		if field.Interface != nil {
			return []byte(fmt.Sprint(field.Interface)), true
		}
	default:
		return trackFieldJSONBytes(field)
	}
	return nil, false
}

func trackFieldJSONBytes(field zapcore.Field) ([]byte, bool) {
	defer func() {
		_ = recover()
	}()

	enc := zapcore.NewMapObjectEncoder()
	field.AddTo(enc)
	v, ok := enc.Fields[field.Key]
	if !ok {
		return nil, false
	}
	if s, ok := v.(string); ok {
		return []byte(s), true
	}
	b, err := jsonLib.Marshal(v)
	if err != nil {
		return nil, false
	}
	return b, true
}

func trackFieldString(field zapcore.Field) (string, bool) {
	switch field.Type {
	case zapcore.StringType:
		return field.String, true
	case zapcore.ByteStringType:
		if b, ok := field.Interface.([]byte); ok {
			return string(b), true
		}
	case zapcore.StringerType:
		if field.Interface != nil {
			return fmt.Sprint(field.Interface), true
		}
	}
	return "", false
}
