package logbus

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"testing"
	"time"

	"github.com/sandwich-go/logbus/glog"

	"github.com/sandwich-go/logbus/thinkingdata"
)

func TestGlogLogger(t *testing.T) {
	Debug("id explosion", zap.Int("int", 123), zap.String("str", "foo"), zap.Error(nil))
	data := &thinkingdata.Data{
		AccountId:  "111",
		Type:       "track",
		Properties: map[string]interface{}{"zhangsan": 17},
	}
	Info("tga", zap.Float64("float64", 8.99), zap.Bool("bool", true), zap.Object("tga", data), zap.Ints("tga2", []int{1}))
	Warn("duration", zap.Time("time", time.Now()), zap.Duration("duration", 1*time.Second))
	Warn("duration", zap.Time("time", time.Now()), zap.Duration("duration", 1*time.Second), zap.Error(fmt.Errorf("fmt error")))
	Error("", zap.Binary("binary", []byte{'x'}), zap.Reflect("reflect", &glog.Field{Key: "key", Integer: 122, Interface: "xxx"}), zap.Error(fmt.Errorf("fmt error")))
	//glog.Panic("", glog.Any("any1", 1499), glog.Any("any2", glog.Field{Key: "key", Integer: 122, Interface: "xxx"}))
	//glog.Fatal("", glog.Any("any1", 1499), glog.Any("any2", []int{1, 2, 3, 4}))

}

func TestGlogArrayLogger(t *testing.T) {
	SetGlobalGLogger(nil, "siid", false)
	Debug("debug array", zap.Strings("strings", []string{"a", "b", "c"}), zap.Uint64s("uints64", []uint64{1, 2, 3}), zap.Error(errors.New("an error")))
}
