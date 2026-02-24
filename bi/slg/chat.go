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
}

func (o *UserChatOpts) fpid() string {
	if o == nil {
		return ""
	}
	return o.FpID
}

func (o *UserChatOpts) toProperties(p map[string]interface{}) map[string]interface{} {
	if o == nil {
		return nil
	}
	p["ip"] = o.IP
	p["server_id"] = xconv.String(o.ServerID)
	p["gameusername"] = o.GameUser
	p["total_power"] = xconv.String(o.TotalPower)
	p["trans_lang"] = o.TransLang
	return p
}

// TrackUserChatPrivate 私聊打点
// chat_type=1, to_roleid=recipientRoleID, chat_room_id=0, channel=0
func TrackUserChatPrivate(roleID uint64, toRoleID uint64, content string, opts *UserChatOpts) error {
	p := map[string]interface{}{
		"chat_type":    chatTypePrivate,
		"channel":      chatChannelPrivate,
		"to_roleid":    xconv.String(toRoleID),
		"chat_room_id": emptyValue,
		"alliance_id":  emptyValue,
		"content":      content,
	}
	return track(xconv.String(roleID), EventUserChat, opts.fpid(), opts.toProperties(p))
}

// TrackUserChatWorldChannel 频道-世界打点
// chat_type=2, channel=1(世界), to_roleid=0, chat_room_id=0
func TrackUserChatWorldChannel(roleID uint64, content string, opts *UserChatOpts) error {
	p := map[string]interface{}{
		"chat_type":    chatTypeChannel,
		"channel":      chatChannelWorld,
		"to_roleid":    emptyValue,
		"chat_room_id": emptyValue,
		"alliance_id":  emptyValue,
		"content":      content,
	}
	return track(xconv.String(roleID), EventUserChat, opts.fpid(), opts.toProperties(p))
}

// TrackUserChatAllianceChannel 频道-联盟打点
// chat_type=2, channel=2(联盟)
func TrackUserChatAllianceChannel(roleID uint64, allianceID uint64, content string, opts *UserChatOpts) error {
	p := map[string]interface{}{
		"chat_type":    chatTypeChannel,
		"channel":      chatChannelAlliance,
		"to_roleid":    emptyValue,
		"chat_room_id": emptyValue,
		"alliance_id":  xconv.String(allianceID),
		"content":      content,
	}
	return track(xconv.String(roleID), EventUserChat, opts.fpid(), opts.toProperties(p))
}

// TrackUserChatAllianceRoom 群组-联盟聊天室打点
// chat_type=3, channel=1(联盟聊天室), chat_room_id=allianceID, to_roleid=0
func TrackUserChatAllianceRoom(roleID uint64, allianceID uint64, chatRoomID, content string, opts *UserChatOpts) error {
	p := map[string]interface{}{
		"chat_type":    chatTypeGroup,
		"channel":      chatChannelAllianceRoom,
		"alliance_id":  xconv.String(allianceID),
		"chat_room_id": chatRoomID,
		"to_roleid":    emptyValue,
		"content":      content,
	}
	return track(xconv.String(roleID), EventUserChat, opts.fpid(), opts.toProperties(p))
}

// TrackUserChatOtherGroup 群组-其他群聊打点
// chat_type=3, channel=2(其他群聊), chat_room_id=roomID, to_roleid=0
func TrackUserChatOtherGroup(roleID uint64, chatRoomID, content string, opts *UserChatOpts) error {
	p := map[string]interface{}{
		"chat_type":    chatTypeGroup,
		"channel":      chatChannelOtherGroup,
		"chat_room_id": chatRoomID,
		"alliance_id":  emptyValue,
		"to_roleid":    emptyValue,
		"content":      content,
	}
	return track(xconv.String(roleID), EventUserChat, opts.fpid(), opts.toProperties(p))
}

// TrackUserChatGroupSend 群发打点
// chat_type=4, channel=1(群发), to_roleid=0, chat_room_id=0
func TrackUserChatGroupSend(roleID uint64, content string, opts *UserChatOpts) error {
	p := map[string]interface{}{
		"chat_type":    chatTypeGroupSend,
		"channel":      chatChannelGroupSend,
		"to_roleid":    emptyValue,
		"chat_room_id": emptyValue,
		"alliance_id":  emptyValue,
		"content":      content,
	}
	return track(xconv.String(roleID), EventUserChat, opts.fpid(), opts.toProperties(p))
}
