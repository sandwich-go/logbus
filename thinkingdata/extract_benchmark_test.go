package thinkingdata

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 固定事件时间；准备字段不计入耗时。覆盖属性数量、复杂字段回退与非法字段过滤。
func BenchmarkExtractFields(b *testing.B) {
	for _, tc := range []struct {
		name  string
		count int
		extra []zap.Field
	}{
		{name: "Scalar3", count: 3},
		{name: "Scalar32", count: 32},
		{name: "Scalar128", count: 128},
		{name: "Object32", count: 32, extra: []zap.Field{zap.Object("object", benchObject{})}},
		{name: "Invalid32", count: 32, extra: []zap.Field{zap.String("invalid-key", "ignored")}},
		{name: "UserAdd32", count: 32, extra: []zap.Field{zap.String(TYPE, USER_ADD)}},
	} {
		b.Run(tc.name, func(b *testing.B) {
			fields := []zap.Field{
				zap.String(ACCOUNT, "100001"),
				zap.String(TIME, "2026-09-20 00:00:00.000"),
				zap.String(TYPE, TRACK),
			}
			if tc.name != "UserAdd32" {
				fields = append(fields, zap.String(EVENT, "battle_finish"))
			}
			for i := 0; i < tc.count; i++ {
				fields = append(fields, zap.Int(fmt.Sprintf("attribute_%02d", i), i))
			}
			fields = append(fields, tc.extra...)
			if _, err := ExtractFields(fields); err != nil {
				b.Fatal(err)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if _, err := ExtractFields(fields); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

type benchObject struct{}

func (benchObject) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("name", "object")
	return nil
}

// 公开入口仍需校验属性名；作为对照组，map 的准备不计入耗时。
func BenchmarkDirectTrack32(b *testing.B) {
	properties := make(map[string]interface{}, 33)
	for i := 0; i < 32; i++ {
		properties[fmt.Sprintf("attribute_%02d", i)] = i
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		properties[TIME] = "2026-09-20 00:00:00.000"
		if _, err := Track("100001", "", "battle_finish", "", "", properties); err != nil {
			b.Fatal(err)
		}
	}
}
