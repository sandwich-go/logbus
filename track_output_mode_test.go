package logbus

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap/zapcore"
)

// initWithBuf 用 bytes.Buffer 替换 WriteSyncer，方便捕获 stdout 输出。
// 返回 buf 和 cleanup 函数（恢复默认配置）。
func initWithBuf(extraOpts ...ConfOption) (*bytes.Buffer, func()) {
	buf := &bytes.Buffer{}
	opts := append([]ConfOption{WithWriteSyncer(zapcore.AddSync(buf))}, extraOpts...)
	Init(NewConf(opts...))
	return buf, resetLogBus
}

// writeTrackLog 直接调 gStdLogger 写一条 track 日志，绕过 Tracker 的序列化层，
// 专注于测试 LevelEnabler 路由和文件写入行为。
func writeTrackLog(channel, msg string) {
	gStdLogger.TrackWithChannel(channel, String(MsgBody, msg))
}

// writeNormalLog 写一条普通 info 日志
func writeNormalLog(msg string) {
	gStdLogger.InfoWithChannel(Setting.DefaultChannel, String(MsgBody, msg))
}

// syncAndRead 等待 gStdLogger 所有 core（含 TrackFileWriteSyncer）flush 完成，
// 然后读取 channel 对应的文件内容（非空行列表）。
// 路径按 ops 规范：{dir}/{channel}/{yyyymmdd}/{channel}_{hostname}_{slot}.log
// channel 自动应用全局配置的文件落盘别名（例如 thinkingdata→tga），仅用于拼读文件路径。
func syncAndRead(dir, channel string) ([]string, error) {
	_ = gStdLogger.Sync() // 触发 Tee → TrackFileWriteSyncer.Sync() → worker drain+fsync
	if alias, ok := Setting.TrackChannelAlias[channel]; ok && alias != "" {
		channel = alias
	}
	fname := expectedTrackFilePath(dir, channel, runtimeHostName(), HourlyRotation, time.Now())
	data, err := os.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	var lines []string
	for _, line := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	return lines, nil
}

// ── TrackOutputStdout（默认，向后兼容）─────────────────────────────────────

func TestTrackOutput_Stdout_Default(t *testing.T) {
	Convey("TrackOutputStdout（默认）: track 日志写 stdout，行为与旧版完全一致", t, func() {
		buf, cleanup := initWithBuf() // 不传 WithTrackOutput，使用默认值
		defer cleanup()

		writeTrackLog(THINKINGDATA, `{"#account_id":"111"}`)
		_ = gStdLogger.Sync()

		out := buf.String()
		// stdout 包含外层 JSON：dd_meta_channel=thinkingdata，msg 内容被 JSON 转义
		So(out, ShouldContainSubstring, THINKINGDATA)
		So(out, ShouldContainSubstring, `#account_id`) // 转义后仍可搜到 key
	})
}

func TestTrackOutput_Stdout_NormalLog(t *testing.T) {
	Convey("TrackOutputStdout: 普通日志正常写 stdout", t, func() {
		buf, cleanup := initWithBuf()
		defer cleanup()

		writeNormalLog("hello_stdout")
		_ = gStdLogger.Sync()

		So(buf.String(), ShouldContainSubstring, "hello_stdout")
	})
}

// ── TrackOutputFile（只写文件）────────────────────────────────────────────

func TestTrackOutput_FileOnly_TrackNotInStdout(t *testing.T) {
	Convey("TrackOutputFile: track 日志不出现在 stdout", t, func() {
		dir := t.TempDir()
		buf, cleanup := initWithBuf(
			WithTrackOutput(TrackOutputFile),
			WithTrackFileDir(dir),
		)
		defer cleanup()

		writeTrackLog(THINKINGDATA, `{"#account_id":"file_only"}`)
		_ = gStdLogger.Sync()

		So(buf.String(), ShouldNotContainSubstring, THINKINGDATA)
		So(buf.String(), ShouldNotContainSubstring, "file_only")
	})
}

