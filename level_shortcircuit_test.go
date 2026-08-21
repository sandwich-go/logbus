package logbus

import (
	"bytes"
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 本文件锁住「级别关闭时不构造字段」这个不变量，并给出前后对比的 benchmark。
//
// 动机：GLogger 的各级别方法先 append msg 字段、StdLogger.*WithChannel 先调
// StdLogger.fields（内含 FetchLogContext 钩子 + 新切片分配），两者都发生在 zap 自身的级别
// 判定之前。业务侧实测（ria worldmap 压测 2026-08-21）GLogger.Debug 与 StdLogger.fields
// 合计占进程总分配量 18.7%，其中相当一部分是级别关闭时也照样付的。
//
// 跑法：go test -run=NONE -bench=BenchmarkDisabledLevel -benchmem ./

func benchFields() []Field {
	return []Field{String("k1", "v1"), Int("k2", 2), Uint64("k3", 3)}
}

// TestDisabledLevelSkipsFetchLogContext 锁住核心不变量：级别关闭时连 FetchLogContext
// 都不该被调用（它是 fields 构造链的入口，被调用即说明短路失效）。
func TestDisabledLevelSkipsFetchLogContext(t *testing.T) {
	origin := Setting.FetchLogContext
	t.Cleanup(func() { Setting.FetchLogContext = origin })

	var fetched int
	Setting.FetchLogContext = func() []Field {
		fetched++
		return []Field{String("ctx", "1")}
	}

	originLevel := GetLogLevel()
	t.Cleanup(func() { SetLogLevel(originLevel) })
	SetLogLevel(zap.ErrorLevel)

	lg := NewScopeLogger("shortcircuit")
	lg.Debug("msg", benchFields()...)
	lg.Info("msg", benchFields()...)
	lg.Warn("msg", benchFields()...)
	lg.DebugWithContext(context.Background(), "msg", benchFields()...)

	if fetched != 0 {
		t.Errorf("级别为 Error 时 FetchLogContext 被调用了 %d 次，短路失效", fetched)
	}
}

// TestEnabledLevelStillLogs 反面用例：级别放行时链路照常走通，短路没有把正常日志吃掉。
func TestEnabledLevelStillLogs(t *testing.T) {
	origin := Setting.FetchLogContext
	t.Cleanup(func() { Setting.FetchLogContext = origin })

	var fetched int
	Setting.FetchLogContext = func() []Field {
		fetched++
		return []Field{String("ctx", "1")}
	}

	originLevel := GetLogLevel()
	t.Cleanup(func() { SetLogLevel(originLevel) })
	SetLogLevel(zap.DebugLevel)

	NewScopeLogger("shortcircuit").Debug("msg", benchFields()...)

	if fetched == 0 {
		t.Error("级别为 Debug 时 FetchLogContext 未被调用，日志链路被短路误伤")
	}
}

// BenchmarkDisabledLevelDebug 量级别关闭时一次 Debug 调用的开销。
// 短路生效后应为 0 allocs/op。
func BenchmarkDisabledLevelDebug(b *testing.B) {
	originLevel := GetLogLevel()
	b.Cleanup(func() { SetLogLevel(originLevel) })
	SetLogLevel(zap.ErrorLevel)

	origin := Setting.FetchLogContext
	b.Cleanup(func() { Setting.FetchLogContext = origin })
	Setting.FetchLogContext = func() []Field { return []Field{String("ctx", "1")} }

	lg := NewScopeLogger("bench")
	f := benchFields()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lg.Debug("some message", f...)
	}
}

// BenchmarkEnabledLevelDebug 对照组：级别放行时的开销，确认短路没有拖慢正常路径。
func BenchmarkEnabledLevelDebug(b *testing.B) {
	originLevel := GetLogLevel()
	b.Cleanup(func() { SetLogLevel(originLevel) })
	SetLogLevel(zap.DebugLevel)

	lg := NewScopeLogger("bench")
	f := benchFields()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lg.Debug("some message", f...)
	}
}

// newLevelledStdLogger 造一个带独立 AtomicLevel 与独立输出缓冲的 StdLogger，
// 用于构造「实例级别与全局默认 logger 级别不一致」的场景。
func newLevelledStdLogger(lvl zapcore.Level) (*StdLogger, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.AddSync(buf),
		zap.NewAtomicLevelAt(lvl),
	)
	return newStdLogger(zap.New(core), nil), buf
}

