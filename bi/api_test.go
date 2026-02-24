package bi

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestTrack(t *testing.T) {
	Convey("Track", t, func() {
		properties := map[string]interface{}{
			"gameserver_id": "1",
			"alliance_id":   "1",
		}
		data, err := Track("sow.global.prod", "3555555", "alliance_create", "122222", properties)
		So(err, ShouldBeNil)
		So(data.AppID, ShouldEqual, "sow.global.prod")
		So(data.RoleID, ShouldEqual, "3555555")
		So(data.Event, ShouldEqual, "alliance_create")
		So(data.Fpid, ShouldEqual, "122222")
		So(data.Ts, ShouldNotBeZeroValue)
		So(data.Properties["gameserver_id"], ShouldEqual, "1")
		So(data.Properties["alliance_id"], ShouldEqual, "1")

		// 验证 JSON 输出格式
		bytes, err := data.MarshalAsJson()
		So(err, ShouldBeNil)

		var decoded map[string]interface{}
		So(json.Unmarshal(bytes, &decoded), ShouldBeNil)
		So(decoded["app_id"], ShouldEqual, "sow.global.prod")
		So(decoded["event"], ShouldEqual, "alliance_create")
		So(decoded["role_id"], ShouldEqual, "3555555")
		So(decoded["fpid"], ShouldEqual, "122222")
		So(decoded["ts"], ShouldNotBeNil)
		props, ok := decoded["properties"].(map[string]interface{})
		So(ok, ShouldBeTrue)
		So(props["gameserver_id"], ShouldEqual, "1")
		So(props["alliance_id"], ShouldEqual, "1")
	})
}

func TestTrack_NilProperties(t *testing.T) {
	Convey("Track with nil properties", t, func() {
		data, err := Track("sow.global.prod", "3555555", "login", "122222", nil)
		So(err, ShouldBeNil)
		So(data.Properties, ShouldNotBeNil)
		So(len(data.Properties), ShouldEqual, 0)

		bytes, err := data.MarshalAsJson()
		So(err, ShouldBeNil)
		So(string(bytes), ShouldContainSubstring, `"properties":{}`)
	})
}

func TestTrack_RequiredParams(t *testing.T) {
	Convey("Track validates required params", t, func() {
		_, err := Track("", "3555555", "event", "122222", nil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "app_id")

		_, err = Track("app", "", "event", "122222", nil)
		So(err, ShouldBeNil)

		_, err = Track("app", "role", "", "122222", nil)
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "event")

		_, err = Track("app", "role", "event", "", nil)
		So(err, ShouldBeNil)
	})
}

func TestExtractEncoder(t *testing.T) {
	Convey("ExtractEncoder", t, func() {
		enc := zapcore.NewMapObjectEncoder()
		zap.String("app_id", "sow.global.prod").AddTo(enc)
		zap.String("role_id", "3555555").AddTo(enc)
		zap.String("event", "alliance_create").AddTo(enc)
		zap.String("fpid", "122222").AddTo(enc)
		zap.String("gameserver_id", "1").AddTo(enc)
		zap.String("alliance_id", "1").AddTo(enc)

		data, err := ExtractEncoder(enc)
		So(err, ShouldBeNil)
		So(data.AppID, ShouldEqual, "sow.global.prod")
		So(data.RoleID, ShouldEqual, "3555555")
		So(data.Event, ShouldEqual, "alliance_create")
		So(data.Fpid, ShouldEqual, "122222")
		So(data.Properties["gameserver_id"], ShouldEqual, "1")
		So(data.Properties["alliance_id"], ShouldEqual, "1")
	})
}

func TestExtractEncoder_MissingRequired(t *testing.T) {
	Convey("ExtractEncoder with missing required fields", t, func() {
		enc := zapcore.NewMapObjectEncoder()
		zap.String("app_id", "app").AddTo(enc)
		// missing role_id, event, fpid
		_, err := ExtractEncoder(enc)
		So(err, ShouldNotBeNil)
	})
}
