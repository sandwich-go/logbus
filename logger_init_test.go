package logbus

import (
	"errors"
	"os"
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
	Init(NewConf(WithLogLevel(zap.InfoLevel), WithBufferedStdout(true)))
	defer resetLogBus()
	defer Close()
	/*Convey("test PrintThingkingData to stdout\n", t, func() {
		properties := map[string]interface{}{"#ip": "10.0.0.2", "player_name": "zhang si", "level": 8, "data": []string{"x", "y"}}
		data, err := thinkingdata.User("111", "", thinkingdata.USER_SET_ONCE, properties)
		So(err, ShouldBeNil)
		getStdLogger().PrintThingkingData(data)
	})*/
	Convey("test server log to stdout\n", t, func() {
		Debug("", Int("int", 111))
		Info("", Int("int", 111), String("str", "222"))
		Warn("", Int("int", 111), String("str", "222"), Bool("b", true))
		Error("", Int("int", 111), String("str", "222"), Bool("b", true), ErrorField(errors.New("this is a test error")))
		//StdLogger().WithOptions(zap.AddCallerSkip(10)).Fatal("fatal", zap.Int("int", 111), zap.String("str", "222"), zap.Bool("b", true), zap.Error(nil))
	})
}

func TestPrintComplexTag(t *testing.T) {
	Init(NewConf(WithLogLevel(zap.DebugLevel), WithCallerSkip(2)))
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
			String(thinkingdata.EVENT, "login"),
			String("$user_id", "111"), Time("$optime", time.Now()), String(bigquery.TableNameKey, "oplog"),
			String("player_name", "zhang liu"), Int("level", 11), Bool("bool", true), Strings("strings", []string{"x", "y"}))
		So(err, ShouldBeNil)
	})
	Convey("test tga and bigquery - UseRecord true to stdout\n", t, func() {
		bigquery.UseRecord = true
		err := Tracker(THINKINGDATA, BIGQUERY).Track(String(thinkingdata.ACCOUNT, "111"), String(thinkingdata.TYPE, thinkingdata.USER_SET_ONCE),
			String("$user_id", "111"), Time("$optime", time.Now()), String(bigquery.TableNameKey, "oplog"),
			String("player_name", "zhang liu"), Int("level", 11), Bool("bool", true), Strings("strings", []string{"x", "y"}))
		So(err, ShouldBeNil)
	})
}
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
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
			data, err := thinkingdata.Track("111", "", "login", properties)
			So(err, ShouldBeNil)
			Info("", zap.Object("tga", data))
		})
	})
}
