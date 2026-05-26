package slg

import (
	"sync"
	"testing"

	"github.com/sandwich-go/logbus"
	"github.com/sandwich-go/logbus/bi"
	"github.com/sandwich-go/logbus/thinkingdata"

	. "github.com/smartystreets/goconvey/convey"
)

type captureTracker struct {
	mu     sync.Mutex
	biData *bi.Data
}

func (c *captureTracker) Track(...logbus.Field) error              { return nil }
func (c *captureTracker) TrackWithTGAData(thinkingdata.Data) error { return nil }
func (c *captureTracker) TrackWithBIData(data bi.Data) error {
	c.mu.Lock()
	c.biData = &data
	c.mu.Unlock()
	return nil
}

func (c *captureTracker) getData() *bi.Data {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.biData
}

func TestTrackUserChatPrivate(t *testing.T) {
	Convey("私聊打点", t, func() {
		ct := &captureTracker{}
		MustInitialize(ct)
		So(TrackUserChatPrivate(1, 2, "你好", nil), ShouldBeNil)
		data := ct.getData()
		So(data, ShouldNotBeNil)
		So(data.Event, ShouldEqual, EventUserChat)
		So(data.Properties["chat_type"], ShouldEqual, chatTypePrivate)
		So(data.Properties["to_roleid"], ShouldEqual, "2")
		So(data.Properties["chat_room_id"], ShouldEqual, emptyValue)
		So(data.Properties["channel"], ShouldEqual, chatChannelPrivate)
		So(data.Properties["content"], ShouldEqual, "你好")
	})
}

func TestTrackUserChatWorldChannel(t *testing.T) {
	Convey("频道-世界打点", t, func() {
		ct := &captureTracker{}
		MustInitialize(ct)
		So(TrackUserChatWorldChannel(1, "世界消息", nil), ShouldBeNil)
		data := ct.getData()
		So(data, ShouldNotBeNil)
		So(data.Properties["chat_type"], ShouldEqual, chatTypeChannel)
		So(data.Properties["channel"], ShouldEqual, chatChannelWorld)
		So(data.Properties["to_roleid"], ShouldEqual, emptyValue)
	})
}

func TestTrackUserChatAllianceRoom(t *testing.T) {
	Convey("群组-联盟聊天室打点", t, func() {
		ct := &captureTracker{}
		MustInitialize(ct)
		So(TrackUserChatAllianceRoom(1, 123, "ally_123", "联盟消息", nil), ShouldBeNil)
		data := ct.getData()
		So(data, ShouldNotBeNil)
		So(data.Properties["chat_type"], ShouldEqual, chatTypeGroup)
		So(data.Properties["channel"], ShouldEqual, chatChannelAllianceRoom)
		So(data.Properties["alliance_id"], ShouldEqual, "123")
		So(data.Properties["chat_room_id"], ShouldEqual, "ally_123")
	})
}

func TestTrackUserChatOpts(t *testing.T) {
	Convey("可选属性透传", t, func() {
		ct := &captureTracker{}
		MustInitialize(ct)
		opts := &UserChatOpts{
			FpID:       "fp1",
			IP:         "1.2.3.4",
			ServerID:   1,
			GameUser:   "玩家A",
			TotalPower: 10000.1,
			TransLang:  "zh",
			TownLvl:    25,
		}
		So(TrackUserChatPrivate(1, 2, "hi", opts), ShouldBeNil)
		data := ct.getData()
		So(data, ShouldNotBeNil)
		So(data.FpID, ShouldEqual, "fp1")
		So(data.Properties["ip"], ShouldEqual, "1.2.3.4")
		So(data.Properties["server_id"], ShouldEqual, "1")
		So(data.Properties["gameusername"], ShouldEqual, "玩家A")
		So(data.Properties["total_power"], ShouldEqual, "10000.1")
		So(data.Properties["trans_lang"], ShouldEqual, "zh")
		So(data.Properties["town_lvl"], ShouldEqual, "25")
	})
}
