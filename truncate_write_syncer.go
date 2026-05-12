package logbus

import (
	"bytes"
	"io"
	"strings"
	"unicode/utf8"

	jsoniter "github.com/json-iterator/go"
	"github.com/sandwich-go/boost/xconv"
	"github.com/sandwich-go/logbus/monitor"
	"go.uber.org/zap/zapcore"
)

const (
	// TruncateFlag 超出长度限制时输出日志的标志字段，便于检索和过滤
	TruncateFlag = "logbus_truncated"

	DefaultTruncateMaxSize      = 800 * 1024 // 800KB
	DefaultTruncateMsgPrefixLen = 4 * 1024   // 4KB
	DefaultTruncateMsgSuffixLen = 4 * 1024   // 4KB

	// MetricLogTruncated 截断发生时递增的 counter 指标名
	MetricLogTruncated = "logbus_log_truncated"
	// MetricLogStripped 字段剔除（瘦身）发生时递增的 counter 指标名
	MetricLogStripped = "logbus_log_stripped"

	// StripFlag 字段剔除后写出的日志中标记已剔除字段路径的字段名
	StripFlag = "logbus_stripped_fields"

	// OriginalSizeKey 剔除/截断发生后写出日志中记录原始字节数的字段名，
	// 统一前缀避免与业务字段命名冲突。
	OriginalSizeKey = "logbus_original_size"
)

// defaultStripFields 默认在截断前优先剔除的字段路径。
// 私有变量 + 副本导出，避免外部通过 append 共享底层数组污染全局默认值。
var defaultStripFields = []string{
	"api_call.request_body",
	"api_call.response_body",
}

// DefaultStripFields 返回默认剔除字段路径的副本，可通过 WithStripFields 传入自定义列表覆盖。
// 返回副本以防止调用方 append 污染全局默认值。
func DefaultStripFields() []string {
	out := make([]string, len(defaultStripFields))
	copy(out, defaultStripFields)
	return out
}

// json 使用 jsoniter 替代 encoding/json， marshal/unmarshal 性能更好
var jsonLib = jsoniter.ConfigCompatibleWithStandardLibrary

// TruncateWriteSyncer 包装 WriteSyncer，当日志超出限制时输出带标志的截断日志
type TruncateWriteSyncer struct {
	ws zapcore.WriteSyncer
	cc *TruncateWriteSyncerOption
}

// NewTruncateWriteSyncer 创建 TruncateWriteSyncer
func NewTruncateWriteSyncer(ws zapcore.WriteSyncer, cc *TruncateWriteSyncerOption) *TruncateWriteSyncer {
	if org, ok := ws.(*TruncateWriteSyncer); ok {
		// 避免重复包装
		org.cc = cc
		return org
	}
	return &TruncateWriteSyncer{
		ws: ws,
		cc: cc,
	}
}

