package logbus

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"go.uber.org/zap/zapcore"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTruncateWriteSyncer_WithinLimit(t *testing.T) {
	Convey("未超出限制时透传", t, func() {
		buf := &bytes.Buffer{}
		cc := NewTruncateWriteSyncerOption(WithTruncateMaxSize(1024))
		ws := NewTruncateWriteSyncer(zapcore.AddSync(buf), cc)
		msg := `{"level":"info","msg":"hello"}`
		n, err := ws.Write([]byte(msg))
		So(err, ShouldBeNil)
		So(n, ShouldEqual, len(msg))
		So(buf.String(), ShouldEqual, msg)
	})
}

func TestTruncateWriteSyncer_ExceedLimit(t *testing.T) {
	Convey("超出限制时输出截断 JSON", t, func() {
		buf := &bytes.Buffer{}
		cc := NewTruncateWriteSyncerOption(
			WithTruncateMaxSize(50),
			WithMsgPrefixLen(20),
			WithMsgSuffixLen(10),
		)
		ws := NewTruncateWriteSyncer(zapcore.AddSync(buf), cc)
		original := `{"level":"info","msg":"this is a very long message that exceeds the limit","ts":123}`
		n, err := ws.Write([]byte(original))
		So(err, ShouldBeNil)
		So(n, ShouldEqual, len(original))

		out := buf.String()
		So(out, ShouldContainSubstring, TruncateFlag)
		So(out, ShouldContainSubstring, "true")

		var parsed map[string]interface{}
		So(json.Unmarshal([]byte(strings.TrimSpace(out)), &parsed), ShouldBeNil)
		So(parsed["logbus_truncated"], ShouldBeTrue)
		So(parsed["partial_msg"], ShouldNotBeNil)
		So(parsed["original_size"], ShouldEqual, float64(len(original)))
		partialMsg, _ := parsed["partial_msg"].(string)
		So(partialMsg, ShouldContainSubstring, "original_size")
	})
}

func TestTruncateWriteSyncer_ExceedLimit_WithExtraFields(t *testing.T) {
	Convey("超出限制时透传 date、tags、caller 等定位字段", t, func() {
		buf := &bytes.Buffer{}
		cc := NewTruncateWriteSyncerOption(
			WithTruncateMaxSize(50),
			WithMsgPrefixLen(10),
			WithMsgSuffixLen(5),
		)
		ws := NewTruncateWriteSyncer(zapcore.AddSync(buf), cc)
		original := `{"log_level":"info","date":"2025-02-24T10:00:00Z","tags":"server","log_xid":"abc123","caller":"main.go:42","msg":"long message xxx"}`
		_, err := ws.Write([]byte(original))
		So(err, ShouldBeNil)

		var parsed map[string]interface{}
		So(json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &parsed), ShouldBeNil)
		So(parsed["logbus_truncated"], ShouldBeTrue)
		So(parsed["log_level"], ShouldEqual, "error")
		So(parsed["date"], ShouldEqual, "2025-02-24T10:00:00Z")
		So(parsed["tags"], ShouldEqual, "server")
		So(parsed["caller"], ShouldEqual, "main.go:42")
		So(parsed["msg"], ShouldEqual, "long message xxx")
		So(parsed["partial_msg"], ShouldNotBeNil)
	})
}

func TestTruncateWriteSyncer_PartialFromBytes(t *testing.T) {
	Convey("非 JSON 时使用前后缀", t, func() {
		buf := &bytes.Buffer{}
		cc := NewTruncateWriteSyncerOption(
			WithTruncateMaxSize(50),
			WithMsgPrefixLen(10),
			WithMsgSuffixLen(5),
		)
		ws := NewTruncateWriteSyncer(zapcore.AddSync(buf), cc)
		original := `not valid json {{{{{{`
		for len(original) < 100 {
			original += "x"
		}
		_, err := ws.Write([]byte(original))
		So(err, ShouldBeNil)

		var parsed map[string]interface{}
		So(json.Unmarshal([]byte(strings.TrimSpace(buf.String())), &parsed), ShouldBeNil)
		So(parsed["logbus_truncated"], ShouldBeTrue)
		partialMsg, _ := parsed["partial_msg"].(string)
		So(partialMsg, ShouldContainSubstring, "not valid ")
		So(partialMsg, ShouldContainSubstring, "original_size")
	})
}

func TestTruncateWriteSyncer_UTF8Safe(t *testing.T) {
	Convey("UTF-8 边界安全截断", t, func() {
		buf := &bytes.Buffer{}
		cc := NewTruncateWriteSyncerOption(
			WithTruncateMaxSize(20),
			WithMsgPrefixLen(5),
			WithMsgSuffixLen(3),
		)
		ws := NewTruncateWriteSyncer(zapcore.AddSync(buf), cc)
		original := `{"msg":"你好世界Hello"}`
		_, err := ws.Write([]byte(original))
		So(err, ShouldBeNil)

		out := strings.TrimSpace(buf.String())
		So(out, ShouldContainSubstring, TruncateFlag)
		var parsed map[string]interface{}
		So(json.Unmarshal([]byte(out), &parsed), ShouldBeNil)
		So(parsed["logbus_truncated"], ShouldBeTrue)
	})
}
