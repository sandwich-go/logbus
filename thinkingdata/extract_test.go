package thinkingdata

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type testObject struct{}

func (testObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", "object")
	return nil
}

type testArray struct{}

func (testArray) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	enc.AppendString("array")
	return nil
}

type testStringer string

func (s testStringer) String() string {
	return string(s)
}

func TestExtractFieldsMatchesMapObjectEncoder(t *testing.T) {
	fixedTime := time.Date(2026, 9, 3, 12, 0, 0, 0, time.UTC)
	tests := []struct {
		name     string
		fields   []zap.Field
		wantFast bool
	}{
		{
			name: "scalar fields and protocol fields",
			fields: []zap.Field{
				zap.String(ACCOUNT, "100001"),
				zap.String(TYPE, TRACK),
				zap.String(EVENT, "battle_finish"),
				zap.String(EVENT_ID, "event-001"),
				zap.String(APPID, "gm15.pressure"),
				zap.Time(TIME, fixedTime),
				zap.String(IP, "10.0.0.1"),
				zap.String(UUID, "51c7b43d-f8ca-447f-ae05-dfc5fcc184c9"),
				zap.String(FIRST_CHECK_ID, "first-check-001"),
				zap.Bool("success", true),
				zap.Int("score", 1280),
				zap.String("result", "first"),
				zap.String("result", "last"),
			},
			wantFast: true,
		},
		{
			name: "invalid key falls back",
			fields: []zap.Field{
				zap.String(ACCOUNT, "100001"),
				zap.String(TYPE, TRACK),
				zap.String(EVENT, "battle_finish"),
				zap.Time(TIME, fixedTime),
				zap.String("invalid-key", "ignored"),
			},
		},
		{
			name: "reflected field",
			fields: []zap.Field{
				zap.String(ACCOUNT, "100001"),
				zap.String(TYPE, TRACK),
				zap.String(EVENT, "battle_finish"),
				zap.Time(TIME, fixedTime),
				zap.Reflect("metadata", map[string]int{"round": 3}),
			},
			wantFast: true,
		},
		{
			name: "object field falls back",
			fields: []zap.Field{
				zap.String(ACCOUNT, "100001"),
				zap.String(TYPE, TRACK),
				zap.String(EVENT, "battle_finish"),
				zap.Time(TIME, fixedTime),
				zap.Object("metadata", testObject{}),
			},
		},
		{
			name: "array field falls back",
			fields: []zap.Field{
				zap.String(ACCOUNT, "100001"),
				zap.String(TYPE, TRACK),
				zap.String(EVENT, "battle_finish"),
				zap.Time(TIME, fixedTime),
				zap.Array("metadata", testArray{}),
			},
		},
		{
			name: "invalid event id type",
			fields: []zap.Field{
				zap.String(ACCOUNT, "100001"),
				zap.String(TYPE, TRACK),
				zap.String(EVENT, "battle_finish"),
				zap.Int(EVENT_ID, 1),
			},
		},
		{
			name: "user add",
			fields: []zap.Field{
				zap.String(DISTINCT, "device-001"),
				zap.String(TYPE, USER_ADD),
				zap.Time(TIME, fixedTime),
				zap.Int("score", 1),
			},
			wantFast: true,
		},
		{
			name: "namespace falls back",
			fields: []zap.Field{
				zap.String(ACCOUNT, "100001"),
				zap.String(TYPE, TRACK),
				zap.String(EVENT, "battle_finish"),
				zap.Time(TIME, fixedTime),
				zap.Namespace("nested"),
				zap.String("value", "nested"),
			},
		},
		{
			name: "inline falls back",
			fields: []zap.Field{
				zap.String(ACCOUNT, "100001"),
				zap.String(TYPE, TRACK),
				zap.String(EVENT, "battle_finish"),
				zap.Time(TIME, fixedTime),
				zap.Inline(testObject{}),
			},
		},
		{
			name: "error falls back",
			fields: []zap.Field{
				zap.String(ACCOUNT, "100001"),
				zap.String(TYPE, TRACK),
				zap.String(EVENT, "battle_finish"),
				zap.Time(TIME, fixedTime),
				zap.Error(errors.New("failed")),
			},
		},
		{
			name: "stringer falls back",
			fields: []zap.Field{
				zap.String(ACCOUNT, "100001"),
				zap.String(TYPE, TRACK),
				zap.String(EVENT, "battle_finish"),
				zap.Time(TIME, fixedTime),
				zap.Stringer("name", testStringer("stringer")),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, gotFast := extractScalarFields(tt.fields)
			if gotFast != tt.wantFast {
				t.Fatalf("fast path = %t, want %t", gotFast, tt.wantFast)
			}
			got, gotErr := ExtractFields(tt.fields)
			want, wantErr := extractWithMapObjectEncoder(tt.fields)
			if !sameError(gotErr, wantErr) {
				t.Fatalf("error = %v, want %v", gotErr, wantErr)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("data = %#v, want %#v", got, want)
			}
		})
	}
}

func TestScalarFieldValueMatchesMapObjectEncoder(t *testing.T) {
	fixedTime := time.Date(2026, 9, 3, 12, 0, 0, 0, time.FixedZone("UTC+8", 8*60*60))
	fields := []zap.Field{
		zap.Binary("binary", []byte("value")),
		zap.Bool("bool", true),
		zap.ByteString("byte_string", []byte("value")),
		zap.Complex128("complex128", 1+2i),
		zap.Complex64("complex64", 3+4i),
		zap.Duration("duration", time.Second),
		zap.Float64("float64", 1.5),
		zap.Float32("float32", 2.5),
		zap.Int64("int64", 1),
		zap.Int32("int32", 2),
		zap.Int16("int16", 3),
		zap.Int8("int8", 4),
		zap.String("string", "value"),
		zap.Time("time", fixedTime),
		{Key: "time_full", Type: zapcore.TimeFullType, Interface: fixedTime.AddDate(300, 0, 0)},
		zap.Uint64("uint64", 1),
		zap.Uint32("uint32", 2),
		zap.Uint16("uint16", 3),
		zap.Uint8("uint8", 4),
		zap.Uintptr("uintptr", 5),
		zap.Reflect("reflect", map[string]int{"value": 1}),
		zap.Skip(),
	}

	for _, field := range fields {
		t.Run(field.Key, func(t *testing.T) {
			got, skip, ok := fieldValue(field)
			if !ok {
				t.Fatal("scalar field is not supported")
			}

			encoder := zapcore.NewMapObjectEncoder()
			field.AddTo(encoder)
			want, exists := encoder.Fields[field.Key]
			if skip {
				if exists {
					t.Fatalf("skip field wrote %#v", want)
				}
				return
			}
			if !exists {
				t.Fatal("map object encoder did not write a value")
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("value = %#v, want %#v", got, want)
			}
		})
	}
}

func extractWithMapObjectEncoder(fields []zap.Field) (Data, error) {
	encoder := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(encoder)
	}
	return ExtractEncoder(encoder)
}

func sameError(got, want error) bool {
	if got == nil || want == nil {
		return got == want
	}
	return got.Error() == want.Error()
}
