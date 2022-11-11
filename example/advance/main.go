package main

import (
	"bitbucket.org/funplus/sandwich/pkg/logbus/config"
	"bitbucket.org/funplus/sandwich/pkg/logbus/globalfields"
	"time"

	"go.uber.org/zap"

	"bitbucket.org/funplus/sandwich/pkg/logbus"
	"bitbucket.org/funplus/sandwich/pkg/logbus/bigquery"
	"bitbucket.org/funplus/sandwich/pkg/logbus/thinkingdata"
)

func main() {
	// close logger before exit
	defer logbus.Close()

	// Init with conf
	logbus.Init(config.NewConf(config.WithOutputStdout(true), config.WithDev(false), config.WithMonitorOutput(logbus.Prometheus), config.WithDefaultChannel("Game")))

	// reason: 两种同等作用的写法，保留一个
	// tag=logbus.DefaultTag  dd_meta_channel=setting.DefaultChannel
	//logbus.Logger().Warn(zap.Int("money", 648)) // Deprecated 还可以使用，但不建议用
	logbus.Warn(zap.Int("money", 648)) // syntactic sugar

	// reason: 打点的推荐方式：使用预定义的tags
	// Print tga log and big query log. Old way
	/*_ = logbus.Logger(logbus.THINKINGDATA, logbus.BIGQUERY).Track(zap.String(thinkingdata.ACCOUNT, "111"), zap.String(thinkingdata.TYPE, thinkingdata.USER_SET_ONCE),
	zap.String("$user_id", "111"), zap.Time("$optime", time.Now()), zap.String(bigquery.TableNameKey, "oplog"),
	zap.String("player_name", "zhang liu"), zap.Int("level", 11), zap.Bool("bool", true), zap.Strings("strings", []string{"x", "y"}))*/ // Deprecated
	// Print tga log and big query log. New way
	_ = logbus.Tracker(config.THINKINGDATA, config.BIGQUERY).Track(zap.String(thinkingdata.ACCOUNT, "111"), zap.String(thinkingdata.TYPE, thinkingdata.USER_SET_ONCE),
		zap.String("$user_id", "111"), zap.Time("$optime", time.Now()), zap.String(bigquery.TableNameKey, "oplog"),
		zap.String("player_name", "zhang liu"), zap.Int("level", 11), zap.Bool("bool", true), zap.Strings("strings", []string{"x", "y"}))

	// hook field: add playerid=gtwefasfwad for all logs below
	globalfields.AppendGlobalFields(zap.String("playerid", "gtwefasfwad"))
	logbus.Warn(zap.Int("money", 648))
}
