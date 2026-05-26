package logbus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/sandwich-go/logbus/bigquery"
	"github.com/sandwich-go/logbus/monitor"
	"go.uber.org/zap/zapcore"
)

// TrackRotation 时间切割粒度
type TrackRotation int

const (
	HourlyRotation TrackRotation = iota // 按小时切割，文件名后缀格式 2006010215
	MinuteRotation                      // 按分钟切割，文件名后缀格式 200601021504
)

// MetricTrackDropped 队列满时丢弃日志的 counter 指标名
const MetricTrackDropped = "logbus_track_dropped"

// MetricTrackWriteFailed 文件打开、写入或 fsync 失败的 counter 指标名
const MetricTrackWriteFailed = "logbus_track_write_failed"

// TrackOutput 控制 track 日志的输出目标
type TrackOutput int

const (
	// TrackOutputStdout track 日志只写 stdout（默认，与引入文件写入前兼容）
	TrackOutputStdout TrackOutput = iota
	// TrackOutputFile track 日志只写文件，不写 stdout
	// 需同时配置 TrackFileDir
	TrackOutputFile
	// TrackOutputBoth track 日志同时写 stdout 和文件
	// 需同时配置 TrackFileDir
	TrackOutputBoth
)

func (r TrackRotation) timeSlot(t time.Time) string {
	switch r {
	case MinuteRotation:
		return t.Format("200601021504")
	default:
		return t.Format("2006010215")
	}
}

// defaultTrackChannelBufSize 每个 channel worker 的 Go channel 缓冲大小
const defaultTrackChannelBufSize = 4096

const (
	defaultTrackBatchMaxBytes      = 64 * 1024
	defaultTrackBatchMaxCount      = 256
	defaultTrackBatchFlushInterval = 10 * time.Millisecond
)

// trackMsg 是投递给 worker 的写入请求
type trackMsg struct {
	data []byte // 已提取的 msg 内容（不含换行）
}

// trackWorker 是每个 channel 独享的写入 goroutine
// 它持有一个文件句柄，串行完成文件切割与写入，无需与其他 channel 竞争。
type trackWorker struct {
	channel        string // sanitize 后的 channel 名（同时作为目录段与文件名前缀）
	baseDir        string
	hostName       string // 文件名 pod 段，自动取 hostname
	rotation       TrackRotation
	keepaliveEvery time.Duration // 0 表示禁用
	keepaliveMsg   []byte        // 心跳写入 payload；可为空字符串

	queue chan trackMsg   // 消息队列
	sync  chan chan error // sync 请求：调用方传入一个 chan error，worker 写完队列后回复
	done  chan struct{}   // worker 退出信号

	// 当前打开的文件句柄
	file *os.File
	slot string

	batch      []byte
	batchCount int
	lastErr    error
}

func newTrackWorker(channel, baseDir, hostName string, rotation TrackRotation, bufSize int, keepaliveEvery time.Duration, keepaliveMsg []byte) *trackWorker {
	w := &trackWorker{
		channel:        channel,
		baseDir:        baseDir,
		hostName:       hostName,
		rotation:       rotation,
		keepaliveEvery: keepaliveEvery,
		keepaliveMsg:   keepaliveMsg,
		queue:          make(chan trackMsg, bufSize),
		sync:           make(chan chan error, 1),
		done:           make(chan struct{}),
	}
	go w.run()
	return w
}