// Write 实现 io.Writer，超出限制时先尝试剔除大字段；若仍超限则输出截断摘要 JSON
func (t *TruncateWriteSyncer) Write(p []byte) (n int, err error) {
	if len(p) <= t.cc.TruncateMaxSize {
		return t.ws.Write(p)
	}
	// 尝试剔除配置的大字段，若剔除后满足大小限制则直接写出
	if len(t.cc.StripFields) > 0 {
		if stripped, ok := t.stripLargeFields(p); ok {
			_ = monitor.Count(MetricLogStripped, 1, map[string]string{})
			if _, err = t.ws.Write(stripped); err != nil {
				return 0, err
			}
			return len(p), nil
		}
	}
	_ = monitor.Count(MetricLogTruncated, 1, map[string]string{})
	truncated := t.buildTruncatedLog(p)
	_, err = t.ws.Write(truncated)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// Sync 实现 zapcore.WriteSyncer
func (t *TruncateWriteSyncer) Sync() error {
	return t.ws.Sync()
}

func (t *TruncateWriteSyncer) buildTruncatedLog(p []byte) []byte {
	msg, extra := t.extractPartialMsgAndExtra(p)
	out := map[string]interface{}{
		LevelKey:        "warn",
		TruncateFlag:    true,
		"partial_msg":   msg,
		OriginalSizeKey: len(p),
	}
	for k, v := range extra {
		out[k] = v
	}
	bytes, _ := jsonLib.Marshal(out)
	return append(bytes, '\n')
}

// stripLargeFields 尝试从 JSON 日志中剔除 cc.StripFields 指定的字段路径，
// 若剔除后字节数 <= TruncateMaxSize，返回新的日志字节（含 StripFlag 标记）和 true；
// 否则返回 nil, false。
// 使用 jsonLib（jsoniter）Unmarshal/Marshal，配合 UseNumber 保留大整数精度、避免与其他路径混用标准库。
func (t *TruncateWriteSyncer) stripLargeFields(p []byte) ([]byte, bool) {
	var raw map[string]interface{}
	decoder := jsonLib.NewDecoder(bytes.NewReader(p))
	decoder.UseNumber()
	if err := decoder.Decode(&raw); err != nil {
		return nil, false
	}
	var stripped []string
	for _, fieldPath := range t.cc.StripFields {
		if deleteNestedField(raw, strings.Split(fieldPath, ".")) {
			stripped = append(stripped, fieldPath)
		}
	}
	if len(stripped) == 0 {
		return nil, false
	}
	raw[StripFlag] = stripped
	raw[OriginalSizeKey] = len(p)
	b, err := jsonLib.Marshal(raw)
	if err != nil {
		return nil, false
	}
	b = append(b, '\n')
	if len(b) <= t.cc.TruncateMaxSize {
		return b, true
	}
	return nil, false
}

// extractPartialMsgAndExtra 从原日志提取 msg 及 level、date、tags、caller 等定位字段
func (t *TruncateWriteSyncer) extractPartialMsgAndExtra(p []byte) (msg string, extra map[string]interface{}) {
	extra = make(map[string]interface{})
	var raw map[string]interface{}
	if err := jsonLib.Unmarshal(p, &raw); err != nil {
		return t.partialFromBytes(p), extra
	}
	// 提取定位字段：level、date、tags、caller、log_xid、dd_meta_channel 等
	for _, key := range []string{TimeKey, "caller", MsgBody, "Tags", "message", Tags, Meta} {
		if v, ok := raw[key]; ok && v != nil {
			extra[key] = v
		}
	}
	msg = t.partialFromBytes(p)
	return msg, extra
}

func (t *TruncateWriteSyncer) partialFromString(s string) string {
	if len(s) <= t.cc.MsgPrefixLen+t.cc.MsgSuffixLen {
		return s
	}
	return s[:t.cc.MsgPrefixLen] + " ...[truncated]... " + s[len(s)-t.cc.MsgSuffixLen:]
}

// partialFromBytes 从字节切片取 UTF-8 安全的前后缀
func (t *TruncateWriteSyncer) partialFromBytes(p []byte) string {
	prefix := safePrefix(p, t.cc.MsgPrefixLen)
	suffix := safeSuffix(p, t.cc.MsgSuffixLen)
	return string(prefix) + " ...[truncated, original_size=" + xconv.String(len(p)) + "]... " + string(suffix)
}

// safePrefix 取前 n 字节，保证在 UTF-8 码点边界截断
func safePrefix(p []byte, n int) []byte {
	if len(p) <= n {
		return p
	}
	p = p[:n]
	for len(p) > 0 && !utf8.RuneStart(p[len(p)-1]) {
		p = p[:len(p)-1]
	}
	return p
}

// safeSuffix 取后 n 字节，保证在 UTF-8 码点边界截断
func safeSuffix(p []byte, n int) []byte {
	if len(p) <= n {
		return p
	}
	start := len(p) - n
	for start < len(p) && !utf8.RuneStart(p[start]) {
		start++
	}
	return p[start:]
}

// deleteNestedField 按路径删除 JSON map 中的嵌套字段，返回是否找到并删除了该字段。
// 例如 path=["api_call","request_body"] 会删除 m["api_call"]["request_body"]。
// 若路径上某节点不存在或类型不匹配则返回 false。
func deleteNestedField(m map[string]interface{}, path []string) bool {
	if len(path) == 0 {
		return false
	}
	if len(path) == 1 {
		if _, ok := m[path[0]]; ok {
			delete(m, path[0])
			return true
		}
		return false
	}
	child, ok := m[path[0]]
	if !ok {
		return false
	}
	childMap, ok := child.(map[string]interface{})
	if !ok {
		return false
	}
	return deleteNestedField(childMap, path[1:])
}

// 确保实现了 io.Writer 和 zapcore.WriteSyncer
var (
	_ io.Writer           = (*TruncateWriteSyncer)(nil)
	_ zapcore.WriteSyncer = (*TruncateWriteSyncer)(nil)
)
