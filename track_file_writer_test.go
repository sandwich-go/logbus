package logbus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sandwich-go/logbus/bigquery"
	. "github.com/smartystreets/goconvey/convey"
)

// buildTrackJSON 构造一条 track 级别的 JSON 日志（模拟 TrackWithChannel 输出）
func buildTrackJSON(channel, msgContent string) []byte {
	return []byte(fmt.Sprintf(`{"log_level":"track","date":"2026-05-09T10:00:00.000+0800","tags":"logbus","%s":"%s","%s":%s}`,
		Meta, channel, MsgBody, msgContent))
}

// readFileLines 读取文件所有非空行
func readFileLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, strings.TrimSpace(line))
		}
	}
	return lines, nil
}

// expectedTrackFilePath 按 ops 规范拼出当前时间下的预期文件路径：
//
//	{baseDir}/{channel}/{yyyymmdd}/{channel}_{hostname}_{slot}.log
//
// channel 与 hostname 经 sanitizeTrackFileComponent 清洗，hostname 为空时
// 由 sanitizeTrackFileComponent 兜底为 "unknown"。
func expectedTrackFilePath(baseDir, channel, hostname string, rot TrackRotation, t time.Time) string {
	chSeg := sanitizeTrackFileComponent(channel)
	hostSeg := sanitizeTrackFileComponent(hostname)
	slot := rot.timeSlot(t)
	date := t.Format("20060102")
	return filepath.Join(baseDir, chSeg, date, fmt.Sprintf("%s_%s_%s.log", chSeg, hostSeg, slot))
}

// runtimeHostName 返回当前进程实际使用的 hostname 段，与 NewTrackFileWriteSyncerFromConfig 一致：
// 优先 os.Hostname，失败则 HOSTNAME env，再失败由 sanitizeTrackFileComponent 兜底 unknown。
func runtimeHostName() string {
	if v, err := os.Hostname(); err == nil && v != "" {
		return v
	}
	if v := os.Getenv("HOSTNAME"); v != "" {
		return v
	}
	return "unknown"
}

// ── extractChannelAndMsg 单元测试 ──────────────────────────────────────────

func TestExtractChannelAndMsg_StringMsg(t *testing.T) {
	Convey("msg 是字符串时提取内容", t, func() {
		raw := fmt.Sprintf(`{"%s":"thinkingdata","%s":"{\"#account_id\":\"111\"}"}`, Meta, MsgBody)
		channel, msg, ok := extractChannelAndMsg([]byte(raw))
		So(ok, ShouldBeTrue)
		So(channel, ShouldEqual, "thinkingdata")
		So(string(msg), ShouldEqual, `{"#account_id":"111"}`)
	})
}

func TestExtractChannelAndMsg_ObjectMsg(t *testing.T) {
	Convey("msg 是 JSON 对象时直接提取 raw bytes", t, func() {
		raw := fmt.Sprintf(`{"%s":"bi","%s":{"app_id":"gof","event":"login"}}`, Meta, MsgBody)
		channel, msg, ok := extractChannelAndMsg([]byte(raw))
		So(ok, ShouldBeTrue)
		So(channel, ShouldEqual, "bi")
		So(string(msg), ShouldContainSubstring, "app_id")
	})
}

func TestExtractChannelAndMsg_BigQueryTableName(t *testing.T) {
	Convey("bigquery 日志提取 channel 时附加表名用于文件名", t, func() {
		raw := fmt.Sprintf(`{"%s":"%s","%s":"oplog","%s":"{\"id\":\"1\"}"}`,
			Meta, BIGQUERY, bigquery.TableNameKey, MsgBody)
		channel, msg, ok := extractChannelAndMsg([]byte(raw))
		So(ok, ShouldBeTrue)
		So(channel, ShouldEqual, BIGQUERY+"_oplog")
		So(string(msg), ShouldEqual, `{"id":"1"}`)
	})
}

