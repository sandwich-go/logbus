package logbus

import (
	"context"
	"testing"

	"go.uber.org/zap"
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
