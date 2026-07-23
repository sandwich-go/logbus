package logbus

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var testContext = context.Background()

func TestLoggerRequiresContext(t *testing.T) {
	loggerType := reflect.TypeOf((*NewILogger)(nil)).Elem()
	contextType := reflect.TypeOf((*context.Context)(nil)).Elem()

	levels := []string{"Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal"}
	for _, level := range levels {
		if _, ok := loggerType.MethodByName(level + "WithContext"); ok {
			t.Errorf("NewILogger.%sWithContext exists; use %s with context as its first argument", level, level)
		}

		method, ok := loggerType.MethodByName(level)
		if !ok {
			t.Errorf("NewILogger.%s is missing", level)
			continue
		}
		if method.Type.NumIn() == 0 || method.Type.In(0) != contextType {
			t.Errorf("NewILogger.%s first argument = %v, want context.Context", level, method.Type.In(0))
		}
	}
}

func TestDepthLoggerRequiresContext(t *testing.T) {
	loggerType := reflect.TypeOf((*NewILogger)(nil)).Elem()
	contextType := reflect.TypeOf((*context.Context)(nil)).Elem()

	levels := []string{"Debug", "Info", "Warn", "Error", "Fatal"}
	for _, level := range levels {
		name := "G" + level + "Depth"
		if _, ok := loggerType.MethodByName(name + "WithContext"); ok {
			t.Errorf("NewILogger.%sWithContext exists; use %s with context as its first argument", name, name)
		}

		method, ok := loggerType.MethodByName(name)
		if !ok {
			t.Errorf("NewILogger.%s is missing", name)
			continue
		}
		if method.Type.NumIn() == 0 || method.Type.In(0) != contextType {
			t.Errorf("NewILogger.%s first argument = %v, want context.Context", name, method.Type.In(0))
		}
	}
}

func TestChannelLoggerRequiresContext(t *testing.T) {
	loggerType := reflect.TypeOf((*NewILogger)(nil)).Elem()
	contextType := reflect.TypeOf((*context.Context)(nil)).Elem()

	levels := []string{"Debug", "Info", "Warn", "Error", "DPanic", "Panic", "Fatal"}
	for _, level := range levels {
		name := level + "WithChannel"
		method, ok := loggerType.MethodByName(name)
		if !ok {
			t.Errorf("NewILogger.%s is missing", name)
			continue
		}
		if method.Type.NumIn() == 0 || method.Type.In(0) != contextType {
			t.Errorf("NewILogger.%s first argument = %v, want context.Context", name, method.Type.In(0))
		}
	}
}

func TestContextTraceID(t *testing.T) {
	const traceID = "2a078316dd16823775edaa075f78eb25"

	if got := TraceIDFromContext(nil); got != "" {
		t.Fatalf("TraceIDFromContext(nil) = %q, want empty", got)
	}
	if got := TraceIDFromContext(context.Background()); got != "" {
		t.Fatalf("TraceIDFromContext(context.Background()) = %q, want empty", got)
	}

	ctx := ContextWithTraceID(context.Background(), traceID)
	if got := TraceIDFromContext(ctx); got != traceID {
		t.Fatalf("TraceIDFromContext(ctx) = %q, want %q", got, traceID)
	}

	ctx = ContextWithTraceID(ctx, "")
	if got := TraceIDFromContext(ctx); got != traceID {
		t.Fatalf("ContextWithTraceID(ctx, empty) changed trace id to %q, want %q", got, traceID)
	}
}

func TestFromContext(t *testing.T) {
	if got := FromContext(nil); got.Type != zap.Skip().Type {
		t.Fatalf("FromContext(nil).Type = %v, want %v", got.Type, zap.Skip().Type)
	}
	if got := FromContext(context.Background()); got.Type != zap.Skip().Type {
		t.Fatalf("FromContext(context.Background()).Type = %v, want %v", got.Type, zap.Skip().Type)
	}

	const traceID = "2a078316dd16823775edaa075f78eb25"
	field := FromContext(ContextWithTraceID(context.Background(), traceID))
	if field.Key != ContextMetaFieldKey {
		t.Fatalf("field.Key = %q, want %q", field.Key, ContextMetaFieldKey)
	}
	if field.Type != zapcore.ObjectMarshalerType {
		t.Fatalf("field.Type = %v, want %v", field.Type, zapcore.ObjectMarshalerType)
	}

	marshaler, ok := field.Interface.(zapcore.ObjectMarshaler)
	if !ok {
		t.Fatalf("field.Interface = %T, want zapcore.ObjectMarshaler", field.Interface)
	}

	enc := zapcore.NewMapObjectEncoder()
	if err := marshaler.MarshalLogObject(enc); err != nil {
		t.Fatalf("MarshalLogObject() error = %v", err)
	}
	if got := enc.Fields[ContextTraceIDKey]; got != traceID {
		t.Fatalf("encoded trace id = %v, want %q", got, traceID)
	}
}