func TestTrackOutput_FileOnly_TrackInFile(t *testing.T) {
	Convey("TrackOutputFile: track 日志写入对应 channel 文件，只含 msg 内容", t, func() {
		dir := t.TempDir()
		_, cleanup := initWithBuf(
			WithTrackOutput(TrackOutputFile),
			WithTrackFileDir(dir),
		)
		defer cleanup()

		payload := `{"#account_id":"222","#type":"track"}`
		writeTrackLog(THINKINGDATA, payload)

		lines, err := syncAndRead(dir, THINKINGDATA)
		So(err, ShouldBeNil)
		So(lines, ShouldHaveLength, 1)
		So(lines[0], ShouldEqual, payload)
	})
}

func TestTrackOutput_FileOnly_NormalLogStillInStdout(t *testing.T) {
	Convey("TrackOutputFile: 普通日志仍写 stdout，不受 track 路由影响", t, func() {
		dir := t.TempDir()
		buf, cleanup := initWithBuf(
			WithTrackOutput(TrackOutputFile),
			WithTrackFileDir(dir),
		)
		defer cleanup()

		writeNormalLog("normal_file_mode")
		_ = gStdLogger.Sync()

		So(buf.String(), ShouldContainSubstring, "normal_file_mode")
	})
}

func TestTrackOutput_FileOnly_MultiChannel(t *testing.T) {
	Convey("TrackOutputFile: 多 channel 各写各的文件", t, func() {
		dir := t.TempDir()
		_, cleanup := initWithBuf(
			WithTrackOutput(TrackOutputFile),
			WithTrackFileDir(dir),
		)
		defer cleanup()

		writeTrackLog(THINKINGDATA, `tga_msg`)
		writeTrackLog(BI, `bi_msg`)

		// 统一 sync 一次即可
		_ = gStdLogger.Sync()

		tgaLines, err := syncAndRead(dir, THINKINGDATA)
		So(err, ShouldBeNil)
		So(tgaLines, ShouldHaveLength, 1)
		So(tgaLines[0], ShouldEqual, `tga_msg`)

		biLines, err := syncAndRead(dir, BI)
		So(err, ShouldBeNil)
		So(biLines, ShouldHaveLength, 1)
		So(biLines[0], ShouldEqual, `bi_msg`)
	})
}

// ── TrackOutputBoth（同时写）─────────────────────────────────────────────

func TestTrackOutput_Both_TrackInStdout(t *testing.T) {
	Convey("TrackOutputBoth: track 日志写 stdout（完整外层 JSON）", t, func() {
		dir := t.TempDir()
		buf, cleanup := initWithBuf(
			WithTrackOutput(TrackOutputBoth),
			WithTrackFileDir(dir),
		)
		defer cleanup()

		writeTrackLog(BI, `{"app_id":"gof"}`)
		_ = gStdLogger.Sync()

		out := buf.String()
		// stdout 包含外层 JSON，dd_meta_channel=bi，msg 内容转义
		So(out, ShouldContainSubstring, BI)
		So(out, ShouldContainSubstring, `app_id`) // 转义后仍可搜到 key
	})
}

func TestTrackOutput_Both_TrackInFile(t *testing.T) {
	Convey("TrackOutputBoth: track 日志同时写文件，文件只含 msg 内容", t, func() {
		dir := t.TempDir()
		_, cleanup := initWithBuf(
			WithTrackOutput(TrackOutputBoth),
			WithTrackFileDir(dir),
		)
		defer cleanup()

		payload := `{"app_id":"gof","event":"login"}`
		writeTrackLog(BI, payload)

		lines, err := syncAndRead(dir, BI)
		So(err, ShouldBeNil)
		So(lines, ShouldHaveLength, 1)
		So(lines[0], ShouldEqual, payload)
	})
}

func TestTrackOutput_Both_NormalLogInStdout(t *testing.T) {
	Convey("TrackOutputBoth: 普通日志仍写 stdout", t, func() {
		dir := t.TempDir()
		buf, cleanup := initWithBuf(
			WithTrackOutput(TrackOutputBoth),
			WithTrackFileDir(dir),
		)
		defer cleanup()

		writeNormalLog("both_normal")
		_ = gStdLogger.Sync()

		So(buf.String(), ShouldContainSubstring, "both_normal")
	})
}

// ── watchdog 配置校验 ─────────────────────────────────────────────────────

