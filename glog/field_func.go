package glog

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

var (
	_minTimeInt64 = time.Unix(0, math.MinInt64)
	_maxTimeInt64 = time.Unix(0, math.MaxInt64)
)

func Err(err error) Field {
	if err == nil {
		return Skip()
	}
	return Field{Key: "error", Type: ErrorType, Interface: err}
}

func Skip() Field {
	return Field{Type: SkipType}
}

// Binary constructs a field that carries an opaque binary blob.
//
// Binary data is serialized in an encoding-appropriate format. For example,
// zap's JSON encoder base64-encodes binary blobs. To log UTF-8 encoded text,
// use ByteString.
func Binary(key string, val []byte) Field {
	return Field{Key: key, Type: BinaryType, Interface: val}
}

// Bool constructs a field that carries a bool.
func Bool(key string, val bool) Field {
	var ival int64
	if val {
		ival = 1
	}
	return Field{Key: key, Type: BoolType, Integer: ival}
}

// Float64 constructs a field that carries a float64. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float64(key string, val float64) Field {
	return Field{Key: key, Type: Float64Type, Integer: int64(math.Float64bits(val))}
}

// Float32 constructs a field that carries a float32. The way the
// floating-point value is represented is encoder-dependent, so marshaling is
// necessarily lazy.
func Float32(key string, val float32) Field {
	return Field{Key: key, Type: Float32Type, Integer: int64(math.Float32bits(val))}
}

// Int constructs a field with the given key and value.
func Int(key string, val int) Field {
	return Int64(key, int64(val))
}

// Int64 constructs a field with the given key and value.
func Int64(key string, val int64) Field {
	return Field{Key: key, Type: Int64Type, Integer: val}
}

// Int32 constructs a field with the given key and value.
func Int32(key string, val int32) Field {
	return Field{Key: key, Type: Int32Type, Integer: int64(val)}
}

// Int16 constructs a field with the given key and value.
func Int16(key string, val int16) Field {
	return Field{Key: key, Type: Int16Type, Integer: int64(val)}
}

// Int8 constructs a field with the given key and value.
func Int8(key string, val int8) Field {
	return Field{Key: key, Type: Int8Type, Integer: int64(val)}
}

// String constructs a field with the given key and value.
func String(key string, val string) Field {
	return Field{Key: key, Type: StringType, String: val}
}

// Stringer constructs a field with the given key and the output of the value's
// String method. The Stringer's String method is called lazily.
func Stringer(key string, val fmt.Stringer) Field {
	return Field{Key: key, Type: StringerType, Interface: val}
}

// Uint constructs a field with the given key and value.
func Uint(key string, val uint) Field {
	return Uint64(key, uint64(val))
}

// Uint64 constructs a field with the given key and value.
func Uint64(key string, val uint64) Field {
	return Field{Key: key, Type: Uint64Type, Integer: int64(val)}
}

// Uint32 constructs a field with the given key and value.
func Uint32(key string, val uint32) Field {
	return Field{Key: key, Type: Uint32Type, Integer: int64(val)}
}

// Uint16 constructs a field with the given key and value.
func Uint16(key string, val uint16) Field {
	return Field{Key: key, Type: Uint16Type, Integer: int64(val)}
}

// Uint8 constructs a field with the given key and value.
func Uint8(key string, val uint8) Field {
	return Field{Key: key, Type: Uint8Type, Integer: int64(val)}
}

// Reflect constructs a field with the given key and an arbitrary object. It uses
// an encoding-appropriate, reflection-based function to lazily serialize nearly
// any object into the logging context, but it's relatively slow and
// allocation-heavy. Outside tests, Any is always a better choice.
//
// If encoding fails (e.g., trying to serialize a map[int]string to JSON), Reflect
// includes the error message in the final log output.
func Reflect(key string, val interface{}) Field {
	return Field{Key: key, Type: ReflectType, Interface: val}
}

// Time constructs a Field with the given key and value. The encoder
// controls how the time is serialized.
func Time(key string, val time.Time) Field {
	if val.Before(_minTimeInt64) || val.After(_maxTimeInt64) {
		return Field{Key: key, Type: TimeFullType, Interface: val}
	}
	return Field{Key: key, Type: TimeType, Integer: val.UnixNano(), Interface: val.Location()}
}

// Duration constructs a field with the given key and value. The encoder
// controls how the duration is serialized.
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: DurationType, Integer: int64(val)}
}

func isNil(val interface{}) bool {
	if val == nil {
		return true
	}
	r := reflect.ValueOf(val)
	switch r.Type().Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Slice:
		if r.IsNil() {
			return true
		}
	}
	return false
}

// Object constructs a field with the given key and ObjectMarshaler. It
// provides a flexible, but still type-safe and efficient, way to add map- or
// struct-like user-defined types to the logging context. The struct's
// MarshalLogObject method is called lazily.
func Object(key string, val interface{}) Field {
	return Field{Key: key, Type: ObjectMarshalerType, Interface: val}
}

// Any takes a key and an arbitrary value and chooses the best way to represent
// them as a field, falling back to a reflection-based approach only if
// necessary.
//
// Since byte/uint8 and rune/int32 are aliases, Any can't differentiate between
// them. To minimize surprises, []byte values are treated as binary blobs, byte
// values are treated as uint8, and runes are always treated as integers.
func Any(key string, value interface{}) Field {
	switch val := value.(type) {
	case bool:
		return Bool(key, val)
	case float64:
		return Float64(key, val)
	case float32:
		return Float32(key, val)
	case int:
		return Int(key, val)
	case int64:
		return Int64(key, val)
	case int32:
		return Int32(key, val)
	case int16:
		return Int16(key, val)
	case int8:
		return Int8(key, val)
	case string:
		return String(key, val)
	case uint:
		return Uint(key, val)
	case uint64:
		return Uint64(key, val)
	case uint32:
		return Uint32(key, val)
	case uint16:
		return Uint16(key, val)
	case uint8:
		return Uint8(key, val)
	case []byte:
		return Binary(key, val)
	case time.Time:
		return Time(key, val)
	case time.Duration:
		return Duration(key, val)
	default:
		return Reflect(key, val)
	}
}
