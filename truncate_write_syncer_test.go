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
		So(parsed[OriginalSizeKey], ShouldEqual, float64(len(original)))
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
		So(parsed["log_level"], ShouldEqual, "warn")
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

// buildLargeAPICallLog 构造一条包含 api_call.request_body/response_body 的超大日志
func buildLargeAPICallLog(maxSize int) string {
	// 填充一个超过 maxSize 的大 body
	bigBody := strings.Repeat("x", maxSize+100)
	m := map[string]interface{}{
		"log_level": "info",
		"date":      "2025-02-24T10:00:00Z",
		"tags":      "server",
		"caller":    "main.go:42",
		"msg":       "api call log",
		"api_call": map[string]interface{}{
			"url":           "https://example.com/api",
			"status":        200,
			"request_body":  bigBody,
			"response_body": bigBody,
		},
	}
	b, _ := json.Marshal(m)
	return string(b)
}

func TestTruncateWriteSyncer_StripFields_ApiCall(t *testing.T) {
	Convey("超出限制时剔除 api_call body 后满足限制，直接写出保留其他字段", t, func() {
		buf := &bytes.Buffer{}
		cc := NewTruncateWriteSyncerOption(
			WithTruncateMaxSize(500),
			// StripFields 默认已包含 api_call.request_body / api_call.response_body
		)
		ws := NewTruncateWriteSyncer(zapcore.AddSync(buf), cc)
		original := buildLargeAPICallLog(500)
		So(len(original), ShouldBeGreaterThan, 500)

		n, err := ws.Write([]byte(original))
		So(err, ShouldBeNil)
		So(n, ShouldEqual, len(original))

		out := strings.TrimSpace(buf.String())
		var parsed map[string]interface{}
		So(json.Unmarshal([]byte(out), &parsed), ShouldBeNil)

		// 未走截断流程，无 logbus_truncated
		So(parsed[TruncateFlag], ShouldBeNil)
		// 标记了哪些字段被剔除
		So(parsed[StripFlag], ShouldNotBeNil)
		// 保留了业务字段
		So(parsed["log_level"], ShouldEqual, "info")
		So(parsed["date"], ShouldEqual, "2025-02-24T10:00:00Z")
		So(parsed["tags"], ShouldEqual, "server")
		So(parsed["msg"], ShouldEqual, "api call log")
		// api_call 的 url/status 应保留，body 被剔除
		apiCall, ok := parsed["api_call"].(map[string]interface{})
		So(ok, ShouldBeTrue)
		So(apiCall["url"], ShouldEqual, "https://example.com/api")
		So(apiCall["status"], ShouldEqual, float64(200))
		So(apiCall["request_body"], ShouldBeNil)
		So(apiCall["response_body"], ShouldBeNil)
		// OriginalSizeKey 记录原始大小
		So(parsed[OriginalSizeKey], ShouldEqual, float64(len(original)))
	})
}

func TestTruncateWriteSyncer_StripFields_FallbackToTruncate(t *testing.T) {
	Convey("剔除字段后仍超限，降级走摘要截断流程", t, func() {
		buf := &bytes.Buffer{}
		// 限制极小（50字节），即使剔除 body 也无法满足
		cc := NewTruncateWriteSyncerOption(
			WithTruncateMaxSize(50),
			WithMsgPrefixLen(10),
			WithMsgSuffixLen(5),
		)
		ws := NewTruncateWriteSyncer(zapcore.AddSync(buf), cc)
		original := buildLargeAPICallLog(500)

		n, err := ws.Write([]byte(original))
		So(err, ShouldBeNil)
		So(n, ShouldEqual, len(original))

		out := strings.TrimSpace(buf.String())
		var parsed map[string]interface{}
		So(json.Unmarshal([]byte(out), &parsed), ShouldBeNil)
		// 应走截断摘要流程
		So(parsed[TruncateFlag], ShouldBeTrue)
		So(parsed["partial_msg"], ShouldNotBeNil)
	})
}

