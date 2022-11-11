package benchmark

import (
	"bitbucket.org/funplus/sandwich/pkg/logbus/config"
	"testing"

	"bitbucket.org/funplus/sandwich/pkg/logbus"

	"go.uber.org/zap"
)

var testFields = []zap.Field{zap.Int("int", 123), zap.String("string", "string"), zap.Int("int", 123),
	zap.String("string", "string"), zap.Int("int", 123), zap.String("string", "string"), zap.Int("int", 123),
	zap.String("string", "string"), zap.Int("int", 123), zap.String("string", "string")}

func BenchmarkStdLogger1(b *testing.B) {
	// 27.4 ns/op	      48 B/op	       2 allocs/op
	b.Run("DisableOutput", func(b *testing.B) {
		logbus.Init(config.NewConf(config.WithLogLevel(zap.FatalLevel), config.WithOutputStdout(true)))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logbus.Info(testFields...)
			}
		})
	})
}

func BenchmarkStdLogger2(b *testing.B) {
	// {"log_level":"info","date":"2020-10-30T15:42:48.690+0800","dd_meta_channel":"server","tags":["logbus"],"int":123,"string":"string","int":123,"string":"string","int":123,"string":"string","int":123,"string":"string","int":123,"string":"string"}
	// 9785 ns/op  1579 B/op	       6 allocs/op
	b.Run("EnableStdout", func(b *testing.B) {
		logbus.Init(config.NewConf(config.WithLogLevel(zap.DebugLevel), config.WithOutputStdout(true)))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logbus.Info(testFields...)
			}
		})
	})
}

func BenchmarkStdLogger3(b *testing.B) {
	// {"tags":["logbus"],"string":"string","int":123}
	// 6406 ns/op
	b.Run("EnableStdout", func(b *testing.B) {
		config.EncodeConfig.TimeKey = ""
		config.EncodeConfig.LevelKey = ""
		config.EncodeConfig.MessageKey = ""
		logbus.Init(config.NewConf(config.WithLogLevel(zap.DebugLevel), config.WithOutputStdout(true) /*, logbus.WithLogId(false)*/))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logbus.InfoWithChannel("msg", zap.String("string", "string"), zap.Int("int", 123))
			}
		})
	})
}

func BenchmarkOutputFile(b *testing.B) {
	// {"log_level":"info","date":"2020-10-30T15:45:49.020+0800","dd_meta_channel":"server","int":123,"string":"string","int":123,"string":"string","int":123,"string":"string","int":123,"string":"string","int":123,"string":"string","log_xid":"buds9bbc1osh4qf798a0"}
	//  7189 ns/op	    1579 B/op	       6 allocs/op
	b.Run("EnableFile", func(b *testing.B) {
		logbus.Init(config.NewConf(config.WithLogLevel(zap.DebugLevel), config.WithOutputStdout(false), config.WithOutputLocalFile(true)))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logbus.Info(testFields...)
			}
		})
	})
}

func BenchmarkOutputFluentdSync(b *testing.B) {
	// 87689 ns/op	    3179 B/op	      38 allocs/op
	b.Run("SyncFluentd", func(b *testing.B) {
		logbus.Init(config.NewConf(config.WithLogLevel(zap.DebugLevel), config.WithOutputStdout(false), config.WithOutputLocalFile(false), config.WithOutputFluentd(true)))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logbus.Logger("debug").Info(testFields...)
			}
		})
	})
}

func BenchmarkOutputFluentdASync(b *testing.B) {
	// 1989 ns/op	    3433 B/op	      46 allocs/op
	b.Run("ASyncFluentd", func(b *testing.B) {
		logbus.Init(config.NewConf(config.WithLogLevel(zap.DebugLevel), config.WithOutputStdout(false), config.WithOutputLocalFile(false), config.WithOutputFluentd(true), config.WithFluentdAsync(true)))
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logbus.Logger("debug").Info(testFields...)
			}
		})
	})
}
