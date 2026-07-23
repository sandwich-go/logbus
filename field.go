package logbus

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/sandwich-go/boost/xerror"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zap.Field

const defaultMaskN = 3

func maskN(val string, n int) string {
	// 如果字符串长度小于等于 3*n，返回等长的掩码
	if len(val) <= 3*n {
		return strings.Repeat("*", len(val))
	}
	return val[:n] + strings.Repeat("*", len(val)-2*n) + val[len(val)-n:]
}

// MaskString 隐藏字符串中间的字符，只显示开头和结尾的字符，可指定收尾掩码长度，如果字符串长度小于等于 3*n[0]则使用等长掩码
func MaskString(key string, val string, n ...int) Field {
	keep := defaultMaskN
	if len(n) != 0 {
		keep = n[0]
	}
	return zap.String(key, maskN(val, keep))
}

// Password 使用等长掩码替代
func Password(key string, val string) Field {

	return zap.String(key, strings.Repeat("*", len(val)))
}
func P(key string, val string) Field { return Password(key, val) }

func String(key string, val string) Field {
	return zap.String(key, val)
}
func S(key string, val string) Field {
	return zap.String(key, val)
}

func Binary(key string, val []byte) Field {
	return zap.Binary(key, val)
}
func Bin(key string, val []byte) Field {
	return zap.Binary(key, val)
}

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}
func B(key string, val bool) Field {
	return zap.Bool(key, val)
}

func Float64(key string, val float64) Field {
	return zap.Float64(key, val)
}
func F64(key string, val float64) Field {
	return zap.Float64(key, val)
}

func Float32(key string, val float32) Field {
	return zap.Float32(key, val)
}
func F32(key string, val float32) Field {
	return zap.Float32(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}
func I(key string, val int) Field { return zap.Int(key, val) }

func Int64(key string, val int64) Field {
	return zap.Int64(key, val)
}
func I64(key string, val int64) Field {
	return zap.Int64(key, val)
}

func Int32(key string, val int32) Field {
	return zap.Int32(key, val)
}
func I32(key string, val int32) Field {
	return zap.Int32(key, val)
}

func Uint(key string, val uint) Field {
	return zap.Uint(key, val)
}
func U(key string, val uint) Field {
	return zap.Uint(key, val)
}

func Uint64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}
func U64(key string, val uint64) Field {
	return zap.Uint64(key, val)
}

func Uint32(key string, val uint32) Field {
	return zap.Uint32(key, val)
}
func U32(key string, val uint32) Field {
	return zap.Uint32(key, val)
}

func Uint16(key string, val uint16) Field {
	return zap.Uint16(key, val)
}
func U16(key string, val uint16) Field {
	return zap.Uint16(key, val)
}

func Uint8(key string, val uint8) Field {
	return zap.Uint8(key, val)
}
func U8(key string, val uint8) Field {
	return zap.Uint8(key, val)
}

func Reflect(key string, val interface{}) Field {
	return zap.Reflect(key, val)
}
func R(key string, val interface{}) Field {
	return zap.Reflect(key, val)
}

func Stringer(key string, val fmt.Stringer) Field {
	return zap.Stringer(key, val)
}

func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}
func T(key string, val time.Time) Field {
	return zap.Time(key, val)
}

func Stack(key string) Field {
	return zap.Stack(key)
}

func Duration(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}
func D(key string, val time.Duration) Field {
	return zap.Duration(key, val)
}

func Object(key string, val zapcore.ObjectMarshaler) Field {
	return zap.Object(key, val)
}
func O(key string, val zapcore.ObjectMarshaler) Field {
	return zap.Object(key, val)
}

func ErrorField(err error) Field {
	if err == nil {
		return zap.Skip()
	}
	if ignoreLogicalError.Load() && xerror.Logic(err) {
		return zap.String("error", err.Error()) // 逻辑错误不作为error输出
	}
	return zap.Error(err)
}

var ignoreLogicalError atomic.Bool

var E = ErrorField

func Any(key string, value interface{}) Field {
	return zap.Any(key, value)
}
func A(key string, value interface{}) Field {
	return zap.Any(key, value)
}
