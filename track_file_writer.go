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
	channel  string
	baseDir  string
	rotation TrackRotation

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

func newTrackWorker(channel, baseDir string, rotation TrackRotation, bufSize int) *trackWorker {
	w := &trackWorker{
		channel:  channel,
		baseDir:  baseDir,
		rotation: rotation,
		queue:    make(chan trackMsg, bufSize),
		sync:     make(chan chan error, 1),
		done:     make(chan struct{}),
	}
	go w.run()
	return w
}

func (w *trackWorker) run() {
	timer := time.NewTimer(defaultTrackBatchFlushInterval)
	if !timer.Stop() {
		<-timer.C
	}
	var timerCh <-chan time.Time
	startTimer := func() {
		if timerCh != nil {
			return
		}
		timer.Reset(defaultTrackBatchFlushInterval)
		timerCh = timer.C
	}
	stopTimer := func() {
		if timerCh == nil {
			return
		}
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timerCh = nil
	}
	flushBatch := func() error {
		stopTimer()
		return w.flushBatch()
	}

	defer func() {
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
			wasEmpty := w.batchCount == 0
			if w.appendBatch(msg.data) {
				w.recordErr(flushBatch())
			} else if wasEmpty {
				startTimer()
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

		case <-timerCh:
			timerCh = nil
			w.recordErr(w.flushBatch())
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
		f, err := openTrackFile(w.baseDir, w.channel, slot)
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
func openTrackFile(baseDir, channel, slot string) (*os.File, error) {
	if err := os.MkdirAll(baseDir, 0o755); err != nil {
		return nil, fmt.Errorf("trackFileWriter: mkdir %s: %w", baseDir, err)
	}
	name := filepath.Join(baseDir, fmt.Sprintf("%s_%s.log", sanitizeTrackFileComponent(channel), slot))
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
	if b.Len() == 0 {
		return "unknown"
	}
	return b.String()
}

// TrackFileWriteSyncer 是一个 zapcore.WriteSyncer，它：
//  1. 从 JSON 字节中提取 dd_meta_channel（channel）和 msg 字段
//  2. 每个 channel 独享一个 goroutine（trackWorker），通过带缓冲的 Go channel 接收消息
//  3. worker 串行完成文件切割与写入，彻底消除跨 channel 的锁竞争
//  4. 写入调用方只做 JSON 解析 + 投递，不阻塞
type TrackFileWriteSyncer struct {
	mu       sync.RWMutex
	baseDir  string
	rotation TrackRotation
	bufSize  int
	workers  map[string]*trackWorker // key = channel，按需创建

	closed bool
}

// NewTrackFileWriteSyncer 创建 TrackFileWriteSyncer，使用默认缓冲大小
func NewTrackFileWriteSyncer(baseDir string, rotation TrackRotation) *TrackFileWriteSyncer {
	return NewTrackFileWriteSyncerWithBufSize(baseDir, rotation, defaultTrackChannelBufSize)
}

// NewTrackFileWriteSyncerWithBufSize 创建 TrackFileWriteSyncer，可自定义每 channel 缓冲大小
func NewTrackFileWriteSyncerWithBufSize(baseDir string, rotation TrackRotation, bufSize int) *TrackFileWriteSyncer {
	return &TrackFileWriteSyncer{
		baseDir:  baseDir,
		rotation: rotation,
		bufSize:  bufSize,
		workers:  make(map[string]*trackWorker),
	}
}

// Write 实现 io.Writer。
// 解析 JSON 提取 channel + msg，然后非阻塞投递到对应 worker 的队列。
// 队列满时丢弃该条日志（记 metric），不阻塞调用方。
func (w *TrackFileWriteSyncer) Write(p []byte) (n int, err error) {
	channel, msg, ok := extractChannelAndMsg(p)
	if !ok {
		return len(p), nil
	}

	// msg 可能引用 p 的内存（jsoniter.RawMessage 是 slice），需要复制
	data := make([]byte, len(msg))
	copy(data, msg)
	w.send(channel, trackMsg{data: data})
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

// SyncChannel 仅对指定 channel 的 worker 执行 sync，若该 worker 尚未创建则 no-op。
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
		wk = newTrackWorker(channel, w.baseDir, w.rotation, w.bufSize)
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
