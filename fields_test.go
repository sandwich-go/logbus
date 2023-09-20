package logbus

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/sandwich-go/logbus/thinkingdata"
)

func TestLoggerFields(t *testing.T) {
	Debug("id explosion", Int("int", 123), String("str", "foo"), ErrorField(nil))
	data := &thinkingdata.Data{
		AccountId:  "111",
		Type:       "track",
		Properties: map[string]interface{}{"zhangsan": 17},
	}
	Info("tga", Float64("float64", 8.99), Bool("bool", true), Object("tga", data), Ints("tga2", []int{1}))
	Warn("duration", Time("time", time.Now()), Duration("duration", 1*time.Second))
	Warn("duration", Time("time", time.Now()), Duration("duration", 1*time.Second), ErrorField(fmt.Errorf("fmt error")))
	Error("", Binary("binary", []byte{'x'}), Reflect("reflect", &Field{Key: "key", Integer: 122, Interface: "xxx"}), ErrorField(fmt.Errorf("fmt error")))
	//glog.Panic("", glog.Any("any1", 1499), glog.Any("any2", glog.Field{Key: "key", Integer: 122, Interface: "xxx"}))
	//glog.Fatal("", glog.Any("any1", 1499), glog.Any("any2", []int{1, 2, 3, 4}))
	Error("", Binary("binary", []byte{'x'}), Reflect("reflect", Getit()), ErrorField(fmt.Errorf("fmt error")))
}

func TestArrayFields(t *testing.T) {
	Debug("debug array", Ints("ints", []int{1, 2, 3}), Int64s("int64s", []int64{4, 5, 6}), Float32s("float32s", []float32{1.2, 1.3}),
		Strings("strings", []string{"a", "b", "c"}), Times("k", []time.Time{time.Now()}),
		Uintptrs("ptr", []uintptr{1, 2}), Errors("errors", []error{errors.New("err1"), errors.New("err2")}))
}

type it interface {
	Do()
}

type it2 interface {
	Do2()
}

func Getit() it {
	return &ts{t: &ts2{B: 1}}
}

type ts struct {
	t it2
}

func (ts) Do() {

}

type ts2 struct {
	B int
}

func (ts2) Do2() {

}