func TestGLoggerInfo(t *testing.T) {
	const traceID = "2a078316dd16823775edaa075f78eb25"

	core, logs := observer.New(zapcore.DebugLevel)
	logger := &GLogger{
		stdLogger: newStdLogger(zap.New(core), currentConfig(), nil),
	}

	logger.Info(ContextWithTraceID(context.Background(), traceID), "info", String("k", "v"))

	if logs.Len() != 1 {
		t.Fatalf("logs.Len() = %d, want 1", logs.Len())
	}

	fields := logs.All()[0].ContextMap()
	meta, ok := fields[ContextMetaFieldKey].(map[string]interface{})
	if !ok {
		t.Fatalf("fields[%q] = %T, want map[string]interface{}", ContextMetaFieldKey, fields[ContextMetaFieldKey])
	}
	if got := meta[ContextTraceIDKey]; got != traceID {
		t.Fatalf("meta[%q] = %v, want %q", ContextTraceIDKey, got, traceID)
	}
	if got := fields[MsgBody]; got != "info" {
		t.Fatalf("fields[%q] = %v, want %q", MsgBody, got, "info")
	}
	if got := fields["k"]; got != "v" {
		t.Fatalf("fields[%q] = %v, want %q", "k", got, "v")
	}
}

func TestGLoggerInfoWithChannel(t *testing.T) {
	const (
		traceID = "2a078316dd16823775edaa075f78eb25"
		channel = "custom"
	)

	core, logs := observer.New(zapcore.DebugLevel)
	logger := &GLogger{
		stdLogger: newStdLogger(zap.New(core), currentConfig(), nil),
	}

	ctx := ContextWithTraceID(context.Background(), traceID)
	logger.InfoWithChannel(ctx, channel, "info", String("k", "v"))

	if logs.Len() != 1 {
		t.Fatalf("logs.Len() = %d, want 1", logs.Len())
	}
	entry := logs.All()[0]
	if entry.Message != channel {
		t.Fatalf("entry.Message = %q, want channel %q", entry.Message, channel)
	}
	fields := entry.ContextMap()
	meta, ok := fields[ContextMetaFieldKey].(map[string]interface{})
	if !ok {
		t.Fatalf("fields[%q] = %T, want map[string]interface{}", ContextMetaFieldKey, fields[ContextMetaFieldKey])
	}
	if got := meta[ContextTraceIDKey]; got != traceID {
		t.Fatalf("meta[%q] = %v, want %q", ContextTraceIDKey, got, traceID)
	}
	if got := fields[MsgBody]; got != "info" {
		t.Fatalf("fields[%q] = %v, want %q", MsgBody, got, "info")
	}
}

func writeDepthInfo(ctx context.Context) {
	InfoDepth(ctx, 1, "depth info")
}

func TestInfoDepth(t *testing.T) {
	const traceID = "2a078316dd16823775edaa075f78eb25"
	var output bytes.Buffer

	Init(NewConf(
		WithDev(true),
		WithWriteSyncer(zapcore.AddSync(&output)),
		WithDisableTruncateWriteSyncer(true),
		WithUseSystemClock(true),
	))
	defer func() {
		Close()
		Init(NewConf())
	}()

	ctx := ContextWithTraceID(context.Background(), traceID)
	_, _, line, _ := runtime.Caller(0)
	writeDepthInfo(ctx)
	if err := DefaultLogger().Sync(); err != nil {
		t.Fatalf("Sync() error = %v", err)
	}

	got := output.String()
	if !strings.Contains(got, fmt.Sprintf(`"%s": "%s"`, ContextTraceIDKey, traceID)) {
		t.Fatalf("output = %q, want trace id %q", got, traceID)
	}

	wantCaller := fmt.Sprintf("context_test.go:%d", line+1)
	if !strings.Contains(got, wantCaller) {
		t.Fatalf("output = %q, want caller %q", got, wantCaller)
	}
}
