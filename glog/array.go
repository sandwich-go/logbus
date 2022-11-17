package glog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

/*// Ints constructs a field that carries a slice of integers.
func Ints(key string, nums []int) Field {
	return Field{Key: key, Type: ArrayMarshalerType, Interface: nums, ArrayType: IntsType}
}

// Int64s constructs a field that carries a slice of integers.
func Int64s(key string, nums []int64) Field {
	return Field{Key: key, Type: ArrayMarshalerType, Interface: nums, ArrayType: Int64sType}
}

// Int32s constructs a field that carries a slice of integers.
func Int32s(key string, nums []int32) Field {
	return Field{Key: key, Type: ArrayMarshalerType, Interface: nums, ArrayType: Int32sType}
}

// Strings constructs a field that carries a slice of strings.
func Strings(key string, ss []string) Field {
	return Field{Key: key, Type: ArrayMarshalerType, Interface: ss, ArrayType: StringsType}
}*/

type stringerArray struct {
	se StringerEach
}

func (ss stringerArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	ss.se.Each(func(fs fmt.Stringer) {
		arr.AppendString(fs.String())
	})
	return nil
}

type StringerEach interface {
	Each(handler func(stringer fmt.Stringer))
}

func Stringers(key string, se StringerEach) zap.Field {
	return zap.Array(key, stringerArray{se: se})
	//return Field{Key: key, Type: ArrayMarshalerType, Interface: se, ArrayType: StringersType}
}

/*// Uints constructs a field that carries a slice of unsigned integers.
func Uints(key string, nums []uint) Field {
	return Field{Key: key, Type: ArrayMarshalerType, Interface: nums, ArrayType: UintsType}
}

// Uint64s constructs a field that carries a slice of unsigned integers.
func Uint64s(key string, nums []uint64) Field {
	return Field{Key: key, Type: ArrayMarshalerType, Interface: nums, ArrayType: Uint64sType}
}

// Uint32s constructs a field that carries a slice of unsigned integers.
func Uint32s(key string, nums []uint32) Field {
	return Field{Key: key, Type: ArrayMarshalerType, Interface: nums, ArrayType: Uint32sType}
}*/