func TestExtractChannelAndMsg_MissingFields(t *testing.T) {
	Convey("缺少 channel 或 msg 字段时返回 false", t, func() {
		Convey("缺少 channel", func() {
			raw := fmt.Sprintf(`{"%s":"hello"}`, MsgBody)
			_, _, ok := extractChannelAndMsg([]byte(raw))
			So(ok, ShouldBeFalse)
		})
		Convey("缺少 msg", func() {
			raw := fmt.Sprintf(`{"%s":"thinkingdata"}`, Meta)
			_, _, ok := extractChannelAndMsg([]byte(raw))
			So(ok, ShouldBeFalse)
		})
		Convey("非 JSON", func() {
			_, _, ok := extractChannelAndMsg([]byte("not json"))
			So(ok, ShouldBeFalse)
		})
	})
}

// ── TrackFileWriteSyncer 功能测试 ──────────────────────────────────────────

func TestTrackFileWriteSyncer_WriteToFile(t *testing.T) {
	Convey("Write 写入正确的文件，只含 msg 内容", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncer(dir, HourlyRotation)

		payload := `{"#account_id":"111","#type":"track"}`
		raw := fmt.Sprintf(`{"%s":"thinkingdata","%s":"%s"}`,
			Meta, MsgBody, strings.ReplaceAll(payload, `"`, `\"`))

		n, err := ws.Write([]byte(raw))
		So(err, ShouldBeNil)
		So(n, ShouldEqual, len(raw))

		So(ws.Sync(), ShouldBeNil)
		So(ws.Close(), ShouldBeNil)

		fname := expectedTrackFilePath(dir, "thinkingdata", runtimeHostName(), HourlyRotation, time.Now())
		lines, readErr := readFileLines(fname)
		So(readErr, ShouldBeNil)
		So(lines, ShouldHaveLength, 1)
		So(lines[0], ShouldEqual, payload)
	})
}

func TestTrackFileWriteSyncer_BigQueryTableNameInFileName(t *testing.T) {
	Convey("BigQuery track 文件名包含表名", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncer(dir, HourlyRotation)

		payload := `{"id":"1","event":"login"}`
		raw := fmt.Sprintf(`{"%s":"%s","%s":"oplog","%s":"%s"}`,
			Meta, BIGQUERY, bigquery.TableNameKey, MsgBody, strings.ReplaceAll(payload, `"`, `\"`))

		n, err := ws.Write([]byte(raw))
		So(err, ShouldBeNil)
		So(n, ShouldEqual, len(raw))

		So(ws.Sync(), ShouldBeNil)
		So(ws.Close(), ShouldBeNil)

		fname := expectedTrackFilePath(dir, BIGQUERY+"_oplog", runtimeHostName(), HourlyRotation, time.Now())
		lines, readErr := readFileLines(fname)
		So(readErr, ShouldBeNil)
		So(lines, ShouldHaveLength, 1)
		So(lines[0], ShouldEqual, payload)
	})
}

func TestTrackFileWriteSyncer_ChannelRouting(t *testing.T) {
	Convey("不同 channel 写入不同文件", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncer(dir, HourlyRotation)

		writeChannel := func(channel, msg string) {
			raw := fmt.Sprintf(`{"%s":"%s","%s":"%s"}`, Meta, channel, MsgBody, msg)
			_, err := ws.Write([]byte(raw))
			So(err, ShouldBeNil)
		}

		writeChannel("thinkingdata", "tga_data")
		writeChannel("bi", "bi_data")
		writeChannel("bigquery", "bq_data")

		So(ws.Sync(), ShouldBeNil)
		So(ws.Close(), ShouldBeNil)

		for _, tc := range []struct{ channel, content string }{
			{"thinkingdata", "tga_data"},
			{"bi", "bi_data"},
			{"bigquery", "bq_data"},
		} {
			fname := expectedTrackFilePath(dir, tc.channel, runtimeHostName(), HourlyRotation, time.Now())
			lines, readErr := readFileLines(fname)
			So(readErr, ShouldBeNil)
			So(lines, ShouldHaveLength, 1)
			So(lines[0], ShouldEqual, tc.content)
		}
	})
}

