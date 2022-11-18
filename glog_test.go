package logbus

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/sandwich-go/logbus/glog"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFieldsAndMsg(t *testing.T) {
	Debug("debug1", Int("int", 123), String("str", "foo"), ErrorField(nil))
	Info("info2", Float64("float64", 8.99), Bool("bool", true))
	Warn("", Time("time", time.Now()), Duration("duration", 1*time.Second), Reflect("obj", map[string]string{"11": "22"}))
	Debug("error3", ErrorField(errors.New("rrrrrr")))
	Error("error4", Binary("binary", []byte{'x'}), Reflect("reflect", Field{Key: "key", Integer: 122, Interface: "xxx"}), ErrorField(errors.New("eeeee")))
	Convey("add duplicated command flag should panic", t, func() {
		So(func() {
			Panic("panic5", Any("any1", 1499), Any("any2", Field{Key: "key", Integer: 122, Interface: "xxx"}))
		}, ShouldPanic)
	})

	//Fatal("fatal6", Any("any1", 1499), Any("any2", []int{1, 2, 3, 4}))
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
	DebugDepth(1, "debug1", Int("int", 123), String("str", "foo"), ErrorField(nil))
	InfoDepth(2, "info2", Float64("float64", 8.99), Bool("bool", true))
	WarnDepth(3, "", Time("time", time.Now()), Duration("duration", 1*time.Second), Reflect("obj", map[string]string{"1": "2"}))
	ErrorDepth(99, "error4", Binary("binary", []byte{'x'}), Reflect("reflect", Field{Key: "key", Integer: 122, Interface: "xxx"}), ErrorField(errors.New("eeeee")))
	ErrorDepth(1, "error4", Binary("binary", []byte{'x'}), Reflect("reflect", Field{Key: "key", Integer: 122, Interface: "xxx"}), ErrorField(errors.New("eeeee")))
	ErrorDepth(2, "error4", Binary("binary", []byte{'x'}), Reflect("reflect", Field{Key: "key", Integer: 122, Interface: "xxx"}), ErrorField(errors.New("eeeee")))

	Error("error4", Binary("binary", []byte{'x'}), Reflect("reflect", Field{Key: "key", Integer: 122, Interface: "xxx"}), ErrorField(errors.New("eeeee")))
}
