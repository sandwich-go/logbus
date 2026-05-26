package slg

import (
	"github.com/sandwich-go/boost/xconv"
)

const (
	// chatType 聊天类型 (内部使用)
	chatTypePrivate   = "1" // 私聊
	chatTypeChannel   = "2" // 频道聊天
	chatTypeGroup     = "3" // 聊天群组
	chatTypeGroupSend = "4" // 群发聊天

	// chatChannel 频道 (chat_type=1 私聊)
	chatChannelPrivate = "1" // 私聊

	// chatChannel 频道 (chat_type=2 频道聊天)
	chatChannelWorld    = "1" // 世界
	chatChannelAlliance = "2" // 联盟

	// chatChannel 频道 (chat_type=3 聊天群组)
	chatChannelAllianceRoom = "1" // 联盟聊天室
	chatChannelOtherGroup   = "2" // 其他群聊群组

	// chatChannel 频道 (chat_type=4 群发)
	chatChannelGroupSend = "1" // 群发
)

// UserChatOpts 聊天打点可选属性，业务按需传入
type UserChatOpts struct {
	FpID       string  // 玩家唯一标识符
	IP         string  // 当前玩家 IP
	ServerID   int32   // 玩家所在服务器 ID
	GameUser   string  // 玩家昵称
	TotalPower float64 // 玩家战力
	TransLang  string  // 发送时使用的语言
	TownLvl    int32   // 玩家主城等级
}

func (o *UserChatOpts) fpid() string {
	if o == nil {
		return ""
	}
	return o.FpID
}

func (o *UserChatOpts) toProperties(p map[string]interface{}) map[string]interface{} {
	if o == nil {
		return p
	}
	p["ip"] = o.IP
	p["server_id"] = xconv.String(o.ServerID)
	p["gameusername"] = o.GameUser
	p["total_power"] = xconv.String(o.TotalPower)
	p["trans_lang"] = o.TransLang
	p["town_lvl"] = xconv.String(o.TownLvl)
	return p
}

// TrackUserChatPrivate 私聊打点
// chat_type=1, to_roleid=recipientRoleID, chat_room_id=0, channel=0
func TrackUserChatPrivate(roleID uint64, toRoleID uint64, content string, opts *UserChatOpts) error {
	return trackUserChat(roleID, chatTypePrivate, chatChannelPrivate,
		xconv.String(toRoleID), emptyValue, emptyValue, content, opts)
}

// TrackUserChatWorldChannel 频道-世界打点
// chat_type=2, channel=1(世界), to_roleid=0, chat_room_id=0
func TrackUserChatWorldChannel(roleID uint64, content string, opts *UserChatOpts) error {
	return trackUserChat(roleID, chatTypeChannel, chatChannelWorld,
		emptyValue, emptyValue, emptyValue, content, opts)
}

// TrackUserChatAllianceChannel 频道-联盟打点
// chat_type=2, channel=2(联盟)
func TrackUserChatAllianceChannel(roleID uint64, allianceID uint64, content string, opts *UserChatOpts) error {
	return trackUserChat(roleID, chatTypeChannel, chatChannelAlliance,
		emptyValue, emptyValue, xconv.String(allianceID), content, opts)
}

// TrackUserChatAllianceRoom 群组-联盟聊天室打点
// chat_type=3, channel=1(联盟聊天室), chat_room_id=allianceID, to_roleid=0
func TrackUserChatAllianceRoom(roleID uint64, allianceID uint64, chatRoomID, content string, opts *UserChatOpts) error {
	return trackUserChat(roleID, chatTypeGroup, chatChannelAllianceRoom,
		emptyValue, chatRoomID, xconv.String(allianceID), content, opts)
}

// TrackUserChatOtherGroup 群组-其他群聊打点
// chat_type=3, channel=2(其他群聊), chat_room_id=roomID, to_roleid=0
func TrackUserChatOtherGroup(roleID uint64, chatRoomID, content string, opts *UserChatOpts) error {
	return trackUserChat(roleID, chatTypeGroup, chatChannelOtherGroup,
		emptyValue, chatRoomID, emptyValue, content, opts)
}

// TrackUserChatGroupSend 群发打点
// chat_type=4, channel=1(群发), to_roleid=0, chat_room_id=0
func TrackUserChatGroupSend(roleID uint64, content string, opts *UserChatOpts) error {
	return trackUserChat(roleID, chatTypeGroupSend, chatChannelGroupSend,
		emptyValue, emptyValue, emptyValue, content, opts)
}

// trackUserChat 通用的聊天打点函数
func trackUserChat(roleID uint64, chatType, channel, toRoleID, chatRoomID, allianceID, content string, opts *UserChatOpts) error {
	p := map[string]interface{}{
		"chat_type":    chatType,
		"channel":      channel,
		"to_roleid":    toRoleID,
		"chat_room_id": chatRoomID,
		"alliance_id":  allianceID,
		"content":      content,
	}
	return track(xconv.String(roleID), EventUserChat, opts.fpid(), opts.toProperties(p))
}