func (w *trackWorker) run() {
	// flushTimer 控制 batch 攒到一定时间未满就 flush；按需启动
	flushTimer := time.NewTimer(defaultTrackBatchFlushInterval)
	if !flushTimer.Stop() {
		<-flushTimer.C
	}
	var flushTimerCh <-chan time.Time
	startFlushTimer := func() {
		if flushTimerCh != nil {
			return
		}
		flushTimer.Reset(defaultTrackBatchFlushInterval)
		flushTimerCh = flushTimer.C
	}
	stopFlushTimer := func() {
		if flushTimerCh == nil {
			return
		}
		if !flushTimer.Stop() {
			select {
			case <-flushTimer.C:
			default:
			}
		}
		flushTimerCh = nil
	}
	flushBatch := func() error {
		stopFlushTimer()
		return w.flushBatch()
	}

	// keepaliveTicker：当配置 keepaliveEvery>0 时持续 tick；
	// 每条真实写入会刷新 lastWriteAt，tick 触发时若距离 lastWriteAt 超过 keepaliveEvery 则写一条心跳。
	var keepaliveCh <-chan time.Time
	var keepaliveTicker *time.Ticker
	lastWriteAt := time.Now()
	if w.keepaliveEvery > 0 {
		keepaliveTicker = time.NewTicker(w.keepaliveEvery)
		keepaliveCh = keepaliveTicker.C
	}

	defer func() {
		if keepaliveTicker != nil {
			keepaliveTicker.Stop()
		}
		w.recordErr(flushBatch())
		if w.file != nil {
			w.recordErr(w.file.Sync())
			w.recordErr(w.file.Close())
			w.file = nil
		}
		close(w.done)
	}()

	for {
		select {
		case msg, ok := <-w.queue:
			if !ok {
				// queue 已关闭，退出
				return
			}
			lastWriteAt = time.Now()
			wasEmpty := w.batchCount == 0
			if w.appendBatch(msg.data) {
				w.recordErr(flushBatch())
			} else if wasEmpty {
				startFlushTimer()
			}

		case replyCh, ok := <-w.sync:
			if !ok {
				return
			}
			// 先排空 queue 里已到达的消息，再 fsync
			drained := w.drainQueue()
			var syncErr error
			if err := flushBatch(); err != nil {
				syncErr = err
				w.recordErr(err)
			}
			if w.file != nil {
				if err := w.file.Sync(); err != nil {
					syncErr = err
					w.recordErr(err)
				}
			}
			if drained != nil {
				syncErr = drained
			}
			if syncErr == nil {
				syncErr = w.lastErr
			}
			w.lastErr = nil
			replyCh <- syncErr

		case <-flushTimerCh:
			flushTimerCh = nil
			w.recordErr(w.flushBatch())

		case now := <-keepaliveCh:
			// 距离上次实际写入超过 keepaliveEvery 才写心跳，避免与正常写入叠加
			if now.Sub(lastWriteAt) < w.keepaliveEvery {
				continue
			}
			lastWriteAt = now
			wasEmpty := w.batchCount == 0
			if w.appendBatch(w.keepaliveMsg) {
				w.recordErr(flushBatch())
			} else if wasEmpty {
				startFlushTimer()
			}
		}
	}
}

func (w *trackWorker) recordErr(err error) {
	if err == nil {
		return
	}
	w.lastErr = err
	_ = monitor.Count(MetricTrackWriteFailed, 1, map[string]string{"channel": w.channel})
}

// drainQueue 在 sync 信号到达时将 queue 中已积压的消息全部写出，返回最后一个写入错误
func (w *trackWorker) drainQueue() error {
	var lastErr error
	for {
		select {
		case msg, ok := <-w.queue:
			if !ok {
				return lastErr
			}
			if w.appendBatch(msg.data) {
				if err := w.flushBatch(); err != nil {
					w.recordErr(err)
					lastErr = err
				}
			}
		default:
			return lastErr
		}
	}
}

func (w *trackWorker) appendBatch(data []byte) bool {
	w.batch = append(w.batch, data...)
	w.batch = append(w.batch, '\n')
	w.batchCount++
	return len(w.batch) >= defaultTrackBatchMaxBytes || w.batchCount >= defaultTrackBatchMaxCount
}

func (w *trackWorker) flushBatch() error {
	if w.batchCount == 0 {
		return nil
	}
	err := w.writeBatch(w.batch)
	w.batch = w.batch[:0]
	w.batchCount = 0
	return err
}

func (w *trackWorker) writeBatch(data []byte) error {
	now := time.Now()
	slot := w.rotation.timeSlot(now)
	if w.file == nil || w.slot != slot {
		if w.file != nil {
			w.recordErr(w.file.Sync())
			w.recordErr(w.file.Close())
		}
		f, err := openTrackFile(w.baseDir, w.channel, w.hostName, slot, dateFromSlot(now))
		if err != nil {
			w.file = nil
			return err
		}
		w.file = f
		w.slot = slot
	}
	_, err := w.file.Write(data)
	return err
}

// dateFromSlot 从当前时间生成日期子目录段（yyyymmdd），与 slot 解耦保证小时/分钟切割模式下日期一致
func dateFromSlot(t time.Time) string { return t.Format("20060102") }