func TestTrackFileWriteSyncer_SanitizeChannelFileName(t *testing.T) {
	Convey("channel 文件名会清洗路径分隔符，避免逃逸 TrackFileDir", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncer(dir, HourlyRotation)

		channel := "../evil/path"
		raw := fmt.Sprintf(`{"%s":"%s","%s":"safe_data"}`, Meta, channel, MsgBody)
		_, err := ws.Write([]byte(raw))
		So(err, ShouldBeNil)

		So(ws.Sync(), ShouldBeNil)
		So(ws.Close(), ShouldBeNil)

		fname := expectedTrackFilePath(dir, channel, runtimeHostName(), HourlyRotation, time.Now())
		// 必须仍位于 baseDir 内（路径分隔符已被清洗）
		rel, err := filepath.Rel(dir, fname)
		So(err, ShouldBeNil)
		So(strings.HasPrefix(rel, ".."), ShouldBeFalse)
		lines, readErr := readFileLines(fname)
		So(readErr, ShouldBeNil)
		So(lines, ShouldHaveLength, 1)
		So(lines[0], ShouldEqual, "safe_data")
	})
}

func TestTrackFileWriteSyncer_TimeSlot(t *testing.T) {
	Convey("timeSlot 格式正确", t, func() {
		now := time.Date(2026, 5, 9, 15, 23, 0, 0, time.UTC)
		So(HourlyRotation.timeSlot(now), ShouldEqual, "2026050915")
		So(MinuteRotation.timeSlot(now), ShouldEqual, "202605091523")
	})
}

func TestTrackFileOpenEvent(t *testing.T) {
	Convey("track 文件生命周期事件区分创建、复用和轮转", t, func() {
		So(trackFileOpenEvent(false, true), ShouldEqual, "created")
		So(trackFileOpenEvent(false, false), ShouldEqual, "opened")
		So(trackFileOpenEvent(true, true), ShouldEqual, "rotated")
		So(trackFileOpenEvent(true, false), ShouldEqual, "rotated")
	})
}

func TestTrackFileWriteSyncer_SyncReturnsWriteError(t *testing.T) {
	Convey("文件打开失败会在 Sync 返回错误", t, func() {
		dir := t.TempDir()
		baseFile := filepath.Join(dir, "not_a_dir")
		So(os.WriteFile(baseFile, []byte("x"), 0o644), ShouldBeNil)
		ws := NewTrackFileWriteSyncer(baseFile, HourlyRotation)

		raw := fmt.Sprintf(`{"%s":"bi","%s":"data"}`, Meta, MsgBody)
		_, err := ws.Write([]byte(raw))
		So(err, ShouldBeNil)
		So(ws.Sync(), ShouldNotBeNil)
		So(ws.Close(), ShouldBeNil)
	})
}

func TestTrackFileWriteSyncer_InvalidJSON(t *testing.T) {
	Convey("非 JSON 或缺失字段时静默忽略，不报错", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncer(dir, HourlyRotation)
		defer ws.Close()

		n, err := ws.Write([]byte("not valid json"))
		So(err, ShouldBeNil)
		So(n, ShouldEqual, len("not valid json"))

		n, err = ws.Write([]byte(`{"other":"field"}`))
		So(err, ShouldBeNil)
		So(n, ShouldBeGreaterThan, 0)
	})
}

func TestTrackFileWriteSyncer_CloseConcurrentWithWrite(t *testing.T) {
	Convey("Close 与 Write 并发不会 panic", t, func() {
		for iter := 0; iter < 20; iter++ {
			dir := t.TempDir()
			ws := NewTrackFileWriteSyncerWithBufSize(dir, HourlyRotation, 8)
			start := make(chan struct{})
			panicCh := make(chan interface{}, 1)
			var wg sync.WaitGroup

			for g := 0; g < 8; g++ {
				wg.Add(1)
				go func(id int) {
					defer wg.Done()
					defer func() {
						if r := recover(); r != nil {
							select {
							case panicCh <- r:
							default:
							}
						}
					}()
					<-start
					for i := 0; i < 100; i++ {
						raw := fmt.Sprintf(`{"%s":"bi","%s":"data_%d_%d"}`, Meta, MsgBody, id, i)
						_, _ = ws.Write([]byte(raw))
					}
				}(g)
			}

			close(start)
			_ = ws.Close()
			wg.Wait()
			select {
			case r := <-panicCh:
				t.Fatalf("Write panicked while Close was running: %v", r)
			default:
			}
		}
	})
}

