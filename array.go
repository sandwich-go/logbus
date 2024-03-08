package logbus

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Array(key string, val zapcore.ArrayMarshaler) Field {
	return zap.Array(key, val)
}

func Bools(key string, bs []bool) Field {
	return zap.Bools(key, bs)
}

func ByteStrings(key string, bss [][]byte) Field {
	return zap.ByteStrings(key, bss)
}

func Durations(key string, ds []time.Duration) Field {
	return zap.Durations(key, ds)
}

func Float64s(key string, nums []float64) Field {
	return zap.Float64s(key, nums)
}

func Float32s(key string, nums []float32) Field {
	return zap.Float32s(key, nums)
}

func Ints(key string, nums []int) Field {
	return zap.Ints(key, nums)
}

func Int64s(key string, nums []int64) Field {
	return zap.Int64s(key, nums)
}

func Int32s(key string, nums []int32) Field {
	return zap.Int32s(key, nums)
}

func Int16s(key string, nums []int16) Field {
	return zap.Int16s(key, nums)
}

func Int8s(key string, nums []int8) Field {
	return zap.Int8s(key, nums)
}

func Strings(key string, ss []string) Field {
	return zap.Strings(key, ss)
}

func Times(key string, ts []time.Time) Field {
	return zap.Times(key, ts)
}

func Uints(key string, nums []uint) Field {
	return zap.Uints(key, nums)
}

func Uint64s(key string, nums []uint64) Field {
	return zap.Uint64s(key, nums)
}

func Uint32s(key string, nums []uint32) Field {
	return zap.Uint32s(key, nums)
}

func Uint16s(key string, nums []uint16) Field {
	return zap.Uint16s(key, nums)
}

func Uint8s(key string, nums []uint8) Field {
	return zap.Uint8s(key, nums)
}

func Errors(key string, errs []error) Field {
	return zap.Errors(key, errs)
}

func Uintptrs(key string, us []uintptr) Field {
	return zap.Uintptrs(key, us)
}