// send 非阻塞投递，队列满时丢弃并记 metric
func (w *trackWorker) send(msg trackMsg) {
	select {
	case w.queue <- msg:
	default:
		_ = monitor.Count(MetricTrackDropped, 1, map[string]string{"channel": w.channel})
	}
}

// doSync 向 worker 发出 sync 请求并等待完成
func (w *trackWorker) doSync() error {
	replyCh := make(chan error, 1)
	select {
	case w.sync <- replyCh:
		select {
		case err := <-replyCh:
			return err
		case <-w.done:
			return w.lastErr
		}
	case <-w.done:
		return w.lastErr
	}
}

// stop 关闭 queue，等待 worker 退出
func (w *trackWorker) stop() error {
	close(w.queue)
	<-w.done
	return w.lastErr
}

// openTrackFile 打开（或追加）一个时间切片文件
//
// 路径布局（ops 规范）：
//
//	{baseDir}/{channel}/{date}/{channel}_{podName}_{slot}.log
//
// 其中 baseDir 已是“该 channel 对应的基础目录”（外层 TrackChannelDirs 覆盖完毕），
// channel 既作为子目录段也作为文件名前缀；podName 来自配置或 HOSTNAME，
// 缺省回退 "unknown" 由 sanitizeTrackFileComponent 兜底。
func openTrackFile(baseDir, channel, podName, slot, date string) (*os.File, error) {
	chSeg := sanitizeTrackFileComponent(channel)
	podSeg := sanitizeTrackFileComponent(podName)
	dir := filepath.Join(baseDir, chSeg, date)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("trackFileWriter: mkdir %s: %w", dir, err)
	}
	name := filepath.Join(dir, fmt.Sprintf("%s_%s_%s.log", chSeg, podSeg, slot))
	f, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("trackFileWriter: open %s: %w", name, err)
	}
	return f, nil
}

func sanitizeTrackFileComponent(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '_' || r == '-' || r == '.':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	out := strings.TrimLeft(b.String(), ".") // 防止 "../" 输入留下以 . 开头的路径段
	if out == "" {
		return "unknown"
	}
	return out
}

// TrackFileWriteSyncerConfig 控制文件落盘行为，路径布局固定为：
//
//	{BaseDir}/{channel}/{yyyymmdd}/{channel}_{hostname}_{slot}.log
//
// 其中 channel 是文件落盘段名，受 ChannelAlias 影响：例如代码侧 dd_meta_channel="thinkingdata" 但 ops
// 期望目录与文件名前缀都用 "tga"，可设 ChannelAlias["thinkingdata"]="tga"。
// ChannelAlias 仅影响文件路径和文件名前缀，不改变原始 dd_meta_channel 或其他输出目标。
//
//   - BaseDir：所有 channel 的根目录（已由调用方在 PMT 环境下从 logbus_track_file_dir 读取）
//   - ChannelAlias：原始 channel 名 → 文件落盘段名映射，未命中保留原 channel
//   - Rotation：时间切割粒度
//   - BufSize：每 channel 队列缓冲大小，默认 defaultTrackChannelBufSize
//   - KeepaliveEvery：心跳间隔；<=0 表示禁用
//   - KeepaliveMessage：心跳触发时落盘的 msg 字符串；可为空
type TrackFileWriteSyncerConfig struct {
	BaseDir          string
	ChannelAlias     map[string]string
	Rotation         TrackRotation
	BufSize          int
	KeepaliveEvery   time.Duration
	KeepaliveMessage string
}

// TrackFileWriteSyncer 是一个 zapcore.WriteSyncer，它：
//  1. 从 JSON 字节中提取 dd_meta_channel（channel）和 msg 字段
//  2. 每个 channel 独享一个 goroutine（trackWorker），通过带缓冲的 Go channel 接收消息
//  3. worker 串行完成文件切割与写入，彻底消除跨 channel 的锁竞争
//  4. 写入调用方只做 JSON 解析 + 投递，不阻塞
type TrackFileWriteSyncer struct {
	mu               sync.RWMutex
	baseDir          string
	channelAlias     map[string]string
	hostName         string
	rotation         TrackRotation
	bufSize          int
	keepaliveEvery   time.Duration
	keepaliveMessage []byte
	workers          map[string]*trackWorker // key = 文件落盘段名（alias 后），按需创建

	closed bool
}