func TestTrackFileWriteSyncer_ConcurrentWrite(t *testing.T) {
	Convey("多 goroutine 并发写入同一 channel，所有行均落盘", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncer(dir, HourlyRotation)

		const goroutines = 20
		const perGoroutine = 50

		var wg sync.WaitGroup
		for i := range goroutines {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := range perGoroutine {
					msg := fmt.Sprintf(`msg_%d_%d`, id, j)
					raw := fmt.Sprintf(`{"%s":"thinkingdata","%s":"%s"}`, Meta, MsgBody, msg)
					ws.Write([]byte(raw)) //nolint:errcheck
				}
			}(i)
		}
		wg.Wait()

		So(ws.Sync(), ShouldBeNil)
		So(ws.Close(), ShouldBeNil)

		fname := expectedTrackFilePath(dir, "thinkingdata", runtimeHostName(), HourlyRotation, time.Now())
		lines, readErr := readFileLines(fname)
		So(readErr, ShouldBeNil)
		// 所有消息均落盘（队列足够大，无丢弃）
		So(len(lines), ShouldEqual, goroutines*perGoroutine)
	})
}

func TestTrackFileWriteSyncer_ConcurrentMultiChannel(t *testing.T) {
	Convey("多 goroutine 并发写入多个 channel，各 channel 文件内容独立", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncer(dir, HourlyRotation)

		channels := []string{"thinkingdata", "bi", "bigquery"}
		const perChannel = 30

		var wg sync.WaitGroup
		for _, ch := range channels {
			wg.Add(1)
			go func(channel string) {
				defer wg.Done()
				for i := range perChannel {
					msg := fmt.Sprintf(`data_%s_%d`, channel, i)
					raw := fmt.Sprintf(`{"%s":"%s","%s":"%s"}`, Meta, channel, MsgBody, msg)
					ws.Write([]byte(raw)) //nolint:errcheck
				}
			}(ch)
		}
		wg.Wait()

		So(ws.Sync(), ShouldBeNil)
		So(ws.Close(), ShouldBeNil)

		for _, ch := range channels {
			fname := expectedTrackFilePath(dir, ch, runtimeHostName(), HourlyRotation, time.Now())
			lines, readErr := readFileLines(fname)
			So(readErr, ShouldBeNil)
			So(len(lines), ShouldEqual, perChannel)
			// 每行前缀正确
			for _, line := range lines {
				So(line, ShouldStartWith, "data_"+ch+"_")
			}
		}
	})
}

func TestTrackFileWriteSyncer_SyncFlushesQueue(t *testing.T) {
	Convey("Sync 调用后队列已排空，文件已 fsync", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncer(dir, HourlyRotation)

		const count = 100
		for i := range count {
			msg := fmt.Sprintf(`line_%d`, i)
			raw := fmt.Sprintf(`{"%s":"bi","%s":"%s"}`, Meta, MsgBody, msg)
			ws.Write([]byte(raw)) //nolint:errcheck
		}

		So(ws.Sync(), ShouldBeNil)

		fname := expectedTrackFilePath(dir, "bi", runtimeHostName(), HourlyRotation, time.Now())
		lines, readErr := readFileLines(fname)
		So(readErr, ShouldBeNil)
		So(len(lines), ShouldEqual, count)

		So(ws.Close(), ShouldBeNil)
	})
}

