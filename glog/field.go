package glog

import (
	"fmt"
	"math"
	"time"
)

/*
* definitions copied from zap
 */

type Field struct {
	Key  string
	Type FieldType
	ArrayType
	Integer   int64
	String    string
	Interface interface{}
}

type FieldType uint8

const (
	// UnknownType is the default field type. Attempting to add it to an encoder will panic.
	UnknownType FieldType = iota
	// ArrayMarshalerType indicates that the field carries an ArrayMarshaler.
	ArrayMarshalerType
	// ObjectMarshalerType indicates that the field carries an ObjectMarshaler.
	ObjectMarshalerType
	// BinaryType indicates that the field carries an opaque binary blob.
	BinaryType
	// BoolType indicates that the field carries a bool.
	BoolType
	// ByteStringType indicates that the field carries UTF-8 encoded bytes.
	ByteStringType
	// Complex128Type indicates that the field carries a complex128.
	Complex128Type
	// Complex64Type indicates that the field carries a complex128.
	Complex64Type
	// DurationType indicates that the field carries a time.Duration.
	DurationType
	// Float64Type indicates that the field carries a float64.
	Float64Type
	// Float32Type indicates that the field carries a float32.
	Float32Type
	// Int64Type indicates that the field carries an int64.
	Int64Type
	// Int32Type indicates that the field carries an int32.
	Int32Type
	// Int16Type indicates that the field carries an int16.
	Int16Type
	// Int8Type indicates that the field carries an int8.
	Int8Type
	// StringType indicates that the field carries a string.
	StringType
	// TimeType indicates that the field carries a time.Time that is
	// representable by a UnixNano() stored as an int64.
	TimeType
	// TimeFullType indicates that the field carries a time.Time stored as-is.
	TimeFullType
	// Uint64Type indicates that the field carries a uint64.
	Uint64Type
	// Uint32Type indicates that the field carries a uint32.
	Uint32Type
	// Uint16Type indicates that the field carries a uint16.
	Uint16Type
	// Uint8Type indicates that the field carries a uint8.
	Uint8Type
	// UintptrType indicates that the field carries a uintptr.
	UintptrType
	// ReflectType indicates that the field carries an interface{}, which should
	// be serialized using reflection.
	ReflectType
	// NamespaceType signals the beginning of an isolated namespace. All
	// subsequent fields should be added to the new namespace.
	NamespaceType
	// StringerType indicates that the field carries a fmt.Stringer.
	StringerType
	// ErrorType indicates that the field carries an error.
	ErrorType
	// SkipType indicates that the field is a no-op.
	SkipType
)

type ArrayType uint8

const (
	UnknownArrayType ArrayType = iota
	StringsType
	StringersType
	IntsType
	Int64sType
	Int32sType
	UintsType
	Uint64sType
	Uint32sType
)

func (f Field) AddTo(w *bytesFiledWriter) {
	if needCheckNil() && (f.Type == ObjectMarshalerType || f.Type == StringerType) && isNil(f.Interface) {
		return
	}

	var err error
	var ret string
	switch f.Type {
	case ObjectMarshalerType:
		ret, err = encodeDeferPanic("%s: %v", f.Key, f.Interface)
	case BinaryType:
		ret = fmt.Sprintf("%s: %v", f.Key, f.Interface.([]byte))
	case BoolType:
		ret = fmt.Sprintf("%s: %t", f.Key, f.Integer == 1)
	case DurationType:
		ret = fmt.Sprintf("%s: %s", f.Key, time.Duration(f.Integer).String())
	case Float64Type:
		ret = fmt.Sprintf("%s: %f", f.Key, math.Float64frombits(uint64(f.Integer)))
	case Float32Type:
		ret = fmt.Sprintf("%s: %f", f.Key, math.Float32frombits(uint32(f.Integer)))
	case Int64Type, Int32Type, Int16Type, Int8Type, Uint64Type, Uint32Type, Uint16Type, Uint8Type:
		ret = fmt.Sprintf("%s: %d", f.Key, f.Integer)
	case StringType:
		ret = fmt.Sprintf("%s: %s", f.Key, f.String)
	case StringerType:
		ret, err = encodeDeferPanic("%s: %s", f.Key, f.Interface.(fmt.Stringer).String())
	case TimeType:
		if f.Interface != nil {
			ret = fmt.Sprintf("%s: %v", f.Key, time.Unix(0, f.Integer).In(f.Interface.(*time.Location)))
		} else {
			// Fall back to UTC if location is nil.
			ret = fmt.Sprintf("%s: %v", f.Key, time.Unix(0, f.Integer))
		}
	case TimeFullType:
		ret = fmt.Sprintf("%s: %v", f.Key, f.Interface.(time.Time))
	case ReflectType:
		ret = fmt.Sprintf("%s: %v", f.Key, f.Interface)
	case SkipType:
		return
	case ErrorType:
		ret = fmt.Sprintf("%s: %s", f.Key, f.Interface.(error).Error())
	default:
		ret = fmt.Sprintf("%s: %v", f.Key, f.Interface)
	}
	w.WriteString(ret, f.Key)
	w.InternalError(f.Key, err)
}

func encodeDeferPanic(fmtString, key string, val interface{}) (ret string, err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("PANIC=%v", v)
		}
	}()
	ret = fmt.Sprintf(fmtString, key, val)
	return
}
