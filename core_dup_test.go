package logbus

import (
	"bytes"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 创建一个简单的核心实现来验证 Write 调用
type writeTrackerCore struct {
	zapcore.Core
	writeCount int
	lastEntry  zapcore.Entry
	lastFields []zap.Field
}

func (w *writeTrackerCore) Write(ent zapcore.Entry, fields []zap.Field) error {
	w.writeCount++
	w.lastEntry = ent
	w.lastFields = fields
	return w.Core.Write(ent, fields)
}

func (w *writeTrackerCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if w.Enabled(ent.Level) {
		return ce.AddCore(ent, w)
	}
	return ce
}

func TestDedupCoreWriteCalled(t *testing.T) {
	// 创建一个带写入跟踪的核心
	var buffer bytes.Buffer
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	baseCore := zapcore.NewCore(encoder, zapcore.AddSync(&buffer), zapcore.DebugLevel)
	tracker := &writeTrackerCore{Core: baseCore}
	core := NewDupCore(tracker)

	logger := zap.New(core)
	defer logger.Sync()

	// 记录一条日志
	logger.Info("test message", zap.String("key", "value"))

	// 验证 Write 被调用
	if tracker.writeCount == 0 {
		t.Fatal("Write method was not called")
	}

	// 验证日志内容
	if tracker.lastEntry.Message != "test message" {
		t.Errorf("expected message 'test message', got '%s'", tracker.lastEntry.Message)
	}

	// 验证字段
	if len(tracker.lastFields) != 1 || tracker.lastFields[0].Key != "key" {
		t.Errorf("unexpected fields: %v", tracker.lastFields)
	}
}

func TestDedupCoreWithRealWrite(t *testing.T) {
	// 创建一个真实的核心来验证输出
	var buffer bytes.Buffer
	encoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	baseCore := zapcore.NewCore(encoder, zapcore.AddSync(&buffer), zapcore.DebugLevel)
	core := NewDupCore(baseCore)

	logger := zap.New(core)
	defer logger.Sync()

	// 记录3条相同的日志
	for i := 0; i < 3; i++ {
		logger.Info("duplicate message")
	}

	// 同步以刷新任何缓冲的日志
	logger.Sync()

	// 检查输出
	output := buffer.String()
	t.Log("Log output:", output)

	// 应该有一条原始日志和一条重复计数日志
	if strings.Count(output, "duplicate message") != 2 {
		t.Errorf("expected 2 occurrences of 'duplicate message', got %d", strings.Count(output, "duplicate message"))
	}

	if !strings.Contains(output, "[Repeated 2 times]") {
		t.Error("missing repeated message count")
	}
}
