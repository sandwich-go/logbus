package logbus

import (
	"errors"
	"github.com/sandwich-go/boost/xcmd"
	"github.com/sandwich-go/boost/xos"
	"github.com/sandwich-go/logbus/basics"
	"github.com/sandwich-go/logbus/config"
	"os"
	"testing"
	"time"

	"github.com/sandwich-go/logbus/bigquery"
	"github.com/sandwich-go/logbus/thinkingdata"
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	defer Close()
	config.EncodeConfig.CallerKey = "caller"
	m.Run()
}

func TestStdLogger(t *testing.T) {
	Init(config.NewConf(config.WithLogLevel(zap.InfoLevel), config.WithOutputFluentd(false), config.WithOutputStdout(true), config.WithOutputLocalFile(false), config.WithBufferedStdout(true)))
	defer basics.ResetLogBus()
	defer Close()
	Convey("test PrintThingkingData to stdout\n", t, func() {
		properties := map[string]interface{}{"#ip": "10.0.0.2", "player_name": "zhang si", "level": 8, "data": []string{"x", "y"}}
		data, err := thinkingdata.User("111", "", thinkingdata.USER_SET_ONCE, properties)
		So(err, ShouldBeNil)
		GetStdLogger().PrintThingkingData(data)
	})
	Convey("test server log to stdout\n", t, func() {
		Debug(zap.Int("int", 111))
		Info(zap.Int("int", 111), zap.String("str", "222"))
		Warn(zap.Int("int", 111), zap.String("str", "222"), zap.Bool("b", true))
		Error(zap.Int("int", 111), zap.String("str", "222"), zap.Bool("b", true), zap.Error(errors.New("this is a test error")))
		//StdLogger().WithOptions(zap.AddCallerSkip(10)).Fatal("fatal", zap.Int("int", 111), zap.String("str", "222"), zap.Bool("b", true), zap.Error(nil))
	})
}