func TestTrackFileWriteSyncer_BatchFlushByTimer(t *testing.T) {
	Convey("未调用 Sync 时，worker 也会按定时器批量 flush", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncer(dir, HourlyRotation)
		defer ws.Close()

		raw := fmt.Sprintf(`{"%s":"bi","%s":"timer_flush"}`, Meta, MsgBody)
		_, err := ws.Write([]byte(raw))
		So(err, ShouldBeNil)

		fname := expectedTrackFilePath(dir, "bi", runtimeHostName(), HourlyRotation, time.Now())
		var lines []string
		deadline := time.Now().Add(500 * time.Millisecond)
		for time.Now().Before(deadline) {
			lines, err = readFileLines(fname)
			if err == nil && len(lines) == 1 {
				break
			}
			time.Sleep(5 * time.Millisecond)
		}
		So(err, ShouldBeNil)
		So(lines, ShouldHaveLength, 1)
		So(lines[0], ShouldEqual, "timer_flush")
	})
}

func TestTrackFileWriteSyncer_CloseIdempotent(t *testing.T) {
	Convey("Close 幂等，多次调用不 panic", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncer(dir, HourlyRotation)
		So(ws.Close(), ShouldBeNil)
		So(ws.Close(), ShouldBeNil)
	})
}

func TestTrackFileWriteSyncer_WriteAfterClose(t *testing.T) {
	Convey("Close 后写入静默丢弃，不报错", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncer(dir, HourlyRotation)
		So(ws.Close(), ShouldBeNil)

		raw := fmt.Sprintf(`{"%s":"bi","%s":"data"}`, Meta, MsgBody)
		n, err := ws.Write([]byte(raw))
		So(err, ShouldBeNil)
		So(n, ShouldEqual, len(raw))
	})
}

// ── 路径布局 / channel alias / keepalive ────────────────────────────────

func TestTrackFileWriteSyncer_PathLayout(t *testing.T) {
	Convey("路径布局符合 ops 规范 {base}/{channel}/{yyyymmdd}/{channel}_{hostname}_{slot}.log", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncer(dir, HourlyRotation)

		raw := fmt.Sprintf(`{"%s":"bi","%s":"layout_data"}`, Meta, MsgBody)
		_, err := ws.Write([]byte(raw))
		So(err, ShouldBeNil)
		So(ws.Sync(), ShouldBeNil)
		So(ws.Close(), ShouldBeNil)

		expected := expectedTrackFilePath(dir, "bi", runtimeHostName(), HourlyRotation, time.Now())
		_, statErr := os.Stat(expected)
		So(statErr, ShouldBeNil)
		lines, _ := readFileLines(expected)
		So(lines, ShouldHaveLength, 1)
		So(lines[0], ShouldEqual, "layout_data")
	})
}

func TestTrackFileWriteSyncer_ChannelAlias(t *testing.T) {
	Convey("ChannelAlias 重写目录段与文件名前缀", t, func() {
		base := t.TempDir()
		ws := NewTrackFileWriteSyncerFromConfig(TrackFileWriteSyncerConfig{
			BaseDir:      base,
			Rotation:     HourlyRotation,
			ChannelAlias: map[string]string{"thinkingdata": "thinkdata"},
		})

		write := func(channel, extra, msg string) {
			raw := fmt.Sprintf(`{"%s":"%s",%s"%s":"%s"}`, Meta, channel, extra, MsgBody, msg)
			_, err := ws.Write([]byte(raw))
			So(err, ShouldBeNil)
		}

		write("thinkingdata", ``, "tga_data")
		write("bi", ``, "bi_data")
		write(BIGQUERY, fmt.Sprintf(`"%s":"oplog",`, bigquery.TableNameKey), "bq_data")

		So(ws.Sync(), ShouldBeNil)
		So(ws.Close(), ShouldBeNil)

		now := time.Now()
		check := func(channel, content string) {
			fname := expectedTrackFilePath(base, channel, runtimeHostName(), HourlyRotation, now)
			lines, err := readFileLines(fname)
			So(err, ShouldBeNil)
			So(lines, ShouldHaveLength, 1)
			So(lines[0], ShouldEqual, content)
		}
		// thinkingdata 被映射成 thinkdata
		check("thinkdata", "tga_data")
		// 其他 channel 未受影响
		check("bi", "bi_data")
		check(BIGQUERY+"_oplog", "bq_data")
	})
}