func TestTrackOutput_WatchdogPanicsWhenDirEmpty(t *testing.T) {
	Convey("FileOnly/Both 但 TrackFileDir 为空时 Init 应 panic", t, func() {
		for _, mode := range []TrackOutput{TrackOutputFile, TrackOutputBoth} {
			So(func() {
				Init(NewConf(WithTrackOutput(mode), WithTrackFileDir("")))
			}, ShouldPanic)
		}
		// panic 后 logbus 状态可能不一致，重新初始化
		Init(NewConf())
		resetLogBus()
	})
}

// ── PMT 环境变量覆写校验 ──────────────────────────────────────────────────

func TestTrackOutput_PMTEnvOverrideAppliesValues(t *testing.T) {
	Convey("PMT 环境变量覆写 TrackFileDir / TrackOutput", t, func() {
		dir := t.TempDir()
		t.Setenv("sys_cd_env", "prod")
		t.Setenv("logbus_track_output", "file")
		t.Setenv("logbus_track_file_dir", dir)

		// 代码侧显式启用文件输出时，PMT 环变覆写为 File + dir。
		Init(NewConf(
			WithTrackOutput(TrackOutputFile),
			WithTrackFileDir("/tmp/should_be_overwritten"),
		))
		defer resetLogBus()

		So(Setting.TrackOutput, ShouldEqual, TrackOutputFile)
		So(Setting.TrackFileDir, ShouldEqual, dir)
	})
}

func TestTrackOutput_PMTEnvIgnoredWhenCodeStdout(t *testing.T) {
	Convey("PMT 环境下代码侧保持 stdout 时不启用文件输出", t, func() {
		dir := t.TempDir()
		t.Setenv("sys_cd_env", "prod")
		t.Setenv("logbus_track_output", "file")
		t.Setenv("logbus_track_file_dir", dir)

		Init(NewConf())
		defer resetLogBus()

		So(Setting.TrackOutput, ShouldEqual, TrackOutputStdout)
		So(Setting.TrackFileDir, ShouldEqual, "./tracklogs")
	})
}

// 注意：Convey 用例结尾的 Init(NewConf())+resetLogBus() 用于把全局状态恢复到默认；
// 如果让 PMT 环变在恢复期间仍生效，恢复 Init 自身就会 panic。所以 panic 用例必须在
// 断言完成后立即清空相关 env，再做最终恢复。
func TestTrackOutput_PMTEnvOverrideRejectsBadOutput(t *testing.T) {
	Convey("PMT 环境下 logbus_track_output 非法值时 Init panic", t, func() {
		t.Setenv("sys_cd_env", "prod")
		t.Setenv("logbus_track_output", "garbage")
		So(func() { Init(NewConf(WithTrackOutput(TrackOutputFile))) }, ShouldPanic)
		t.Setenv("sys_cd_env", "")
		Init(NewConf())
		resetLogBus()
	})
}

func TestTrackOutput_PMTEnvOverrideRequiresOutput(t *testing.T) {
	Convey("PMT 环境下未设置 logbus_track_output 时 Init panic", t, func() {
		t.Setenv("sys_cd_env", "prod")
		// 显式确保 logbus_track_output 不存在
		t.Setenv("logbus_track_output", "")
		So(func() { Init(NewConf(WithTrackOutput(TrackOutputFile))) }, ShouldPanic)
		t.Setenv("sys_cd_env", "")
		Init(NewConf())
		resetLogBus()
	})
}

func TestTrackOutput_PMTEnvOverrideRequiresDirWhenFile(t *testing.T) {
	Convey("PMT 环境 + output=file 但 logbus_track_file_dir 缺失时 panic", t, func() {
		t.Setenv("sys_cd_env", "prod")
		t.Setenv("logbus_track_output", "file")
		t.Setenv("logbus_track_file_dir", "")
		So(func() { Init(NewConf(WithTrackOutput(TrackOutputFile))) }, ShouldPanic)
		t.Setenv("sys_cd_env", "")
		Init(NewConf())
		resetLogBus()
	})
}

func TestTrackOutput_PMTEnvOverrideStdoutNoDirRequired(t *testing.T) {
	Convey("PMT 环境 + output=stdout 时不要求 dir，且 TrackOutput 被覆写", t, func() {
		t.Setenv("sys_cd_env", "prod")
		t.Setenv("logbus_track_output", "stdout")
		t.Setenv("logbus_track_file_dir", "")

		Init(NewConf(WithTrackOutput(TrackOutputFile), WithTrackFileDir("/tmp/should_be_overwritten")))
		defer resetLogBus()
		So(Setting.TrackOutput, ShouldEqual, TrackOutputStdout)
	})
}

