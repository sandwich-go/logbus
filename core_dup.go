package logbus

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sandwich-go/boost/xconv"
	"github.com/sandwich-go/boost/xos"
	"github.com/sandwich-go/boost/xslice"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	dupIgnoreFields = strings.Split(xos.EnvGetCaseInsensitive("logbus_core_dup_ignores"), ",")
	v := xos.EnvGetCaseInsensitive("logbus_core_dup")
	if v == "" {
		return
	}
	dupDuration = time.Duration(xconv.Int32(v)) * time.Millisecond
}

var dupDuration = time.Second
var dupIgnoreFields []string

type dupCore struct {
	zapcore.Core
	mu          sync.Mutex
	lastEntry   lastLogEntry
	repeatCount int
}

type lastLogEntry struct {
	fields   []zap.Field
	level    zapcore.Level
	message  string
	time     time.Time
	timeLast time.Time
}

func NewDupCore(core zapcore.Core) zapcore.Core {
	return &dupCore{Core: core}
}

func (c *dupCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *dupCore) isEqual(a, b lastLogEntry) bool {
	if a.level != b.level || a.message != b.message || len(a.fields) != len(b.fields) {
		return false
	}

	for i := range a.fields {
		af := a.fields[i]
		bf := b.fields[i]
		if xslice.StringsContain(dupIgnoreFields, af.Key) {
			continue
		}

		// Compare field basic properties
		if af.Key != bf.Key || af.Type != bf.Type {
			return false
		}
		// Compare values based on type
		switch af.Type {
		case zapcore.StringType:
			if af.String != bf.String {
				return false
			}
		case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
			if af.Integer != bf.Integer {
				return false
			}
		case zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type:
			if af.Integer != bf.Integer {
				return false
			}
		case zapcore.Float64Type, zapcore.Float32Type:
			if af.Integer != bf.Integer { // zap uses Integer to store floats internally
				return false
			}
		case zapcore.BoolType:
			if af.Integer != bf.Integer {
				return false
			}
		case zapcore.TimeType:
			if !time.Unix(0, af.Integer).Equal(time.Unix(0, bf.Integer)) {
				return false
			}
		case zapcore.DurationType:
			if af.Integer != bf.Integer {
				return false
			}
		default:
			// For complex types, don't compare
			return false
		}
	}
	return true
}

func (c *dupCore) Write(ent zapcore.Entry, fields []zap.Field) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	current := lastLogEntry{
		fields:  fields,
		level:   ent.Level,
		message: ent.Message,
		time:    ent.Time,
	}

	// Check if the current entry is the same as the last one AND within 1 second
	if c.isEqual(c.lastEntry, current) && ent.Time.Sub(c.lastEntry.timeLast) <= dupDuration {
		c.repeatCount++
		c.lastEntry.timeLast = ent.Time
		return nil
	}

	var err error
	if c.repeatCount > 0 {
		repeatEntry := c.lastEntry
		newEntry := zapcore.Entry{
			Level:   repeatEntry.level,
			Message: "[Repeated " + strconv.Itoa(c.repeatCount) + " times] " + repeatEntry.message,
			Time:    repeatEntry.time,
		}
		if err = c.Core.Write(newEntry, append(repeatEntry.fields,
			zap.Strings("logbus_dup_ignores", dupIgnoreFields),
			zap.Time("logbus_dup_end_time", c.lastEntry.timeLast))); err != nil {
			return err
		}
		c.repeatCount = 0
	}
	c.lastEntry = current
	c.lastEntry.timeLast = ent.Time // Initialize timeLast for the new entry
	return c.Core.Write(ent, fields)
}

func (c *dupCore) Sync() error {
	c.mu.Lock()
	if c.repeatCount > 0 {
		ent := zapcore.Entry{
			Level:   c.lastEntry.level,
			Message: "[Repeated " + strconv.Itoa(c.repeatCount) + " times] " + c.lastEntry.message,
			Time:    c.lastEntry.time,
		}
		_ = c.Core.Write(ent, append(c.lastEntry.fields, zap.Strings("logbus_dup_ignores", dupIgnoreFields), zap.Time("logbus_dup_end_time", c.lastEntry.timeLast)))
		c.repeatCount = 0
	}
	c.mu.Unlock()

	return c.Core.Sync()
}
