package logbus

import (
	"errors"
	"testing"
	"time"

	"github.com/sandwich-go/logbus/bigquery"
	"github.com/sandwich-go/logbus/thinkingdata"

	"go.uber.org/zap"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	defer Close()
	EncodeConfig.CallerKey = "caller"
	m.Run()
}

func TestStdLogger(t *testing.T) {
	Init(NewConf(WithLogLevel(zap.InfoLevel), WithBufferedStdout(true), WithFetchLogContext(func() []zap.Field {
		return []zap.Field{String("dd_meta_channel", "fetch")}
	}), WithUseSystemClock(true)))
	defer resetLogBus()
	defer Close()
	Convey("test server log to stdout\n", t, func() {
		Debug("", Int("int", 111))
		Info("", Int("int", 111), String("str", "222"))
		Warn("", Int("int", 111), String("str", "222"), Bool("b", true))
		Error("", Int("int", 111), String("str", "222"), Bool("b", true), ErrorField(errors.New("this is a test error")))
		// StdLogger().WithOptions(zap.AddCallerSkip(10)).Fatal("fatal", zap.Int("int", 111), zap.String("str", "222"), zap.Bool("b", true), zap.Error(nil))
		So(gStdLogger.fetch, ShouldEqual, nil)
	})
}

func TestScopeLogger(t *testing.T) {
	refresh = true
	scopLogger1 := NewScopeLoggerWithFetchFunc("test1", nil)
	scopLogger2 := NewScopeLoggerWithFetchFunc("test2", func() []Field {
		return []zap.Field{String("dd_meta_channel", "test2")}
	})
	scopLogger1.Debug("", Int("int", 111))
	Convey("test scope log\n", t, func() {
		So(len(cacheGLogger), ShouldEqual, 2)
		So(refresh, ShouldEqual, true)
	})
	Init(NewConf(WithLogLevel(zap.InfoLevel), WithBufferedStdout(true), WithFetchLogContext(func() []zap.Field {
		return []zap.Field{String("dd_meta_channel", "test1")}
	})))
	scopLogger3 := NewScopeLogger("test3")
	defer resetLogBus()
	defer Close()
	Convey("test scope log\n", t, func() {
		scopLogger1.Debug("", Int("int", 111))                       // should not print
		scopLogger1.Info("", Int("int", 111), String("str", "222"))  // should print with dd_meta_channel test1
		scopLogger2.Warn("", Int("int", 111), String("str", "222"))  // should print with dd_meta_channel test2
		scopLogger3.Error("", Int("int", 111), String("str", "222")) // should print with dd_meta_channel test1
		So(cacheGLogger, ShouldEqual, nil)
		So(refresh, ShouldEqual, false)
	})
}