func TestTrackOutput_PMTEnvOverrideKeepalive(t *testing.T) {
	Convey("PMT 环境下 keepalive 环变覆写 interval 与 message", t, func() {
		dir := t.TempDir()
		t.Setenv("sys_cd_env", "prod")
		t.Setenv("logbus_track_output", "file")
		t.Setenv("logbus_track_file_dir", dir)
		t.Setenv("logbus_track_keepalive_interval", "30s")
		t.Setenv("logbus_track_keepalive_message", `{"hb":1}`)

		// 代码侧给的默认值都应被环变覆盖
		Init(NewConf(
			WithTrackOutput(TrackOutputFile),
			WithTrackKeepaliveInterval(time.Minute),
			WithTrackKeepaliveMessage("code-default"),
		))
		defer resetLogBus()

		So(Setting.TrackKeepaliveInterval, ShouldEqual, 30*time.Second)
		So(Setting.TrackKeepaliveMessage, ShouldEqual, `{"hb":1}`)
	})
}

func TestTrackOutput_PMTEnvOverrideKeepaliveZeroDisables(t *testing.T) {
	Convey("PMT 环境下 keepalive_interval=0 表示禁用", t, func() {
		dir := t.TempDir()
		t.Setenv("sys_cd_env", "prod")
		t.Setenv("logbus_track_output", "file")
		t.Setenv("logbus_track_file_dir", dir)
		t.Setenv("logbus_track_keepalive_interval", "0")

		// 即便代码侧默认 1 分钟，也应被覆盖为 0
		Init(NewConf(
			WithTrackOutput(TrackOutputFile),
			WithTrackKeepaliveInterval(time.Minute),
		))
		defer resetLogBus()
		So(Setting.TrackKeepaliveInterval, ShouldEqual, time.Duration(0))
	})
}

func TestTrackOutput_PMTEnvOverrideKeepaliveEmptyMessage(t *testing.T) {
	Convey("PMT 环境下显式空 keepalive_message 应覆盖代码侧默认", t, func() {
		dir := t.TempDir()
		t.Setenv("sys_cd_env", "prod")
		t.Setenv("logbus_track_output", "file")
		t.Setenv("logbus_track_file_dir", dir)
		t.Setenv("logbus_track_keepalive_message", "") // 显式空串

		Init(NewConf(
			WithTrackOutput(TrackOutputFile),
			WithTrackKeepaliveMessage("code-default"),
		))
		defer resetLogBus()
		So(Setting.TrackKeepaliveMessage, ShouldEqual, "")
	})
}

func TestTrackOutput_PMTEnvOverrideKeepaliveBadInterval(t *testing.T) {
	Convey("PMT 环境下 keepalive_interval 解析失败应 panic", t, func() {
		dir := t.TempDir()
		t.Setenv("sys_cd_env", "prod")
		t.Setenv("logbus_track_output", "file")
		t.Setenv("logbus_track_file_dir", dir)
		t.Setenv("logbus_track_keepalive_interval", "not-a-duration")

		So(func() { Init(NewConf(WithTrackOutput(TrackOutputFile))) }, ShouldPanic)
		t.Setenv("sys_cd_env", "")
		Init(NewConf())
		resetLogBus()
	})
}

func TestTrackOutput_PMTEnvOverrideKeepaliveNegativeInterval(t *testing.T) {
	Convey("PMT 环境下 keepalive_interval 负值应 panic", t, func() {
		dir := t.TempDir()
		t.Setenv("sys_cd_env", "prod")
		t.Setenv("logbus_track_output", "file")
		t.Setenv("logbus_track_file_dir", dir)
		t.Setenv("logbus_track_keepalive_interval", "-1s")

		So(func() { Init(NewConf(WithTrackOutput(TrackOutputFile))) }, ShouldPanic)
		t.Setenv("sys_cd_env", "")
		Init(NewConf())
		resetLogBus()
	})
}