// NewTrackFileWriteSyncer 创建 TrackFileWriteSyncer，使用默认缓冲大小
func NewTrackFileWriteSyncer(baseDir string, rotation TrackRotation) *TrackFileWriteSyncer {
	return NewTrackFileWriteSyncerFromConfig(TrackFileWriteSyncerConfig{
		BaseDir:  baseDir,
		Rotation: rotation,
	})
}

// NewTrackFileWriteSyncerWithBufSize 创建 TrackFileWriteSyncer，可自定义每 channel 缓冲大小
func NewTrackFileWriteSyncerWithBufSize(baseDir string, rotation TrackRotation, bufSize int) *TrackFileWriteSyncer {
	return NewTrackFileWriteSyncerFromConfig(TrackFileWriteSyncerConfig{
		BaseDir:  baseDir,
		Rotation: rotation,
		BufSize:  bufSize,
	})
}

// NewTrackFileWriteSyncerFromConfig 从完整配置创建 TrackFileWriteSyncer。
// hostname 通过 os.Hostname 自动获取，失败时回退 unknown（由 sanitizeTrackFileComponent 兜底）。
func NewTrackFileWriteSyncerFromConfig(cfg TrackFileWriteSyncerConfig) *TrackFileWriteSyncer {
	bufSize := cfg.BufSize
	if bufSize <= 0 {
		bufSize = defaultTrackChannelBufSize
	}
	alias := make(map[string]string, len(cfg.ChannelAlias))
	for k, v := range cfg.ChannelAlias {
		if v == "" {
			continue
		}
		alias[k] = v
	}
	host, _ := os.Hostname()
	if host == "" {
		host = os.Getenv("HOSTNAME")
	}
	return &TrackFileWriteSyncer{
		baseDir:          cfg.BaseDir,
		channelAlias:     alias,
		hostName:         host,
		rotation:         cfg.Rotation,
		bufSize:          bufSize,
		keepaliveEvery:   cfg.KeepaliveEvery,
		keepaliveMessage: []byte(cfg.KeepaliveMessage),
		workers:          make(map[string]*trackWorker),
	}
}

// applyAlias 将原始 channel 名（含 bigquery_xxx 表名后缀）映射为最终文件落盘段：
//  1. 精确匹配 ChannelAlias[channel]
//  2. 对 bigquery_xxx 形态，按前缀 bigquery 命中 alias 后再拼回表名段
//  3. 都没命中保留原值
//
// 该映射只用于文件目录段、文件名前缀和本 writer 的 worker key。
func (w *TrackFileWriteSyncer) applyAlias(channel string) string {
	if v, ok := w.channelAlias[channel]; ok {
		return v
	}
	if idx := strings.IndexByte(channel, '_'); idx > 0 {
		prefix := channel[:idx]
		if v, ok := w.channelAlias[prefix]; ok {
			return v + channel[idx:]
		}
	}
	return channel
}

// Write 实现 io.Writer。
// 解析 JSON 提取 channel + msg，然后非阻塞投递到对应 worker 的队列。
// 队列满时丢弃该条日志（记 metric），不阻塞调用方。
func (w *TrackFileWriteSyncer) Write(p []byte) (n int, err error) {
	channel, msg, ok := extractChannelAndMsg(p)
	if !ok {
		return len(p), nil
	}

	// 计算文件落盘段名后再投递（worker key、文件路径都用映射后的名字）
	finalChannel := w.applyAlias(channel)
	// msg 可能引用 p 的内存（jsoniter.RawMessage 是 slice），需要复制
	data := make([]byte, len(msg))
	copy(data, msg)
	w.send(finalChannel, trackMsg{data: data})
	return len(p), nil
}

