package glog

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
}
