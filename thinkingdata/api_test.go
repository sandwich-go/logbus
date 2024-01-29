package thinkingdata

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUser(t *testing.T) {
	Convey("User set", t, func() {
		properties := map[string]interface{}{"#ip": "10.0.0.1", "player_name": "zhang san", "level": 9}
		data, err := User("111", "", USER_SET, properties)
		So(err, ShouldBeNil)
		So(data.Time, ShouldNotBeEmpty)
		So(data.Type, ShouldEqual, USER_SET)
		So(data.Properties["#ip"], ShouldBeNil)
		So(data.Ip, ShouldEqual, "10.0.0.1")
		So(data.Properties["player_name"], ShouldEqual, "zhang san")
		So(data.Properties["level"], ShouldEqual, 9)
	})
}

func TestTrack(t *testing.T) {
	Convey("User track", t, func() {
		properties := map[string]interface{}{"#ip": "10.0.0.1", "player_name": "zhang san", "level": 9}
		data, err := Track("111", "", "login", "", properties)
		So(err, ShouldBeNil)
		So(data.Time, ShouldNotBeEmpty)
		So(data.Type, ShouldEqual, TRACK)
		So(data.EventName, ShouldEqual, "login")
		So(data.Properties["#ip"], ShouldBeNil)
		So(data.Ip, ShouldEqual, "10.0.0.1")
		So(data.Properties["player_name"], ShouldEqual, "zhang san")
		So(data.Properties["level"], ShouldEqual, 9)
	})
}
