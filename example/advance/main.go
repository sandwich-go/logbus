package main

import (
	"time"

	"github.com/sandwich-go/boost/xpanic"
	bislg "github.com/sandwich-go/logbus/bi/slg"
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
		logbus.WithDefaultTag("Advance"),
	),
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

	ps := make(map[string]interface{})
	ps["level"] = 11
	_ = logbus.Tracker(logbus.THINKINGDATA).TrackWithTGAData(thinkingdata.Data{
		AccountId:    "2222",
		DistinctId:   "",
		Type:         thinkingdata.TRACK,
		Time:         time.Now().String(),
		EventName:    "login",
		EventId:      "ID1",
		FirstCheckId: "",
		Ip:           "",
		UUID:         "",
		Properties:   ps,
	})

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

	logbus.ReservedGlobalFields = nil
	logbus.SetGlobalFields(nil)
	logbus.Info("clean log")

	biTracker := logbus.NewTracker(logbus.WithTags(logbus.BI), logbus.WithBiAppID("gof.prod.global"))
	bislg.MustInitialize(biTracker)
	_ = bislg.TrackUserChatPrivate(1, 2, "hello world", &bislg.UserChatOpts{
		IP:         "1.1.1.1",
		FpID:       "1",
		ServerID:   1,
		GameUser:   "zhang san",
		TotalPower: 100.123,
		TransLang:  "zh",
	})

	logbus.Init(logbus.NewConf(
		logbus.WithDev(true),
		logbus.WithEnableTraceLevel(false),
	),
	)

	err := logbus.Tracker(logbus.THINKINGDATA).Track(logbus.String(thinkingdata.ACCOUNT, "111"), logbus.String(thinkingdata.TYPE, thinkingdata.TRACK),
		logbus.String(thinkingdata.EVENT_ID, "ID1"), logbus.String(thinkingdata.EVENT, "login"))
	xpanic.WhenError(err)
	_ = bislg.TrackUserChatPrivate(3, 4, "hello world", &bislg.UserChatOpts{
		IP:         "3.3.3.3",
		FpID:       "3",
		ServerID:   1,
		GameUser:   "li si",
		TotalPower: 99.21,
		TransLang:  "zh",
	})
}