// Sync 向所有 worker 发出 sync 请求，等待它们排空队列并 fsync
func (w *TrackFileWriteSyncer) Sync() error {
	w.mu.RLock()
	workers := make([]*trackWorker, 0, len(w.workers))
	for _, wk := range w.workers {
		workers = append(workers, wk)
	}
	w.mu.RUnlock()

	var lastErr error
	for _, wk := range workers {
		if err := wk.doSync(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// SyncChannel 仅对指定文件落盘段名的 worker 执行 sync，若该 worker 尚未创建则 no-op。
// channel 参数应使用文件路径中的 channel 段；配置了 ChannelAlias 时，这里传 alias 后的落盘段名。
// 主要供测试和按需 flush 场景使用。
func (w *TrackFileWriteSyncer) SyncChannel(channel string) error {
	w.mu.RLock()
	wk, ok := w.workers[channel]
	w.mu.RUnlock()
	if !ok {
		return nil
	}
	return wk.doSync()
}

// Close 停止所有 worker，等待它们退出
func (w *TrackFileWriteSyncer) Close() error {
	w.mu.Lock()
	if w.closed {
		w.mu.Unlock()
		return nil
	}
	w.closed = true
	workers := make([]*trackWorker, 0, len(w.workers))
	for _, wk := range w.workers {
		workers = append(workers, wk)
	}
	w.workers = make(map[string]*trackWorker)
	w.mu.Unlock()

	var lastErr error
	for _, wk := range workers {
		if err := wk.stop(); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// send 按 channel 名获取或创建 worker，并在同一把锁内投递消息，避免 Close 并发关闭 queue。
func (w *TrackFileWriteSyncer) send(channel string, msg trackMsg) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return
	}
	wk, ok := w.workers[channel]
	if !ok {
		wk = newTrackWorker(
			channel, w.baseDir, w.hostName, w.rotation, w.bufSize,
			w.keepaliveEvery, w.keepaliveMessage,
		)
		w.workers[channel] = wk
	}
	wk.send(msg)
}

// extractChannelAndMsg 从 JSON 字节中提取 dd_meta_channel 和 msg 字段。
// msg 字段值类型可能是字符串或嵌套 JSON；写文件时直接写原始 []byte。
func extractChannelAndMsg(p []byte) (channel string, msg []byte, ok bool) {
	var raw map[string]jsoniter.RawMessage
	if err := jsonLib.Unmarshal(p, &raw); err != nil {
		return "", nil, false
	}
	chRaw, hasChannel := raw[Meta]
	msgRaw, hasMsg := raw[MsgBody]
	if !hasChannel || !hasMsg {
		return "", nil, false
	}
	// channel 是一个 JSON 字符串，需要 unquote
	if err := jsonLib.Unmarshal(chRaw, &channel); err != nil || channel == "" {
		return "", nil, false
	}
	channel = trackFileChannel(raw, channel)
	// msg 可能是 JSON 字符串（其内容是 TGA/BI JSON）或其他类型；
	// 如果是 JSON 字符串，写文件时只取字符串内容（即 JSON 本体，不带外层引号）；
	// 如果不是字符串类型，直接写原始 JSON 值。
	var msgStr string
	if err := jsonLib.Unmarshal(msgRaw, &msgStr); err == nil {
		// msg 是字符串，内容本身就是 JSON（或普通文本），直接写字符串值
		return channel, []byte(msgStr), true
	}
	// msg 是对象/数组/数字等，直接写原始 JSON
	return channel, msgRaw, true
}

func trackFileChannel(raw map[string]jsoniter.RawMessage, channel string) string {
	if channel != BIGQUERY {
		return channel
	}
	tableRaw, ok := raw[bigquery.TableNameKey]
	if !ok {
		return channel
	}
	var tableName string
	if err := jsonLib.Unmarshal(tableRaw, &tableName); err != nil || tableName == "" {
		return channel
	}
	return channel + "_" + tableName
}

// 确保实现 zapcore.WriteSyncer
var _ zapcore.WriteSyncer = (*TrackFileWriteSyncer)(nil)

// trackOnlyLevelEnabler 只允许 TrackLevel 通过（用于文件写入 core）
type trackOnlyLevelEnabler struct{}

func (trackOnlyLevelEnabler) Enabled(level zapcore.Level) bool {
	return runtimeEnableTrackLevel.Load() && level == TrackLevel
}

// excludeTrackLevelEnabler 将 TrackLevel 从 wrapped enabler 中排除（用于 stdout core）
type excludeTrackLevelEnabler struct {
	wrapped zapcore.LevelEnabler
}

func (e *excludeTrackLevelEnabler) Enabled(level zapcore.Level) bool {
	if level == TrackLevel {
		return false
	}
	return e.wrapped.Enabled(level)
}