// ── 全局单例 + 生命周期 ───────────────────────────────────────────────────

// TestTrackOutput_ScopeLoggersShareSameWriter: 多个 ScopeLogger 写同一 track channel
// 时必须落到同一文件，且写入的所有行都必须完整、顺序可被解析（无交错损坏）。
// 修复前每次 newNLoggerInstance 都会新建一个 TrackFileWriteSyncer，多 worker 并发
// append 同一文件会出现行交错；修复后通过 gTrackFileWriteSyncer 单例消除此风险。
func TestTrackOutput_ScopeLoggersShareSameWriter(t *testing.T) {
	Convey("多个 ScopeLogger 共享同一 TrackFileWriteSyncer，写入互不交错", t, func() {
		dir := t.TempDir()
		_, cleanup := initWithBuf(
			WithTrackOutput(TrackOutputFile),
			WithTrackFileDir(dir),
		)
		defer cleanup()

		// 所有 scope logger 必须引用同一个全局 writer
		So(gTrackFileWriteSyncer, ShouldNotBeNil)
		singleton := gTrackFileWriteSyncer

		const scopes = 4
		const perScope = 50
		loggers := make([]NewILogger, 0, scopes)
		for i := 0; i < scopes; i++ {
			loggers = append(loggers, NewScopeLogger(fmt.Sprintf("scope_%d", i)))
		}
		// 再次断言：ScopeLogger 创建过程中没有替换掉全局单例
		So(gTrackFileWriteSyncer, ShouldEqual, singleton)

		var wg sync.WaitGroup
		for si, lg := range loggers {
			wg.Add(1)
			go func(scopeIdx int, logger NewILogger) {
				defer wg.Done()
				for j := 0; j < perScope; j++ {
					// 使用不同的 msg payload，便于后续验证完整性
					payload := fmt.Sprintf("scope_%d_line_%d", scopeIdx, j)
					logger.(*GLogger).GetStdLogger().TrackWithChannel(THINKINGDATA, String(MsgBody, payload))
				}
			}(si, lg)
		}
		wg.Wait()

		lines, err := syncAndRead(dir, THINKINGDATA)
		So(err, ShouldBeNil)
		So(len(lines), ShouldEqual, scopes*perScope)
		// 每行必须是完整的 payload（无交错破坏）
		for _, line := range lines {
			So(line, ShouldStartWith, "scope_")
			So(line, ShouldContainSubstring, "_line_")
		}
	})
}

// TestTrackOutput_OldScopeLoggerWritesAfterReinit 覆盖热切换同源问题：
// 旧 ScopeLogger 的 file core 捕获的是稳定代理，二次 Init 重建具体 writer 后，
// 旧 logger 继续写 track 时必须进入新的 writer，而不是写入已关闭旧 writer 后静默丢弃。
func TestTrackOutput_OldScopeLoggerWritesAfterReinit(t *testing.T) {
	Convey("旧 ScopeLogger 在二次 Init 后继续写 track 不丢失", t, func() {
		firstDir := t.TempDir()
		_, cleanup := initWithBuf(
			WithTrackOutput(TrackOutputFile),
			WithTrackFileDir(firstDir),
		)
		defer cleanup()

		oldScope := NewScopeLogger("old_scope")
		oldStd := oldScope.(*GLogger).GetStdLogger()
		oldGlobal := gStdLogger

		oldStd.TrackWithChannel(BI, String(MsgBody, "before_reinit_scope"))
		oldGlobal.TrackWithChannel(BI, String(MsgBody, "before_reinit_global"))
		So(oldStd.Sync(), ShouldBeNil)

		firstLines, err := syncAndRead(firstDir, BI)
		So(err, ShouldBeNil)
		So(firstLines, ShouldContain, "before_reinit_scope")
		So(firstLines, ShouldContain, "before_reinit_global")

		secondDir := t.TempDir()
		Init(NewConf(
			WithWriteSyncer(zapcore.AddSync(&bytes.Buffer{})),
			WithTrackOutput(TrackOutputFile),
			WithTrackFileDir(secondDir),
		))

		oldStd.TrackWithChannel(BI, String(MsgBody, "after_reinit_scope"))
		oldGlobal.TrackWithChannel(BI, String(MsgBody, "after_reinit_global"))
		So(oldStd.Sync(), ShouldBeNil)
		So(oldGlobal.Sync(), ShouldBeNil)

		secondLines, err := syncAndRead(secondDir, BI)
		So(err, ShouldBeNil)
		So(secondLines, ShouldContain, "after_reinit_scope")
		So(secondLines, ShouldContain, "after_reinit_global")
	})
}