func TestTrackFileWriteSyncer_HostnameAuto(t *testing.T) {
	Convey("hostname 段自动取 os.Hostname", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncer(dir, HourlyRotation)
		raw := fmt.Sprintf(`{"%s":"bi","%s":"hn"}`, Meta, MsgBody)
		_, err := ws.Write([]byte(raw))
		So(err, ShouldBeNil)
		So(ws.Sync(), ShouldBeNil)
		So(ws.Close(), ShouldBeNil)

		fname := expectedTrackFilePath(dir, "bi", runtimeHostName(), HourlyRotation, time.Now())
		_, statErr := os.Stat(fname)
		So(statErr, ShouldBeNil)
	})
}

func TestTrackFileWriteSyncer_KeepaliveOnIdle(t *testing.T) {
	Convey("KeepaliveEvery>0 时闲置超时会落一条 keepalive 日志", t, func() {
		dir := t.TempDir()
		interval := 50 * time.Millisecond
		ws := NewTrackFileWriteSyncerFromConfig(TrackFileWriteSyncerConfig{
			BaseDir:          dir,
			Rotation:         HourlyRotation,
			KeepaliveEvery:   interval,
			KeepaliveMessage: `{"hb":1}`,
		})

		// 触发一次正常写入以创建 worker（未激活 channel 不会自动心跳）
		raw := fmt.Sprintf(`{"%s":"bi","%s":"first"}`, Meta, MsgBody)
		_, err := ws.Write([]byte(raw))
		So(err, ShouldBeNil)
		So(ws.Sync(), ShouldBeNil)

		// 等待至少 2 个间隔，期望 worker 至少补一条心跳
		time.Sleep(interval * 4)
		So(ws.Sync(), ShouldBeNil)
		So(ws.Close(), ShouldBeNil)

		fname := expectedTrackFilePath(dir, "bi", runtimeHostName(), HourlyRotation, time.Now())
		lines, readErr := readFileLines(fname)
		So(readErr, ShouldBeNil)
		So(len(lines), ShouldBeGreaterThanOrEqualTo, 2)
		So(lines[0], ShouldEqual, "first")
		// 后续行应为 keepalive payload
		for _, line := range lines[1:] {
			So(line, ShouldEqual, `{"hb":1}`)
		}
	})
}

func TestTrackFileWriteSyncer_KeepaliveDisabled(t *testing.T) {
	Convey("KeepaliveEvery<=0 时不会有心跳输出", t, func() {
		dir := t.TempDir()
		ws := NewTrackFileWriteSyncerFromConfig(TrackFileWriteSyncerConfig{
			BaseDir:        dir,
			Rotation:       HourlyRotation,
			KeepaliveEvery: 0, // 禁用
		})

		raw := fmt.Sprintf(`{"%s":"bi","%s":"only_one"}`, Meta, MsgBody)
		_, err := ws.Write([]byte(raw))
		So(err, ShouldBeNil)
		So(ws.Sync(), ShouldBeNil)

		time.Sleep(80 * time.Millisecond)
		So(ws.Sync(), ShouldBeNil)
		So(ws.Close(), ShouldBeNil)

		fname := expectedTrackFilePath(dir, "bi", runtimeHostName(), HourlyRotation, time.Now())
		lines, readErr := readFileLines(fname)
		So(readErr, ShouldBeNil)
		So(lines, ShouldHaveLength, 1)
	})
}