func TestPrintComplexTag(t *testing.T) {
	Init(config.NewConf(config.WithLogLevel(zap.DebugLevel), config.WithOutputFluentd(false), config.WithOutputStdout(true), config.WithOutputLocalFile(false), config.WithCallerSkip(2)))
	defer basics.ResetLogBus()
	Convey("test only tga to stdout\n", t, func() {
		err := Tracker(THINKINGDATA).Track(zap.String(thinkingdata.ACCOUNT, "111"), zap.String(thinkingdata.TYPE, thinkingdata.USER_SET_ONCE),
			zap.String("player_name", "zhang liu"), zap.Int("level", 11), zap.Bool("bool", true), zap.Strings("strings", []string{"x", "y"}))
		So(err, ShouldBeNil)
	})
	Convey("test only bigquery to stdout\n", t, func() {
		err := Tracker(BIGQUERY).Track(zap.String("$user_id", "111"), zap.Time("$optime", time.Now()), zap.String(bigquery.TableNameKey, "oplog"),
			zap.String("player_name", "zhang liu"), zap.Int("level", 11), zap.Bool("bool", true), zap.Strings("strings", []string{"x", "y"}))
		So(err, ShouldBeNil)
	})
	Convey("test bigquery and tga - UseRecord false to stdout\n", t, func() {
		err := Tracker(BIGQUERY, THINKINGDATA).Track(zap.String(thinkingdata.ACCOUNT, "111"), zap.String(thinkingdata.TYPE, thinkingdata.TRACK),
			zap.String(thinkingdata.EVENT, "login"),
			zap.String("$user_id", "111"), zap.Time("$optime", time.Now()), zap.String(bigquery.TableNameKey, "oplog"),
			zap.String("player_name", "zhang liu"), zap.Int("level", 11), zap.Bool("bool", true), zap.Strings("strings", []string{"x", "y"}))
		So(err, ShouldBeNil)
	})
	Convey("test tga and bigquery - UseRecord true to stdout\n", t, func() {
		bigquery.UseRecord = true
		err := Tracker(THINKINGDATA, BIGQUERY).Track(zap.String(thinkingdata.ACCOUNT, "111"), zap.String(thinkingdata.TYPE, thinkingdata.USER_SET_ONCE),
			zap.String("$user_id", "111"), zap.Time("$optime", time.Now()), zap.String(bigquery.TableNameKey, "oplog"),
			zap.String("player_name", "zhang liu"), zap.Int("level", 11), zap.Bool("bool", true), zap.Strings("strings", []string{"x", "y"}))
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
func TestTagLoggerConf(t *testing.T) {
	if xos.Exists("/tmp/fun-collector.sock") {
		Convey("test tag logger conf to stdout file and fluentd\n", t, func() {
			Init(config.NewConf(config.WithLogLevel(zap.DebugLevel), config.WithOutputFluentd(true), config.WithOutputStdout(true), config.WithOutputLocalFile(true)))
			defer basics.ResetLogBus()
			Logger("debug", "debug.access").Debug(zap.Int("int", 123), zap.Time("time", time.Now()), zap.Any("any", map[int]int{1: 1}))
			Logger("debug.access").Warn(zap.Error(errors.New("err")), zap.Bool("book", false))
		})
	}
}

func TestTagLoggerThinkingData(t *testing.T) {
	Convey("test tag logger\n", t, func() {
		config.EncodeConfig.LevelKey = ""
		config.EncodeConfig.MessageKey = ""
		config.EncodeConfig.TimeKey = ""
		defer basics.ResetLogBus()
		Convey("test thinkingdata Stdout\n", func() {
			Init(config.NewConf(config.WithLogLevel(zap.DebugLevel), config.WithOutputStdout(true)))
			properties := map[string]interface{}{"#ip": "10.0.0.1", "player_name": "zhang san", "level": 7}
			data, err := thinkingdata.Track("111", "", "login", properties)
			So(err, ShouldBeNil)
			Info(zap.Object("tga", data))
		})
		Convey("test thinkingdata User to file\n", func() {
			Init(config.NewConf(config.WithLogLevel(zap.DebugLevel), config.WithOutputStdout(false), config.WithOutputLocalFile(true)))
			properties := map[string]interface{}{"#ip": "10.0.0.2", "player_name": "zhang si", "level": 8}
			data, err := thinkingdata.User("111", "", thinkingdata.USER_SET_ONCE, properties)
			So(err, ShouldBeNil)
			Logger("debug").Info(zap.Object("tga", data))
		})
		if xcmd.IsTrue(xcmd.GetOptWithEnv("sandwich_test_enable_fluentd")) {
			Convey("test thinkingdata User to fluentd\n", func() {
				Init(config.NewConf(config.WithLogLevel(zap.DebugLevel), config.WithOutputStdout(false), config.WithOutputLocalFile(false), config.WithOutputFluentd(true)))
				properties := map[string]interface{}{"#ip": "10.0.0.3", "player_name": "zhang wu", "level": 9}
				data, err := thinkingdata.User("111", "", thinkingdata.USER_SET_ONCE, properties)
				So(err, ShouldBeNil)
				Logger("debug").Info(zap.Object("tga", data))
			})
		}
	})
}

func TestFluentdAndPrometheus(t *testing.T) {
	Convey("test log with fluentd and prometheus open\n", t, func() {
		Close()
		withFluentd := xcmd.IsTrue(xcmd.GetOptWithEnv("sandwich_test_enable_fluentd"))
		Init(config.NewConf(config.WithLogLevel(zap.DebugLevel), config.WithOutputFluentd(withFluentd), config.WithOutputStdout(false), config.WithOutputLocalFile(false), config.WithMonitorOutput(Prometheus), config.WithDefaultPrometheusListenAddress(":8991")))
		Logger("debug").Info(zap.Bool("prometheus", true))
	})
}