func TestTrackOutput_OldStdoutLoggerWritesFileAfterReinit(t *testing.T) {
	Convey("stdout 模式创建的旧 ScopeLogger 在二次 Init 到 file 后可写文件", t, func() {
		buf, cleanup := initWithBuf(WithTrackOutput(TrackOutputStdout))
		defer cleanup()

		oldScope := NewScopeLogger("old_stdout_scope")
		oldStd := oldScope.(*GLogger).GetStdLogger()
		oldStd.TrackWithChannel(BI, String(MsgBody, "before_file_mode"))
		So(oldStd.Sync(), ShouldBeNil)
		So(buf.String(), ShouldContainSubstring, "before_file_mode")

		secondDir := t.TempDir()
		Init(NewConf(
			WithWriteSyncer(zapcore.AddSync(&bytes.Buffer{})),
			WithTrackOutput(TrackOutputFile),
			WithTrackFileDir(secondDir),
		))

		oldStd.TrackWithChannel(BI, String(MsgBody, "after_file_mode"))
		So(oldStd.Sync(), ShouldBeNil)
		So(buf.String(), ShouldNotContainSubstring, "after_file_mode")

		lines, err := syncAndRead(secondDir, BI)
		So(err, ShouldBeNil)
		So(lines, ShouldContain, "after_file_mode")
	})
}

// TestTrackOutput_InitReleasesOldWriter: 重复 Init 切换输出模式时，
// 旧的 TrackFileWriteSyncer 必须被 Close（释放 worker goroutine 与文件句柄）。
// 修复前每次 newNLoggerInstance 都新建 writer 且无任何地方关闭，重复 Init 会泄漏 goroutine。
func TestTrackOutput_InitReleasesOldWriter(t *testing.T) {
	Convey("重复 Init 时旧的 TrackFileWriteSyncer 必须被关闭", t, func() {
		dir := t.TempDir()

		// 首次 Init 为 File 模式，写一条 track 日志触发 worker 创建
		_, cleanup := initWithBuf(
			WithTrackOutput(TrackOutputFile),
			WithTrackFileDir(dir),
		)
		first := gTrackFileWriteSyncer
		So(first, ShouldNotBeNil)

		writeTrackLog(BI, "first_msg")
		_ = gStdLogger.Sync()

		// 触发 Init 第二次（不同模式），旧 writer 应被 Close
		cleanup() // 先把 buf-based syncer 场景复位
		_, cleanup2 := initWithBuf(WithTrackOutput(TrackOutputStdout))
		defer cleanup2()

		// TrackOutputStdout 模式下不应再有全局单例
		So(gTrackFileWriteSyncer, ShouldBeNil)
		// 旧 writer 的 Close 幂等且已完成：再次调用应立即返回 nil
		So(first.Close(), ShouldBeNil)
		// Close 后再 Write 静默丢弃
		n, err := first.Write([]byte(`{"dd_meta_channel":"bi","msg":"after_close"}`))
		So(err, ShouldBeNil)
		So(n, ShouldBeGreaterThan, 0)
	})
}

// TestDefaultStripFields_ReturnsCopy: DefaultStripFields() 必须返回副本，
// 外部 append 不应污染后续调用的返回值。
func TestDefaultStripFields_ReturnsCopy(t *testing.T) {
	Convey("DefaultStripFields 返回副本，外部 append 不污染全局", t, func() {
		a := DefaultStripFields()
		originalLen := len(a)
		// 尝试通过 append 扩展并改写（若容量足够会踩到底层数组）
		_ = append(a, "extra.field")
		// 尝试直接改元素
		if len(a) > 0 {
			a[0] = "mutated"
		}

		b := DefaultStripFields()
		So(len(b), ShouldEqual, originalLen)
		So(b[0], ShouldEqual, "api_call.request_body")
		So(b[1], ShouldEqual, "api_call.response_body")
	})
}
