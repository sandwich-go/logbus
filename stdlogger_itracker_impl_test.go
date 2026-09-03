package logbus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/sandwich-go/logbus/thinkingdata"
	"go.uber.org/zap/zapcore"
)

func TestThinkingDataOutputSurvivesInnerBufferReuse(t *testing.T) {
	var output bytes.Buffer
	closeDynamicLogLevel()
	t.Setenv("sys_conf_path_env", "")
	t.Setenv("logbus_enable_trace_level", "1")
	Init(NewConf(
		WithWriteSyncer(zapcore.AddSync(&output)),
		WithEnableTraceLevel(true),
	))
	t.Cleanup(Close)

	const events = 128
	for i := 0; i < events; i++ {
		err := Tracker(THINKINGDATA).Track(
			String(thinkingdata.ACCOUNT, "100001"),
			String(thinkingdata.TYPE, thinkingdata.TRACK),
			String(thinkingdata.EVENT, "buffer_lifecycle"),
			String("sequence", fmt.Sprintf("%d", i)),
		)
		if err != nil {
			t.Fatal(err)
		}
	}

	lines := bytes.Split(bytes.TrimSpace(output.Bytes()), []byte{'\n'})
	if len(lines) != events {
		t.Fatalf("got %d log lines, want %d", len(lines), events)
	}
	for i, line := range lines {
		var envelope map[string]json.RawMessage
		if err := json.Unmarshal(line, &envelope); err != nil {
			t.Fatalf("line %d is not valid outer JSON: %v", i, err)
		}
		var message string
		if err := json.Unmarshal(envelope[MsgBody], &message); err != nil {
			t.Fatalf("line %d has invalid inner JSON string: %v", i, err)
		}
		var inner struct {
			Properties map[string]json.RawMessage `json:"properties"`
		}
		if err := json.Unmarshal([]byte(message), &inner); err != nil {
			t.Fatalf("line %d has invalid ThinkingData JSON: %v", i, err)
		}
		var got string
		if err := json.Unmarshal(inner.Properties["sequence"], &got); err != nil {
			t.Fatalf("line %d has invalid sequence property: %v", i, err)
		}
		if want := fmt.Sprintf("%d", i); got != want {
			t.Fatalf("line %d sequence = %q, want %q", i, got, want)
		}
	}
}
