package main

import (
	"bitbucket.org/funplus/sandwich/pkg/logbus"
	"log"
	"os"

	"github.com/myzhan/boomer"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"bitbucket.org/funplus/sandwich/pkg/logbus/thinkingdata"
)

var count atomic.Uint64

func foo() {
	logbus.GetStdLogger().PrintThingkingData(data)
	boomer.RecordSuccess("http", "foo", 1, 10)
	count.Add(1)
}

var properties map[string]interface{}
var data thinkingdata.Data

func main() {
	defer func() {
		logbus.Warn(zap.Uint64("count", count.Load()))
	}()
	log.SetOutput(os.Stderr)
	task1 := &boomer.Task{
		Name: "foo",
		// The weight is used to distribute goroutines over multiple tasks.
		Weight: 10,
		Fn:     foo,
	}
	properties = map[string]interface{}{"#ip": "10.0.0.2", "player_name": "zhang si", "level": 8, "data": []string{"x", "y"}, "rc": 100, "install_ts": "2019-01-31 04:07:23",
		"is_paiduser": false, "idfa": "2965F487-588B-43E9-8AF0-1D6182531AEC", "idfv": "EACB27A3-EC86-4019-BE1D-2851117354B2", "android_id": "", "gl_version": "2.0", "app_version": "6.6.100",
		"version_code": "64277", "device": "iPad11", "os": "ffs.global.iOS", "coins": 1001641, "tc": 126121, "continuous_day": 3, "vip_level": 5}
	data, _ = thinkingdata.User("111", "", thinkingdata.USER_SET_ONCE, properties)
	logbus.GetStdLogger().PrintThingkingData(data) // 730 bytes
	boomer.Run(task1)
}
