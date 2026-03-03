package logbus

import (
	"io"
	"unicode/utf8"

	"github.com/json-iterator/go"
	"github.com/sandwich-go/boost/xconv"
	"go.uber.org/zap/zapcore"
)

const (
	// TruncateFlag 超出长度限制时输出日志的标志字段，便于检索和过滤
	TruncateFlag = "logbus_truncated"

	DefaultTruncateMaxSize      = 800 * 1024 // 800KB
	DefaultTruncateMsgPrefixLen = 4 * 1024   // 4KB
	DefaultTruncateMsgSuffixLen = 4 * 1024   // 4KB
)

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

// Write 实现 io.Writer，超出限制时输出截断 JSON
func (t *TruncateWriteSyncer) Write(p []byte) (n int, err error) {
	if len(p) <= t.cc.TruncateMaxSize {
		return t.ws.Write(p)
	}
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
		LevelKey:        "error",
		TruncateFlag:    true,
		"partial_msg":   msg,
		"original_size": len(p),
	}
	for k, v := range extra {
		out[k] = v
	}
	bytes, _ := jsonLib.Marshal(out)
	return append(bytes, '\n')
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

// 确保实现了 io.Writer 和 zapcore.WriteSyncer
var (
	_ io.Writer           = (*TruncateWriteSyncer)(nil)
	_ zapcore.WriteSyncer = (*TruncateWriteSyncer)(nil)
)