// TestDebugWithChannelUsesInstanceLogger 锁住 P1：DebugWithChannel 的级别判定与写入都必须
// 落在本实例的 stdLogger 上，不能用 gStdLogger。
//
// 触发条件：SetGlobalGLogger 注入一个 Debug 级别开启的自定义 StdLogger，而全局默认 logger
// 处于 Error 级别。按 gStdLogger 判定会在此处直接 return，本该输出的 Debug 日志被丢掉。
func TestDebugWithChannelUsesInstanceLogger(t *testing.T) {
	originLevel := GetLogLevel()
	t.Cleanup(func() { SetLogLevel(originLevel) })
	SetLogLevel(zap.ErrorLevel) // 全局默认 logger 关闭 Debug

	custom, buf := newLevelledStdLogger(zap.DebugLevel) // 实例开启 Debug
	g := &GLogger{channelKey: "ch", stdLogger: custom}

	g.DebugWithChannel("ch", "must be written", String("k", "v"))

	if buf.Len() == 0 {
		t.Error("实例 stdLogger 为 Debug 级别，DebugWithChannel 应写入；" +
			"按 gStdLogger 判定会把本该输出的日志丢掉")
	}
}

// TestDebugWithChannelShortCircuitsOnInstanceLevel 锁住反方向：实例级别关闭时必须短路，
// 不得因为全局默认 logger 开启就继续构造字段、甚至写到全局 logger 去。
func TestDebugWithChannelShortCircuitsOnInstanceLevel(t *testing.T) {
	originLevel := GetLogLevel()
	t.Cleanup(func() { SetLogLevel(originLevel) })
	SetLogLevel(zap.DebugLevel) // 全局默认 logger 开启 Debug

	originFetch := Setting.FetchLogContext
	t.Cleanup(func() { Setting.FetchLogContext = originFetch })
	var fetched int
	Setting.FetchLogContext = func() []Field {
		fetched++
		return []Field{String("ctx", "1")}
	}

	custom, buf := newLevelledStdLogger(zap.ErrorLevel) // 实例关闭 Debug
	g := &GLogger{channelKey: "ch", stdLogger: custom}

	g.DebugWithChannel("ch", "must be dropped", String("k", "v"))

	if buf.Len() != 0 {
		t.Errorf("实例 stdLogger 为 Error 级别，DebugWithChannel 不应写入，实际写了 %d 字节", buf.Len())
	}
	// fetched 走的是全局 gStdLogger 的 fields 链路；被调用说明短路失效且写错了目标。
	if fetched != 0 {
		t.Errorf("实例级别关闭时仍构造了字段（FetchLogContext 被调用 %d 次），"+
			"说明判定用的不是实例的 stdLogger", fetched)
	}
}

// TestDepthLoggersShortCircuit 锁住 GDebugDepth/GInfoDepth/GWarnDepth 的短路。
//
// 这三个方法先 append MsgBody 再 getDepthLogger，而后者未命中时会 WithOptions 克隆一个
// zap logger 并写入 sync.Map，代价都在级别判定之前。
func TestDepthLoggersShortCircuit(t *testing.T) {
	originLevel := GetLogLevel()
	t.Cleanup(func() { SetLogLevel(originLevel) })
	SetLogLevel(zap.ErrorLevel)

	custom, buf := newLevelledStdLogger(zap.ErrorLevel)
	g := &GLogger{channelKey: "ch", stdLogger: custom}

	g.GDebugDepth(1, "msg", String("k", "v"))
	g.GInfoDepth(1, "msg", String("k", "v"))
	g.GWarnDepth(1, "msg", String("k", "v"))

	if buf.Len() != 0 {
		t.Errorf("级别为 Error 时 Depth 系方法不应写入，实际写了 %d 字节", buf.Len())
	}
	// 短路发生在 getDepthLogger 之前，因此不该有 depth logger 被创建。
	var created int
	g.depthLogger.Range(func(any, any) bool { created++; return true })
	if created != 0 {
		t.Errorf("级别关闭时创建了 %d 个 depth logger，短路应在 getDepthLogger 之前", created)
	}
}
