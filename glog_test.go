package logbus

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"testing"
	"time"

	"github.com/sandwich-go/logbus/glog"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFieldsAndMsg(t *testing.T) {
	Debug("debug1", zap.Int("int", 123), zap.String("str", "foo"), zap.Error(nil))
	Info("info2", zap.Float64("float64", 8.99), zap.Bool("bool", true))
	Warn("", zap.Time("time", time.Now()), zap.Duration("duration", 1*time.Second), zap.Reflect("obj", map[string]string{"11": "22"}))
	Debug("error3", zap.Error(errors.New("rrrrrr")))
	Error("error4", zap.Binary("binary", []byte{'x'}), zap.Reflect("reflect", zap.Field{Key: "key", Integer: 122, Interface: "xxx"}), zap.Error(errors.New("eeeee")))
	Convey("add duplicated command flag should panic", t, func() {
		So(func() {
			Panic("panic5", zap.Any("any1", 1499), zap.Any("any2", zap.Field{Key: "key", Integer: 122, Interface: "xxx"}))
		}, ShouldPanic)
	})

	//Fatal("fatal6", Any("any1", 1499), Any("any2", []int{1, 2, 3, 4}))
}

func TestArrayFields(t *testing.T) {
	Debug("debug array", zap.Ints("ints", []int{1, 2, 3}), zap.Int64s("int64s", []int64{4, 5, 6}), zap.Float32s("float32s", []float32{1.2, 1.3}),
		zap.Strings("strings", []string{"a", "b", "c"}), zap.Times("k", []time.Time{time.Now()}),
		zap.Uintptrs("ptr", []uintptr{1, 2}), zap.Errors("errors", []error{errors.New("err1"), errors.New("err2")}))
}

type stringerObject struct {
	value string
}

func (s stringerObject) String() string {
	return s.value
}

type stringers []stringerObject

func (ss stringers) Each(handler func(stringer fmt.Stringer)) {
	for _, curr := range ss {
		handler(curr)
	}
}

func TestStringers(t *testing.T) {
	var sgs stringers
	sgs = append(sgs, stringerObject{value: "1a"}, stringerObject{value: "2a"}, stringerObject{value: "3a"})
	Debug("stringers", glog.Stringers("stringers", sgs))
}

func TestDepthLogger(t *testing.T) {
	DebugDepth(1, "debug1", zap.Int("int", 123), zap.String("str", "foo"), zap.Error(nil))
	InfoDepth(2, "info2", zap.Float64("float64", 8.99), zap.Bool("bool", true))
	WarnDepth(3, "", zap.Time("time", time.Now()), zap.Duration("duration", 1*time.Second) /*, zap.Object("obj", map[string]string{"1": "2"})*/)
	ErrorDepth(99, "error4", zap.Binary("binary", []byte{'x'}), zap.Reflect("reflect", zap.Field{Key: "key", Integer: 122, Interface: "xxx"}), zap.Error(errors.New("eeeee")))
	ErrorDepth(1, "error4", zap.Binary("binary", []byte{'x'}), zap.Reflect("reflect", zap.Field{Key: "key", Integer: 122, Interface: "xxx"}), zap.Error(errors.New("eeeee")))
	ErrorDepth(2, "error4", zap.Binary("binary", []byte{'x'}), zap.Reflect("reflect", zap.Field{Key: "key", Integer: 122, Interface: "xxx"}), zap.Error(errors.New("eeeee")))

	Error("error4", zap.Binary("binary", []byte{'x'}), zap.Reflect("reflect", zap.Field{Key: "key", Integer: 122, Interface: "xxx"}), zap.Error(errors.New("eeeee")))
}
