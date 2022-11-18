package main

import (
	"time"

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

	// default channel, default tag
	logbus.Warn("", logbus.Int("money", 648))

	// reason: 打点的推荐方式：使用预定义的tags
	// Print tga log and big query log. New way
	_ = logbus.Tracker(logbus.THINKINGDATA, logbus.BIGQUERY).Track(logbus.String(thinkingdata.ACCOUNT, "111"), logbus.String(thinkingdata.TYPE, thinkingdata.USER_SET_ONCE),
		logbus.String("$user_id", "111"), logbus.Time("$optime", time.Now()), logbus.String(bigquery.TableNameKey, "oplog"),
		logbus.String("player_name", "zhang liu"), logbus.Int("level", 11), logbus.Bool("bool", true), logbus.Strings("strings", []string{"x", "y"}))

	// hook field: add playerid=gtwefasfwad for all logs below
	//logbus.appendGlobalFields(zap.String("playerid", "gtwefasfwad")) // Deprecated
	logbus.Warn("", logbus.Int("money", 648))
}
