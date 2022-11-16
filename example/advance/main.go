package main

import (
	"time"

	"go.uber.org/zap"

	"github.com/sandwich-go/logbus"
	"github.com/sandwich-go/logbus/bigquery"
	"github.com/sandwich-go/logbus/thinkingdata"
)

func main() {
	// close logger before exit
	defer logbus.Close()

	// Init with conf
	logbus.Init(logbus.NewConf(
		logbus.WithDev(false),
		logbus.WithMonitorOutput(logbus.Prometheus),
		logbus.WithDefaultChannel("Game"),
		logbus.WithDefaultTag("Advance")),
	)

	// reason: 两种同等作用的写法，保留一个
	// tag=logbus.DefaultTag  dd_meta_channel=setting.DefaultChannel
	logbus.Warn("", zap.Int("money", 648))

	// reason: 打点的推荐方式：使用预定义的tags
	// Print tga log and big query log. New way
	_ = logbus.Tracker(logbus.THINKINGDATA, logbus.BIGQUERY).Track(zap.String(thinkingdata.ACCOUNT, "111"), zap.String(thinkingdata.TYPE, thinkingdata.USER_SET_ONCE),
		zap.String("$user_id", "111"), zap.Time("$optime", time.Now()), zap.String(bigquery.TableNameKey, "oplog"),
		zap.String("player_name", "zhang liu"), zap.Int("level", 11), zap.Bool("bool", true), zap.Strings("strings", []string{"x", "y"}))

	// hook field: add playerid=gtwefasfwad for all logs below
	//logbus.appendGlobalFields(zap.String("playerid", "gtwefasfwad")) // Deprecated
	logbus.Warn("", zap.Int("money", 648))
}
