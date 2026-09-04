package logbus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/sandwich-go/logbus/thinkingdata"
	"go.uber.org/zap/zapcore"
)

const (
	thinkingDataBaselineExtraAttributes = 32
	thinkingDataBaselineAttributeBytes  = 256
	thinkingDataBaselineLongTextBytes   = 4096
)

// thinkingDataBaselineFields mirrors the current generic Track call site.
// The conventional profile contains six ThinkingData protocol fields and three
// business properties. The other profiles add business properties only.
func thinkingDataBaselineFields(extraAttributes, attributeBytes int) []Field {
	fields := []Field{
		String(thinkingdata.ACCOUNT, "100001"),
		String(thinkingdata.UUID, "51c7b43d-f8ca-447f-ae05-dfc5fcc184c9"),
		String(thinkingdata.TYPE, thinkingdata.TRACK),
		String(thinkingdata.EVENT, "battle_finish"),
		String(thinkingdata.EVENT_ID, "baseline-001"),
		String(thinkingdata.APPID, "gm15.pressure"),
		String("server_id", "gm15-01"),
		String("result", "win"),
		Int("score", 1280),
	}

	if extraAttributes == 0 {
		return fields
	}

	value := strings.Repeat("payload-", attributeBytes/len("payload-")+1)
	value = value[:attributeBytes]
	for i := 0; i < extraAttributes; i++ {
		fields = append(fields, String(fmt.Sprintf("attribute_%02d", i), value))
	}
	return fields
}

func initThinkingDataBaseline(tb testing.TB, writer zapcore.WriteSyncer) {
	tb.Helper()
	// Package init can have started a PMT config poller before the benchmark
	// installs its isolated environment. Stop it so its work is not sampled.
	closeDynamicLogLevel()
	tb.Setenv("sys_conf_path_env", "")
	tb.Setenv("logbus_enable_trace_level", "1")
	Init(NewConf(
		WithWriteSyncer(writer),
		WithEnableTraceLevel(true),
	))
	tb.Cleanup(Close)
}

// thinkingDataBaselinePipe keeps the ordinary, unbuffered WriteSyncer code
// path while draining output so the benchmark measures this process only.
func thinkingDataBaselinePipe(tb testing.TB) zapcore.WriteSyncer {
	tb.Helper()
	reader, writer, err := os.Pipe()
	if err != nil {
		tb.Fatal(err)
	}
	drained := make(chan struct{})
	go func() {
		_, _ = io.Copy(io.Discard, reader)
		close(drained)
	}()
	tb.Cleanup(func() {
		_ = writer.Close()
		<-drained
		_ = reader.Close()
	})
	return zapcore.AddSync(writer)
}

func benchmarkThinkingDataBaselineTrack(b *testing.B, fields []Field, parallel bool) {
	sink := thinkingDataBaselinePipe(b)
	benchmarkThinkingDataBaselineTrackTo(b, fields, parallel, sink)
}

func benchmarkThinkingDataBaselineTrackTo(b *testing.B, fields []Field, parallel bool, sink zapcore.WriteSyncer) {
	initThinkingDataBaseline(b, sink)
	tracker := Tracker(THINKINGDATA)

	b.ReportAllocs()
	b.ResetTimer()
	if parallel {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				if err := tracker.Track(fields...); err != nil {
					b.Fatal(err)
				}
			}
		})
		return
	}
	for i := 0; i < b.N; i++ {
		if err := tracker.Track(fields...); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkThinkingDataBaselineTrackerCreation(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var tracker ITracker
	for i := 0; i < b.N; i++ {
		tracker = Tracker(THINKINGDATA)
	}
	_ = tracker
}

// Conventional means nine fields: six ThinkingData protocol fields and three
// business properties. It is not a payload-size classification.
func BenchmarkThinkingDataBaselineGenericTrackConventional(b *testing.B) {
	benchmarkThinkingDataBaselineTrack(b, thinkingDataBaselineFields(0, 0), false)
}

// AttributeRich adds 32 business properties of 256 bytes each (about 8 KiB of
// extra source data) to the conventional profile.
func BenchmarkThinkingDataBaselineGenericTrackAttributeRich(b *testing.B) {
	benchmarkThinkingDataBaselineTrack(b, thinkingDataBaselineFields(thinkingDataBaselineExtraAttributes, thinkingDataBaselineAttributeBytes), false)
}

// LongText adds one 4 KiB business property to the conventional profile. It
// separates a large single value from the cost of many property keys.
func BenchmarkThinkingDataBaselineGenericTrackLongText(b *testing.B) {
	benchmarkThinkingDataBaselineTrack(b, thinkingDataBaselineFields(1, thinkingDataBaselineLongTextBytes), false)
}

// Discard removes the OS write from the measurement. It is a process-side
// lower bound, not the default output-path baseline.
func BenchmarkThinkingDataBaselineGenericTrackConventionalDiscard(b *testing.B) {
	benchmarkThinkingDataBaselineTrackTo(b, thinkingDataBaselineFields(0, 0), false, zapcore.AddSync(io.Discard))
}

func BenchmarkThinkingDataBaselineGenericTrackAttributeRichDiscard(b *testing.B) {
	benchmarkThinkingDataBaselineTrackTo(b, thinkingDataBaselineFields(thinkingDataBaselineExtraAttributes, thinkingDataBaselineAttributeBytes), false, zapcore.AddSync(io.Discard))
}

func BenchmarkThinkingDataBaselineGenericTrackConventionalParallel(b *testing.B) {
	benchmarkThinkingDataBaselineTrack(b, thinkingDataBaselineFields(0, 0), true)
}

func BenchmarkThinkingDataBaselineGenericTrackAttributeRichParallel(b *testing.B) {
	benchmarkThinkingDataBaselineTrack(b, thinkingDataBaselineFields(thinkingDataBaselineExtraAttributes, thinkingDataBaselineAttributeBytes), true)
}

func TestThinkingDataBaselineOutput(t *testing.T) {
	cases := []struct {
		name   string
		fields []Field
	}{
		{name: "Conventional", fields: thinkingDataBaselineFields(0, 0)},
		{name: "AttributeRich", fields: thinkingDataBaselineFields(thinkingDataBaselineExtraAttributes, thinkingDataBaselineAttributeBytes)},
		{name: "LongText", fields: thinkingDataBaselineFields(1, thinkingDataBaselineLongTextBytes)},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var output bytes.Buffer
			initThinkingDataBaseline(t, zapcore.AddSync(&output))
			if err := Tracker(THINKINGDATA).Track(tc.fields...); err != nil {
				t.Fatal(err)
			}

			line := output.Bytes()
			if !json.Valid(line) {
				t.Fatalf("output is not valid JSON: %q", line)
			}
			var envelope map[string]json.RawMessage
			if err := json.Unmarshal(line, &envelope); err != nil {
				t.Fatal(err)
			}
			var channel string
			if err := json.Unmarshal(envelope[Meta], &channel); err != nil {
				t.Fatal(err)
			}
			if channel != THINKINGDATA {
				t.Fatalf("unexpected channel %q", channel)
			}
			var message string
			if err := json.Unmarshal(envelope[MsgBody], &message); err != nil {
				t.Fatal(err)
			}
			if !json.Valid([]byte(message)) {
				t.Fatalf("inner ThinkingData message is not valid JSON: %q", message)
			}
			t.Logf("outer_bytes=%d inner_thinkingdata_bytes=%d", len(line), len(message))
		})
	}
}
