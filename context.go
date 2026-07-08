package logbus

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// ContextMetaFieldKey is the top-level log object key for context metadata.
	ContextMetaFieldKey = "meta"
	// ContextTraceIDKey is the trace id key inside the context metadata log object.
	ContextTraceIDKey = "trace_id"
)

type traceIDContextKey struct{}

type contextMeta struct {
	traceID string
}

// ContextWithTraceID stores traceID in context.Value.
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if traceID == "" {
		return ctx
	}
	return context.WithValue(ctx, traceIDContextKey{}, traceID)
}

// TraceIDFromContext returns the trace id stored by ContextWithTraceID.
func TraceIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if traceID, ok := ctx.Value(traceIDContextKey{}).(string); ok {
		return traceID
	}
	return ""
}

// FromContext returns request-scoped log metadata extracted from ctx.
// It returns zap.Skip when ctx has no trace id.
func FromContext(ctx context.Context) Field {
	m := contextMeta{
		traceID: TraceIDFromContext(ctx),
	}
	if m.empty() {
		return zap.Skip()
	}
	return Object(ContextMetaFieldKey, m)
}

func (m contextMeta) empty() bool {
	return m.traceID == ""
}

func (m contextMeta) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if m.traceID != "" {
		enc.AddString(ContextTraceIDKey, m.traceID)
	}
	return nil
}
