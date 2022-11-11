package glog

import (
	"errors"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultLogger(t *testing.T) {
	Debug("debug1", Int("int", 123), String("str", "foo"), Err(nil))
	Info("info2", Float64("float64", 8.99), Bool("bool", true))
	Warn("", Time("time", time.Now()), Duration("duration", 1*time.Second), Object("obj", map[string]string{"1": "2"}))
	Debug("error3", Err(errors.New("rrrrrr")))
	Error("error4", Binary("binary", []byte{'x'}), Reflect("reflect", Field{Key: "key", Integer: 122, Interface: "xxx"}), Err(errors.New("eeeee")))
	Convey("add duplicated command flag should panic", t, func() {
		So(func() {
			Panic("panic5", Any("any1", 1499), Any("any2", Field{Key: "key", Integer: 122, Interface: "xxx"}))
		}, ShouldPanic)
	})

	//Fatal("fatal6", Any("any1", 1499), Any("any2", []int{1, 2, 3, 4}))
}

func TestArrayLogger(t *testing.T) {
	Debug("debug array", Ints("ints", []int{1, 2, 3}), Strings("strings", []string{"a", "b", "c"}))
}

func TestDepthLogger(t *testing.T) {
	DebugDepth(1, "debug1", Int("int", 123), String("str", "foo"), Err(nil))
	InfoDepth(2, "info2", Float64("float64", 8.99), Bool("bool", true))
	WarnDepth(3, "", Time("time", time.Now()), Duration("duration", 1*time.Second), Object("obj", map[string]string{"1": "2"}))
	ErrorDepth(99, "error4", Binary("binary", []byte{'x'}), Reflect("reflect", Field{Key: "key", Integer: 122, Interface: "xxx"}), Err(errors.New("eeeee")))
}
