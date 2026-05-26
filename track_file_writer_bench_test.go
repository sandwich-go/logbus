package logbus

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/sandwich-go/logbus/monitor"
)

// ─────────────────────────────────────────────────────────────────────────────
// track 文件落盘大并发 / 大数据量压测
//
// 设计原则：
//   1. 全部用 Go Benchmark 形式（go test -bench），便于复跑、对比、出 ns/op + allocs/op；
//   2. 输出走 t.TempDir，Benchmark 结束自动清理，不污染本地文件系统；
//   3. 默认 -short 跳过长跑用例，方便 CI 只跑轻量基线。
//
// 跑法示例：
//   go test -run=^$ -bench=BenchmarkTrackFileWriter -benchmem -benchtime=3s ./...
//   go test -run=^$ -bench=BenchmarkTrackFileWriter_Throughput -cpu=1,4,16
//   go test -run=^$ -bench=BenchmarkTrackFileWriter_DropOnSaturate -benchtime=2s
// ─────────────────────────────────────────────────────────────────────────────

// countingReporter 实现 monitor.Reporter，用于在压测中精确捕获 dropped / writeFailed 计数
type countingReporter struct {
	counts sync.Map // metricName -> *atomic.Int64
}

func (r *countingReporter) get(name string) *atomic.Int64 {
	if v, ok := r.counts.Load(name); ok {
		return v.(*atomic.Int64)
	}
	v := new(atomic.Int64)
	actual, _ := r.counts.LoadOrStore(name, v)
	return actual.(*atomic.Int64)
}

func (r *countingReporter) Count(metric string, value int64, _ prometheus.Labels) error {
	r.get(metric).Add(value)
	return nil
}

func (r *countingReporter) Gauge(string, float64, prometheus.Labels) error          { return nil }
func (r *countingReporter) Timing(string, time.Duration, prometheus.Labels) error { return nil }

// withCountingReporter 在 benchmark 期间临时替换全局 reporter，结束后还原。
// 注意：Go benchmark 默认不并发跑多个用例，所以全局替换是安全的。
func withCountingReporter(b *testing.B) *countingReporter {
	b.Helper()
	orig := monitor.DefaultMetricsReporter
	rep := &countingReporter{}
	monitor.DefaultMetricsReporter = rep
	b.Cleanup(func() { monitor.DefaultMetricsReporter = orig })
	return rep
}

// ── payload 构造：复用同一个底层 buffer，避免 benchmark 体内分配影响测量 ─

// makeTrackPayload 生成一条形如 TrackWithChannel 输出的 JSON 日志。
// msgSize 用来灌大字段，便于压"大数据量"维度。
func makeTrackPayload(channel string, msgSize int) []byte {
	// msg 字段是字符串。为模拟真实业务，里面再嵌一段 escape 后的 JSON。
	// 单条整体长度 ≈ msgSize + 固定包头 ~80 字节。
	body := make([]byte, msgSize)
	for i := range body {
		body[i] = 'x'
	}
	return []byte(fmt.Sprintf(
		`{"log_level":"track","tags":"logbus","%s":"%s","%s":"%s"}`,
		Meta, channel, MsgBody, body,
	))
}

// 不同体量的 payload 表，覆盖 small / medium / large 三档
var payloadSizes = []struct {
	name string
	size int
}{
	{"S128B", 128},
	{"M1KB", 1024},
	{"L8KB", 8 * 1024},
}

// ── 1. 吞吐基线：单 channel，串行写入，测纯路径开销 ──────────────────

func BenchmarkTrackFileWriter_Throughput_Serial(b *testing.B) {
	for _, pl := range payloadSizes {
		b.Run(pl.name, func(b *testing.B) {
			dir := b.TempDir()
			ws := NewTrackFileWriteSyncerWithBufSize(dir, HourlyRotation, 1<<14)
			b.Cleanup(func() { _ = ws.Close() })

			payload := makeTrackPayload("bench", pl.size)
			b.SetBytes(int64(len(payload)))
			b.ReportAllocs()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				_, _ = ws.Write(payload)
			}
			b.StopTimer()
			if err := ws.Sync(); err != nil {
				b.Fatalf("sync: %v", err)
			}
		})
	}
}

// ── 2. 大并发：多 goroutine 共写一个 channel ───────────────────────