func TestPrintComplexTag(t *testing.T) {
	Init(NewConf(WithLogLevel(zap.DebugLevel)))
	defer resetLogBus()
	Convey("test only tga to stdout\n", t, func() {
		err := Tracker(THINKINGDATA).Track(String(thinkingdata.ACCOUNT, "111"), String(thinkingdata.TYPE, thinkingdata.USER_SET_ONCE),
			String("player_name", "zhang liu"), Int("level", 11), Bool("bool", true), Strings("strings", []string{"x", "y"}))
		So(err, ShouldBeNil)
	})

	Convey("test only bigquery to stdout\n", t, func() {
		err := Tracker(BIGQUERY).Track(String("$user_id", "111"), Time("$optime", time.Now()), String(bigquery.TableNameKey, "oplog"),
			String("player_name", "zhang liu"), Int("level", 11), Bool("bool", true), Strings("strings", []string{"x", "y"}))
		So(err, ShouldBeNil)
	})

	Convey("test bigquery and tga - UseRecord false to stdout\n", t, func() {
		err := Tracker(BIGQUERY, THINKINGDATA).Track(String(thinkingdata.ACCOUNT, "111"), String(thinkingdata.TYPE, thinkingdata.TRACK),
			String(thinkingdata.EVENT, "login"), String(thinkingdata.EVENT_ID, "ID1"),
			String("$user_id", "111"), Time("$optime", time.Now()), String(bigquery.TableNameKey, "oplog"),
			String("player_name", "zhang liu"), Int("level", 11), Bool("bool", true), Strings("strings", []string{"x", "y"}))
		So(err, ShouldBeNil)
	})

	Convey("test tga", t, func() {
		Convey("without event_id, should be fine", func() {
			err := Tracker(THINKINGDATA).Track(String(thinkingdata.ACCOUNT, "111"), String(thinkingdata.TYPE, thinkingdata.TRACK),
				String(thinkingdata.EVENT, "login"),
				String("player_name", "zhang liu"), Int("level", 11), Bool("bool", true), Strings("strings", []string{"x", "y"}))
			So(err, ShouldBeNil)
		})

		Convey("with legal event_id, should be fine", func() {
			err := Tracker(THINKINGDATA).Track(String(thinkingdata.ACCOUNT, "111"), String(thinkingdata.TYPE, thinkingdata.TRACK),
				String(thinkingdata.EVENT, "login"), String(thinkingdata.EVENT_ID, "ID1"),
				String("player_name", "zhang liu"), Int("level", 11), Bool("bool", true), Strings("strings", []string{"x", "y"}))
			So(err, ShouldBeNil)
		})

		Convey("with illegal event_id, should return error", func() {
			err := Tracker(THINKINGDATA).Track(String(thinkingdata.ACCOUNT, "111"), String(thinkingdata.TYPE, thinkingdata.TRACK),
				String(thinkingdata.EVENT, "login"), String(thinkingdata.EVENT_ID, "_dfa"),
				String("player_name", "zhang liu"), Int("level", 11), Bool("bool", true), Strings("strings", []string{"x", "y"}))
			So(err, ShouldNotBeNil)
		})
	})

	Convey("test tga and bigquery - UseRecord true to stdout\n", t, func() {
		bigquery.UseRecord = true
		err := Tracker(THINKINGDATA, BIGQUERY).Track(String(thinkingdata.ACCOUNT, "111"), String(thinkingdata.TYPE, thinkingdata.USER_SET_ONCE),
			String("$user_id", "111"), Time("$optime", time.Now()), String(bigquery.TableNameKey, "oplog"),
			String("player_name", "zhang liu"), Int("level", 11), Bool("bool", true), Strings("strings", []string{"x", "y"}))
		So(err, ShouldBeNil)
	})

	// Test cases for THINKINGDATACENTRALIZATION tag
	Convey("test tga data centralization", t, func() {
		Convey("valid case with appid", func() {
			// Adding `String("appid", "appid1")` as a proper test case
			err := Tracker(THINKINGDATACENTRALIZATION).Track(String(thinkingdata.TAG_APPID, "appid1"), String(thinkingdata.ACCOUNT, "111"),
				String(thinkingdata.TYPE, thinkingdata.TRACK), String(thinkingdata.EVENT, "login"),
				String("player_name", "zhang liu"), Int("level", 11), Bool("bool", true), Strings("strings", []string{"x", "y"}))
			So(err, ShouldBeNil)
		})

		Convey("missing appid should return error", func() {
			// Intentionally omit `appid` to trigger error
			err := Tracker(THINKINGDATACENTRALIZATION).Track(String(thinkingdata.ACCOUNT, "111"), String(thinkingdata.TYPE, thinkingdata.TRACK),
				String(thinkingdata.EVENT, "login"),
				String("player_name", "zhang liu"), Int("level", 11), Bool("bool", true), Strings("strings", []string{"x", "y"}))
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "tga appid must be assigned")
		})
	})
}

func TestTagLoggerThinkingData(t *testing.T) {
	Convey("test tag logger\n", t, func() {
		EncodeConfig.LevelKey = ""
		EncodeConfig.MessageKey = ""
		EncodeConfig.TimeKey = ""
		defer resetLogBus()
		Convey("test thinkingdata Stdout\n", func() {
			Init(NewConf(WithLogLevel(zap.DebugLevel)))
			properties := map[string]interface{}{"#ip": "10.0.0.1", "player_name": "zhang san", "level": 7}
			data, err := thinkingdata.Track("111", "", "login", "", properties)
			So(err, ShouldBeNil)
			Info("", zap.Object("tga", data))
		})
	})
}
