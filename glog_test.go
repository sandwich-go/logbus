package logbus

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"bitbucket.org/funplus/sandwich/pkg/logbus/glog"

	"bitbucket.org/funplus/sandwich/pkg/logbus/thinkingdata"
)

func TestGlogLogger(t *testing.T) {
	glog.Debug("id explosion", glog.Int("int", 123), glog.String("str", "foo"), glog.Err(nil))
	data := &thinkingdata.Data{
		AccountId:  "111",
		Type:       "track",
		Properties: map[string]interface{}{"zhangsan": 17},
	}
	glog.Info("tga", glog.Float64("float64", 8.99), glog.Bool("bool", true), glog.Object("tga", data), glog.Object("tga2", []int{1}))
	glog.Warn("duration", glog.Time("time", time.Now()), glog.Duration("duration", 1*time.Second))
	glog.Warn("duration", glog.Time("time", time.Now()), glog.Duration("duration", 1*time.Second), glog.Err(fmt.Errorf("fmt error")))
	glog.Error("", glog.Binary("binary", []byte{'x'}), glog.Reflect("reflect", &glog.Field{Key: "key", Integer: 122, Interface: "xxx"}), glog.Err(fmt.Errorf("fmt error")))
	//glog.Panic("", glog.Any("any1", 1499), glog.Any("any2", glog.Field{Key: "key", Integer: 122, Interface: "xxx"}))
	//glog.Fatal("", glog.Any("any1", 1499), glog.Any("any2", []int{1, 2, 3, 4}))

}

func TestGlogArrayLogger(t *testing.T) {
	SetGlobalLogger(GetStdLogger(), "siid", false)
	glog.Debug("debug array", glog.Strings("strings", []string{"a", "b", "c"}), glog.Uint64s("uints64", []uint64{1, 2, 3}), glog.Err(errors.New("an error")))
}
