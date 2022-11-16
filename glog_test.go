package logbus

import (
	"errors"
	"github.com/sandwich-go/logbus/glog"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultLogger(t *testing.T) {
	Debug("debug1", glog.Int("int", 123), glog.String("str", "foo"), glog.Err(nil))
	Info("info2", glog.Float64("float64", 8.99), glog.Bool("bool", true))
	Warn("", glog.Time("time", time.Now()), glog.Duration("duration", 1*time.Second), glog.Object("obj", map[string]string{"1": "2"}))
	Debug("error3", glog.Err(errors.New("rrrrrr")))
	Error("error4", glog.Binary("binary", []byte{'x'}), glog.Reflect("reflect", glog.Field{Key: "key", Integer: 122, Interface: "xxx"}), glog.Err(errors.New("eeeee")))
	Convey("add duplicated command flag should panic", t, func() {
		So(func() {
			Panic("panic5", glog.Any("any1", 1499), glog.Any("any2", glog.Field{Key: "key", Integer: 122, Interface: "xxx"}))
		}, ShouldPanic)
	})

	//Fatal("fatal6", Any("any1", 1499), Any("any2", []int{1, 2, 3, 4}))
}

func TestArrayLogger(t *testing.T) {
	Debug("debug array", glog.Ints("ints", []int{1, 2, 3}), glog.Strings("strings", []string{"a", "b", "c"}))
}

func TestDepthLogger(t *testing.T) {
	DebugDepth(1, "debug1", glog.Int("int", 123), glog.String("str", "foo"), glog.Err(nil))
	InfoDepth(2, "info2", glog.Float64("float64", 8.99), glog.Bool("bool", true))
	WarnDepth(3, "", glog.Time("time", time.Now()), glog.Duration("duration", 1*time.Second), glog.Object("obj", map[string]string{"1": "2"}))
	ErrorDepth(99, "error4", glog.Binary("binary", []byte{'x'}), glog.Reflect("reflect", glog.Field{Key: "key", Integer: 122, Interface: "xxx"}), glog.Err(errors.New("eeeee")))
}
