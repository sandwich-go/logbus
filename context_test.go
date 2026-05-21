package logbus

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

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