func BenchmarkTrackFileWriter_Throughput_Parallel(b *testing.B) {
	for _, pl := range payloadSizes {
		b.Run(pl.name, func(b *testing.B) {
			dir := b.TempDir()
			// 队列足够大，避免压测被 drop 掩盖真实开销
			ws := NewTrackFileWriteSyncerWithBufSize(dir, HourlyRotation, 1<<16)
			b.Cleanup(func() { _ = ws.Close() })

			payload := makeTrackPayload("bench", pl.size)
			b.SetBytes(int64(len(payload)))
			b.ReportAllocs()
			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _ = ws.Write(payload)
				}
			})
			b.StopTimer()
			if err := ws.Sync(); err != nil {
				b.Fatalf("sync: %v", err)
			}
		})
	}
}

// ── 3. 多 channel 并发：worker 之间无竞争路径 ──────────────────────
//
// 业务里 thinkingdata / bi / bigquery_xxx 多 channel 并行写，
// 这里验证 per-channel goroutine 拓扑下吞吐是否随 channel 数线性扩展。
func BenchmarkTrackFileWriter_MultiChannel(b *testing.B) {
	channelCounts := []int{1, 4, 16}
	for _, nc := range channelCounts {
		b.Run(fmt.Sprintf("channels=%d", nc), func(b *testing.B) {
			dir := b.TempDir()
			ws := NewTrackFileWriteSyncerWithBufSize(dir, HourlyRotation, 1<<16)
			b.Cleanup(func() { _ = ws.Close() })

			// 预生成各 channel 的 payload
			payloads := make([][]byte, nc)
			for i := 0; i < nc; i++ {
				payloads[i] = makeTrackPayload(fmt.Sprintf("ch_%02d", i), 512)
			}
			b.SetBytes(int64(len(payloads[0])))
			b.ReportAllocs()
			b.ResetTimer()

			var idx atomic.Uint64
			b.RunParallel(func(pb *testing.PB) {
				// 每个 P 拿一个起始 channel，循环递增模拟"散列到 N 个 channel"
				start := idx.Add(1) - 1
				i := uint64(0)
				for pb.Next() {
					p := payloads[(start+i)%uint64(nc)]
					_, _ = ws.Write(p)
					i++
				}
			})
			b.StopTimer()
			if err := ws.Sync(); err != nil {
				b.Fatalf("sync: %v", err)
			}
		})
	}
}

// ── 4. 切割边界压力：MinuteRotation 下持续写入 ─────────────────────
//
// 注意 benchmark 不能保证恰好横跨整分钟，但 b.N 充分大时会自然覆盖到至少一次切片。
// 这里主要验证切片 + 关旧文件 + 开新文件路径不会成为吞吐瓶颈。
func BenchmarkTrackFileWriter_MinuteRotation(b *testing.B) {
	dir := b.TempDir()
	ws := NewTrackFileWriteSyncerWithBufSize(dir, MinuteRotation, 1<<16)
	b.Cleanup(func() { _ = ws.Close() })

	payload := makeTrackPayload("rot", 512)
	b.SetBytes(int64(len(payload)))
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = ws.Write(payload)
		}
	})
	b.StopTimer()
	if err := ws.Sync(); err != nil {
		b.Fatalf("sync: %v", err)
	}
}

// ── 5. 队列饱和与丢弃：故意把 BufSize 调小，用大量 goroutine 压满 ──
//
// 关注三件事：
//   1. Write 不阻塞（每次调用都立即返回）；
//   2. dropped 计数 = 实际投递 - 落盘；
//   3. 调用方拿到 n == len(p)，err == nil（合约不变）。
//
// 因为是 benchmark，断言只能在 b.Cleanup 里做轻量校验。
func BenchmarkTrackFileWriter_DropOnSaturate(b *testing.B) {
	if testing.Short() {
		b.Skip("short mode")
	}
	rep := withCountingReporter(b)
	dropped := rep.get(MetricTrackDropped)

	dir := b.TempDir()
	// 故意把缓冲设得很小，制造队列饱和
	ws := NewTrackFileWriteSyncerWithBufSize(dir, HourlyRotation, 16)

	payload := makeTrackPayload("flood", 512)
	b.SetBytes(int64(len(payload)))
	b.ReportAllocs()

	var written int64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var local int64
		for pb.Next() {
			n, err := ws.Write(payload)
			if err != nil || n != len(payload) {
				b.Fatalf("Write contract violated: n=%d err=%v", n, err)
			}
			local++
		}
		atomic.AddInt64(&written, local)
	})
	b.StopTimer()

	_ = ws.Sync()
	_ = ws.Close()

	// 落盘行数（聚合所有时间槽 + 文件）
	disk := countDiskLines(b, dir)

	dr := dropped.Load()
	b.ReportMetric(float64(dr)/float64(written)*100, "drop%")
	b.ReportMetric(float64(disk), "lines/disk")

	// 合约校验：dropped + disk 应当等于 written（允许 ±1 是 timer flush 半路 close 的窗口）
	delta := written - (dr + disk)
	if delta < -2 || delta > 2 {
		b.Fatalf("accounting mismatch: written=%d dropped=%d disk=%d delta=%d",
			written, dr, disk, delta)
	}
}

