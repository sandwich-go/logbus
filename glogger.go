package logbus

import (
	"fmt"
	"github.com/sandwich-go/logbus/glog"
	"github.com/sandwich-go/logbus/stdl"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

//GLogger GLogger
type GLogger struct {
	channelKey   string
	printAsError bool
	stdLogger    *stdl.StdLogger // 对外隐藏StdLogger
	depthLogger  sync.Map
}

type stringerArray struct {
	se glog.StringerEach
}

func (ss stringerArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	ss.se.Each(func(fs fmt.Stringer) {
		arr.AppendString(fs.String())
	})
	return nil
}

func NewGLogger(ck string, sl *stdl.StdLogger, ae bool) *GLogger {
	return &GLogger{
		channelKey:   ck,
		printAsError: ae,
		stdLogger:    sl,
	}
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
}

func (s *GLogger) printAsErr(fields ...glog.Field) bool {
	hasErr := false
	for _, v := range fields {
		if v.Type == glog.ErrorType {
			hasErr = true
			break
		}
	}
	return hasErr
}

func (s *GLogger) GDebug(msg string, fields ...glog.Field) {
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, s.convertField(msg, fields...)...)
		return
	}
	s.stdLogger.DebugWithChannel(s.channelKey, s.convertField(msg, fields...)...)
}

func (s *GLogger) GInfo(msg string, fields ...glog.Field) {
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, s.convertField(msg, fields...)...)
		return
	}
	s.stdLogger.InfoWithChannel(s.channelKey, s.convertField(msg, fields...)...)
}

func (s *GLogger) GWarn(msg string, fields ...glog.Field) {
	if s.printAsError && s.printAsErr(fields...) {
		s.stdLogger.ErrorWithChannel(s.channelKey, s.convertField(msg, fields...)...)
		return
	}
	s.stdLogger.WarnWithChannel(s.channelKey, s.convertField(msg, fields...)...)
}
func (s *GLogger) GError(msg string, fields ...glog.Field) {
	s.stdLogger.ErrorWithChannel(s.channelKey, s.convertField(msg, fields...)...)
}

func (s *GLogger) GPanic(msg string, fields ...glog.Field) {
	s.stdLogger.PanicWithChannel(s.channelKey, s.convertField(msg, fields...)...)
}

func (s *GLogger) GFatal(msg string, fields ...glog.Field) {
	s.stdLogger.FatalWithChannel(s.channelKey, s.convertField(msg, fields...)...)
}

func (s *GLogger) getDepthLogger(depth int) *zap.Logger {
	if lg, ok := s.depthLogger.Load(depth); ok {
		return lg.(*zap.Logger)
	}
	cloneLogger := s.stdLogger.WithOptions(zap.AddCallerSkip(depth))
	s.depthLogger.Store(depth, cloneLogger)
	return cloneLogger
}

func (s *GLogger) GDebugDepth(depth int, msg string, fields ...glog.Field) {
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.Error(s.channelKey, s.convertField(msg, fields...)...)
		return
	}
	lg.Debug(s.channelKey, s.convertField(msg, fields...)...)
}

func (s *GLogger) GInfoDepth(depth int, msg string, fields ...glog.Field) {
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.Error(s.channelKey, s.convertField(msg, fields...)...)
		return
	}
	lg.Info(s.channelKey, s.convertField(msg, fields...)...)
}
func (s *GLogger) GWarnDepth(depth int, msg string, fields ...glog.Field) {
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.Error(s.channelKey, s.convertField(msg, fields...)...)
		return
	}
	lg.Warn(s.channelKey, s.convertField(msg, fields...)...)
}
func (s *GLogger) GErrorDepth(depth int, msg string, fields ...glog.Field) {
	lg := s.getDepthLogger(depth)
	if s.printAsError && s.printAsErr(fields...) {
		lg.Error(s.channelKey, s.convertField(msg, fields...)...)
		return
	}
	lg.Error(s.channelKey, s.convertField(msg, fields...)...)
}
func (s *GLogger) GFatalDepth(depth int, msg string, fields ...glog.Field) {
	lg := s.getDepthLogger(depth)
	lg.Fatal(s.channelKey, s.convertField(msg, fields...)...)
}