func TestTruncateWriteSyncer_StripFields_CustomPath(t *testing.T) {
	Convey("自定义 StripFields 剔除顶层字段", t, func() {
		buf := &bytes.Buffer{}
		bigVal := strings.Repeat("y", 600)
		original := `{"level":"info","msg":"test","raw_data":"` + bigVal + `","keep_field":"should_remain"}`
		cc := NewTruncateWriteSyncerOption(
			WithTruncateMaxSize(200),
			WithStripFields("raw_data"),
		)
		ws := NewTruncateWriteSyncer(zapcore.AddSync(buf), cc)
		So(len(original), ShouldBeGreaterThan, 200)

		_, err := ws.Write([]byte(original))
		So(err, ShouldBeNil)

		out := strings.TrimSpace(buf.String())
		var parsed map[string]interface{}
		So(json.Unmarshal([]byte(out), &parsed), ShouldBeNil)
		So(parsed[TruncateFlag], ShouldBeNil)
		So(parsed[StripFlag], ShouldNotBeNil)
		So(parsed["raw_data"], ShouldBeNil)
		So(parsed["keep_field"], ShouldEqual, "should_remain")
	})
}

func TestTruncateWriteSyncer_StripFields_PreserveLargeNumber(t *testing.T) {
	Convey("剔除字段时保留非目标大整数精度", t, func() {
		buf := &bytes.Buffer{}
		bigVal := strings.Repeat("y", 600)
		largeID := "9223372036854775807"
		original := `{"level":"info","order_id":` + largeID + `,"raw_data":"` + bigVal + `","keep_field":"should_remain"}`
		cc := NewTruncateWriteSyncerOption(
			WithTruncateMaxSize(220),
			WithStripFields("raw_data"),
		)
		ws := NewTruncateWriteSyncer(zapcore.AddSync(buf), cc)

		_, err := ws.Write([]byte(original))
		So(err, ShouldBeNil)

		out := strings.TrimSpace(buf.String())
		So(out, ShouldContainSubstring, `"order_id":`+largeID)
		So(out, ShouldNotContainSubstring, `"order_id":9.223372036854776e+18`)

		var parsed map[string]interface{}
		So(json.Unmarshal([]byte(out), &parsed), ShouldBeNil)
		So(parsed[TruncateFlag], ShouldBeNil)
		So(parsed["raw_data"], ShouldBeNil)
	})
}

func TestTruncateWriteSyncer_StripFields_NoMatchField(t *testing.T) {
	Convey("StripFields 中的字段不存在时，仍走截断流程", t, func() {
		buf := &bytes.Buffer{}
		bigVal := strings.Repeat("z", 600)
		original := `{"level":"info","msg":"` + bigVal + `"}`
		cc := NewTruncateWriteSyncerOption(
			WithTruncateMaxSize(200),
			WithMsgPrefixLen(10),
			WithMsgSuffixLen(5),
			// StripFields 中的字段 original 日志里不存在
			WithStripFields("api_call.request_body", "api_call.response_body"),
		)
		ws := NewTruncateWriteSyncer(zapcore.AddSync(buf), cc)

		_, err := ws.Write([]byte(original))
		So(err, ShouldBeNil)

		out := strings.TrimSpace(buf.String())
		var parsed map[string]interface{}
		So(json.Unmarshal([]byte(out), &parsed), ShouldBeNil)
		// 无可剔除字段，走截断摘要流程
		So(parsed[TruncateFlag], ShouldBeTrue)
	})
}

func TestDeleteNestedField(t *testing.T) {
	Convey("deleteNestedField", t, func() {
		Convey("删除顶层字段", func() {
			m := map[string]interface{}{"a": 1, "b": 2}
			So(deleteNestedField(m, []string{"a"}), ShouldBeTrue)
			So(m["a"], ShouldBeNil)
			So(m["b"], ShouldEqual, 2)
		})
		Convey("删除嵌套字段", func() {
			m := map[string]interface{}{
				"api_call": map[string]interface{}{
					"request_body":  "big",
					"response_body": "big",
					"url":           "http://x",
				},
			}
			So(deleteNestedField(m, []string{"api_call", "request_body"}), ShouldBeTrue)
			inner := m["api_call"].(map[string]interface{})
			So(inner["request_body"], ShouldBeNil)
			So(inner["url"], ShouldEqual, "http://x")
		})
		Convey("字段不存在返回 false", func() {
			m := map[string]interface{}{"x": 1}
			So(deleteNestedField(m, []string{"not_exist"}), ShouldBeFalse)
		})
		Convey("路径中间节点类型不匹配返回 false", func() {
			m := map[string]interface{}{"a": "not_a_map"}
			So(deleteNestedField(m, []string{"a", "b"}), ShouldBeFalse)
		})
	})
}
