package logbus

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

// GLogger 全局对象的类型定义
type GLogger struct {
	channelKey   string
	printAsError bool
	stdLogger    *StdLogger // 对外隐藏StdLogger
	depthLogger  sync.Map
}

/*type stringerArray struct {
	se glog.StringerEach
}

func (ss stringerArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	ss.se.Each(func(fs fmt.Stringer) {
		arr.AppendString(fs.String())
	})
	return nil
}

func (s *GLogger) convertField(msg string, fields ...glog.Field) (ret []zap.Field) {
	ret = append(ret, zap.String("glog-msg", msg))
	for _, f := range fields {
		if f.Type == glog.ArrayMarshalerType {
			switch f.ArrayType {
			case glog.StringsType:
				ret = append(ret, zap.Strings(f.Key, f.Interface.([]string)))
			case glog.StringersType:
				ret = append(ret, zap.Array(f.Key, stringerArray{se: f.Interface.(glog.StringerEach)}))
			case glog.IntsType:
				ret = append(ret, zap.Ints(f.Key, f.Interface.([]int)))
			case glog.Int64sType:
				ret = append(ret, zap.Int64s(f.Key, f.Interface.([]int64)))
			case glog.Int32sType:
				ret = append(ret, zap.Int32s(f.Key, f.Interface.([]int32)))
			case glog.UintsType:
				ret = append(ret, zap.Uints(f.Key, f.Interface.([]uint)))
			case glog.Uint64sType:
				ret = append(ret, zap.Uint64s(f.Key, f.Interface.([]uint64)))
			case glog.Uint32sType:
				ret = append(ret, zap.Uint32s(f.Key, f.Interface.([]uint32)))
			default:
				ret = append(ret, zap.Any(f.Key, f.Interface))
			}
			continue
		}
		if f.Type == glog.ObjectMarshalerType {
			if _, ok := f.Interface.(zapcore.ObjectMarshaler); !ok {
				f.Type = glog.ReflectType
			}
		}
		ret = append(ret, zap.Field{
			Key:       f.Key,
			Type:      zapcore.FieldType(f.Type),
			Integer:   f.Integer,
			String:    f.String,
			Interface: f.Interface,
		})
	}
	return ret
}*/

func (s *GLogger) printAsErr(fields ...zap.Field) bool {
	hasErr := false
	for _, v := range fields {
		if v.Type == zapcore.ErrorType {
			hasErr = true
			break
		}
	}
	return hasErr
}

func (s *GLogger) syncDepthLogger() {
	s.depthLogger.Range(func(key, value interface{}) bool {
		_ = value.(*StdLogger).Sync()
		return true
	})
}

func (s *GLogger) Debug(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	s.stdLogger.DebugWithChannel(s.channelKey, fields...)
}

func (s *GLogger) Info(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	s.stdLogger.InfoWithChannel(s.channelKey, fields...)
}

func (s *GLogger) Warn(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
		return
	}
	s.stdLogger.WarnWithChannel(s.channelKey, fields...)
}
func (s *GLogger) Error(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	s.stdLogger.ErrorWithChannel(s.channelKey, fields...)
}

func (s *GLogger) DPanic(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	s.stdLogger.DPanicWithChannel(s.channelKey, fields...)
}

func (s *GLogger) Panic(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	s.stdLogger.PanicWithChannel(s.channelKey, fields...)
}

func (s *GLogger) Fatal(msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	s.stdLogger.FatalWithChannel(s.channelKey, fields...)
}

func (s *GLogger) getDepthLogger(depth int) *zap.Logger {
	if lg, ok := s.depthLogger.Load(depth); ok {
		return lg.(*zap.Logger)
	}
	cloneLogger := s.stdLogger.WithOptions(zap.AddCallerSkip(depth))
	s.depthLogger.Store(depth, cloneLogger)
	return cloneLogger
}

func (s *GLogger) GDebugDepth(depth int, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.Error(s.channelKey, fields...)
		return
	}
	lg.Debug(s.channelKey, fields...)
}

func (s *GLogger) GInfoDepth(depth int, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.Error(s.channelKey, fields...)
		return
	}
	lg.Info(s.channelKey, fields...)
}
func (s *GLogger) GWarnDepth(depth int, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.Error(s.channelKey, fields...)
		return
	}
	lg.Warn(s.channelKey, fields...)
}
func (s *GLogger) GErrorDepth(depth int, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.Error(s.channelKey, fields...)
		return
	}
	lg.Error(s.channelKey, fields...)
}
func (s *GLogger) GFatalDepth(depth int, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	lg := s.getDepthLogger(depth)
	lg.Fatal(s.channelKey, fields...)
}

// WithChannel
func (s *GLogger) DebugWithChannel(c string, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	gStdLogger.DebugWithChannel(c, fields...)
}

func (s *GLogger) InfoWithChannel(c string, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	s.stdLogger.InfoWithChannel(c, fields...)
}

func (s *GLogger) WarnWithChannel(c string, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	s.stdLogger.WarnWithChannel(c, fields...)
}

func (s *GLogger) ErrorWithChannel(c string, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	s.stdLogger.ErrorWithChannel(c, fields...)
}

func (s *GLogger) DPanicWithChannel(c string, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	s.stdLogger.DPanicWithChannel(c, fields...)
}

func (s *GLogger) PanicWithChannel(c string, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	s.stdLogger.PanicWithChannel(c, fields...)
}

func (s *GLogger) FatalWithChannel(c string, msg string, fields ...zap.Field) {
	fields = append(fields, zap.String("glog-msg", msg))
	s.stdLogger.FatalWithChannel(c, fields...)
}