func TestTrackFileWriteSyncer_KeepaliveResetByWrite(t *testing.T) {
	Convey("正常写入会重置心跳计时，活跃 channel 不会触发心跳", t, func() {
		dir := t.TempDir()
		interval := 80 * time.Millisecond
		ws := NewTrackFileWriteSyncerFromConfig(TrackFileWriteSyncerConfig{
			BaseDir:          dir,
			Rotation:         HourlyRotation,
			KeepaliveEvery:   interval,
			KeepaliveMessage: `KEEP`,
		})

		// 持续写入 ~3 个间隔，期间不应产生心跳
		deadline := time.Now().Add(interval * 3)
		count := 0
		for time.Now().Before(deadline) {
			raw := fmt.Sprintf(`{"%s":"bi","%s":"busy_%d"}`, Meta, MsgBody, count)
			_, err := ws.Write([]byte(raw))
			So(err, ShouldBeNil)
			count++
			time.Sleep(interval / 4)
		}
		So(ws.Sync(), ShouldBeNil)
		So(ws.Close(), ShouldBeNil)

		fname := expectedTrackFilePath(dir, "bi", runtimeHostName(), HourlyRotation, time.Now())
		lines, readErr := readFileLines(fname)
		So(readErr, ShouldBeNil)
		for _, line := range lines {
			So(line, ShouldNotEqual, "KEEP")
		}
	})
}

// ── LevelEnabler 单元测试 ──────────────────────────────────────────────────

func TestTrackOnlyLevelEnabler(t *testing.T) {
	Convey("trackOnlyLevelEnabler 只在当前启用文件输出时允许 TrackLevel", t, func() {
		e := trackOnlyLevelEnabler{}
		oldWriter, oldOutput := gTrackFileWriteSyncerProxy.state()
		defer func() {
			_ = gTrackFileWriteSyncerProxy.setCurrent(oldWriter, oldOutput)
			runtimeEnableTrackLevel.Store(true)
		}()

		_ = gTrackFileWriteSyncerProxy.setCurrent(&TrackFileWriteSyncer{}, TrackOutputFile)
		runtimeEnableTrackLevel.Store(true)
		So(e.Enabled(TrackLevel), ShouldBeTrue)
		So(e.Enabled(0), ShouldBeFalse) // DebugLevel
		So(e.Enabled(1), ShouldBeFalse) // InfoLevel

		_ = gTrackFileWriteSyncerProxy.setCurrent(&TrackFileWriteSyncer{}, TrackOutputStdout)
		So(e.Enabled(TrackLevel), ShouldBeFalse)

		_ = gTrackFileWriteSyncerProxy.setCurrent(nil, TrackOutputFile)
		So(e.Enabled(TrackLevel), ShouldBeFalse)

		_ = gTrackFileWriteSyncerProxy.setCurrent(&TrackFileWriteSyncer{}, TrackOutputFile)
		runtimeEnableTrackLevel.Store(false)
		So(e.Enabled(TrackLevel), ShouldBeFalse)
	})
}

func TestTrackStdoutLevelEnabler(t *testing.T) {
	Convey("trackStdoutLevelEnabler 按当前 TrackOutput 路由 TrackLevel", t, func() {
		base := newTrackLevelEnabler()
		e := &trackStdoutLevelEnabler{wrapped: base}
		oldWriter, oldOutput := gTrackFileWriteSyncerProxy.state()
		defer func() { _ = gTrackFileWriteSyncerProxy.setCurrent(oldWriter, oldOutput) }()

		_ = gTrackFileWriteSyncerProxy.setCurrent(&TrackFileWriteSyncer{}, TrackOutputFile)
		So(e.Enabled(TrackLevel), ShouldBeFalse)
		_ = gTrackFileWriteSyncerProxy.setCurrent(nil, TrackOutputStdout)
		So(e.Enabled(TrackLevel), ShouldBeTrue)
		_ = gTrackFileWriteSyncerProxy.setCurrent(&TrackFileWriteSyncer{}, TrackOutputBoth)
		So(e.Enabled(TrackLevel), ShouldBeTrue)
		So(e.Enabled(0), ShouldBeTrue) // DebugLevel 在 runtimeLogLevel=Debug 下通过
	})
}
