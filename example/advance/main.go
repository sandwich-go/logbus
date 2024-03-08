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

	// 非线程安全
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

	// Print tga event log. thinkingdata.EVENT_ID is optional. Provide event_id, if want to update existing events
	_ = logbus.Tracker(logbus.THINKINGDATA).Track(logbus.String(thinkingdata.ACCOUNT, "111"), logbus.String(thinkingdata.TYPE, thinkingdata.TRACK),
		logbus.String(thinkingdata.EVENT_ID, "ID1"), logbus.String(thinkingdata.EVENT, "login"), logbus.Time("$optime", time.Now()),
		logbus.String("player_name", "zhang liu"), logbus.Int("level", 11), logbus.Bool("bool", true), logbus.Strings("strings", []string{"x", "y"}))

	// scope logger
	playerLogger := logbus.NewScopeLogger("Player", zap.String("playername", "zhangsong"), zap.Int("playerid", 123))
	guildLogger := logbus.NewScopeLogger("Guild", zap.String("guildname", "guild1"))
	playerLogger.Info("player gold", logbus.Int("money", 648))
	guildLogger.Info("guild gold", logbus.Int("money", 6480))

	// 增加全局域 非线程安全
	logbus.AppendGlobalFields(logbus.String("playerid", "gtwefasfwad"))
	logbus.Warn("", logbus.Int("money", 648)) // has extra global field

	q := logbus.NewQueue()
	q.Push(logbus.Int("i", 1))
	q.Push(logbus.Int("j", 2))
	logbus.Debug("queue", q.Retrieve()...)
}