// ── 6. 极端并发：goroutine 数 ≫ GOMAXPROCS，验证调度争抢下不退化 ──

func BenchmarkTrackFileWriter_HighGoroutine(b *testing.B) {
	if testing.Short() {
		b.Skip("short mode")
	}
	dir := b.TempDir()
	ws := NewTrackFileWriteSyncerWithBufSize(dir, HourlyRotation, 1<<16)
	b.Cleanup(func() { _ = ws.Close() })

	payload := makeTrackPayload("storm", 256)
	b.SetBytes(int64(len(payload)))
	b.ReportAllocs()

	// 1k goroutine 同时写，远大于 GOMAXPROCS
	const goroutines = 1024
	perG := b.N / goroutines
	if perG <= 0 {
		perG = 1
	}

	var wg sync.WaitGroup
	wg.Add(goroutines)
	b.ResetTimer()
	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < perG; i++ {
				_, _ = ws.Write(payload)
			}
		}()
	}
	wg.Wait()
	b.StopTimer()

	if err := ws.Sync(); err != nil {
		b.Fatalf("sync: %v", err)
	}
	b.ReportMetric(float64(goroutines), "goroutines")
}

// ── 7. 大数据量长跑：模拟连续秒级高 QPS ────────────────────────────
//
// 不依赖 b.N，固定按时间窗压测，最后输出 QPS / drop% / disk lines。
// 跑法：go test -run=^$ -bench=BenchmarkTrackFileWriter_Soak -benchtime=10s
func BenchmarkTrackFileWriter_Soak(b *testing.B) {
	if testing.Short() {
		b.Skip("short mode")
	}
	rep := withCountingReporter(b)
	dropped := rep.get(MetricTrackDropped)

	dir := b.TempDir()
	ws := NewTrackFileWriteSyncerWithBufSize(dir, HourlyRotation, 1<<16)

	payload := makeTrackPayload("soak", 512)
	b.SetBytes(int64(len(payload)))
	b.ReportAllocs()

	// 从 b.N 推一个目标"运行时长"。Go benchmark 会自动把 b.N 加大到 -benchtime 指定时间。
	// 这里只把 b.N 当真实写入次数预算用，不再二次循环。
	var written int64
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var local int64
		for pb.Next() {
			_, _ = ws.Write(payload)
			local++
		}
		atomic.AddInt64(&written, local)
	})
	b.StopTimer()

	_ = ws.Sync()
	_ = ws.Close()

	disk := countDiskLines(b, dir)
	dr := dropped.Load()

	b.ReportMetric(float64(written), "writes")
	b.ReportMetric(float64(disk), "lines/disk")
	b.ReportMetric(float64(dr)/float64(written)*100, "drop%")
	b.ReportMetric(float64(runtime.NumGoroutine()), "g_after")
}

// ── 辅助：递归统计 dir 下所有 .log 文件的总行数 ────────────────────

func countDiskLines(b *testing.B, dir string) int64 {
	b.Helper()
	var total int64
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		data, rErr := os.ReadFile(path)
		if rErr != nil {
			return rErr
		}
		// 每行以 '\n' 结尾，最后一行也是；空文件计 0
		for _, c := range data {
			if c == '\n' {
				total++
			}
		}
		return nil
	})
	if err != nil {
		b.Fatalf("walk %s: %v", dir, err)
	}
	return total
}
