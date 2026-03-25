package slg

import "github.com/sandwich-go/boost/xconv"

const (
	EventAllianceCreateDismiss      = "alliance_create_dismiss"
	EventAllianceLevelUp            = "alliance_level_up"
	EventAllianceUserEvent          = "alliance_user_event"
	EventAllianceChangeLeaderEvent  = "alliance_change_leader_event"
	EventAllianceApplyLog           = "alliance_apply_log"
	EventAllianceInviteLog          = "alliance_invite_log"
	EventAllianceRecommend          = "alliance_recommend"
	EventAllianceModifyEvent        = "alliance_modify_event"
	EventAllianceMemberRankChange   = "alliance_member_rank_change_event"
	EventAllianceScienceResearch    = "alliance_science_research"
	EventAllianceScienceDonate      = "alliance_science_donate"
	EventAllianceBuildingStatus     = "alliance_building_status_change"
	EventAllianceHelpRequest        = "alliance_help_request"
	EventAllianceHelpAlly           = "alliance_help_ally"
	EventAllianceKingdomAppointment = "alliance_kingdom_appointment"
	EventAllianceTab                = "alliance_tab"
	EventAllianceRelocation         = "alliance_relocation"
	EventAllianceChest              = "alliance_chest"
)

// AllianceCreateDismissOperatorType 创建/解散联盟操作类型。
type AllianceCreateDismissOperatorType int32

const (
	// AllianceCreateDismissCreate 表示创建联盟。
	AllianceCreateDismissCreate AllianceCreateDismissOperatorType = 1
	// AllianceCreateDismissDismiss 表示主动解散联盟。
	AllianceCreateDismissDismiss AllianceCreateDismissOperatorType = 2
	// AllianceCreateDismissAutoDismiss 表示因不活跃而自动解散联盟。
	AllianceCreateDismissAutoDismiss AllianceCreateDismissOperatorType = 3
)

// AllianceUserEventOperatorType 联盟成员变动类型。
type AllianceUserEventOperatorType int32

const (
	// AllianceUserEventCreateAndJoin 表示创建联盟并加入。
	AllianceUserEventCreateAndJoin AllianceUserEventOperatorType = 1
	// AllianceUserEventApplyJoin 表示申请后成功加入联盟。
	AllianceUserEventApplyJoin AllianceUserEventOperatorType = 2
	// AllianceUserEventInviteJoin 表示接受邀请后加入联盟。
	AllianceUserEventInviteJoin AllianceUserEventOperatorType = 3
	// AllianceUserEventDirectJoin 表示直接加入联盟。
	AllianceUserEventDirectJoin AllianceUserEventOperatorType = 4
	// AllianceUserEventLeave 表示主动退出联盟。
	AllianceUserEventLeave AllianceUserEventOperatorType = 5
	// AllianceUserEventKicked 表示被踢出联盟。
	AllianceUserEventKicked AllianceUserEventOperatorType = 6
	// AllianceUserEventDismissedLeave 表示联盟解散后离开联盟。
	AllianceUserEventDismissedLeave AllianceUserEventOperatorType = 7
	// AllianceUserEventWhitelistJoin 表示通过白名单直接加入联盟。
	AllianceUserEventWhitelistJoin AllianceUserEventOperatorType = 8
)

// AllianceChangeLeaderOperatorType 盟主变更方式。
type AllianceChangeLeaderOperatorType int32

const (
	// AllianceChangeLeaderInactiveReplace 表示盟主不活跃后被替代。
	AllianceChangeLeaderInactiveReplace AllianceChangeLeaderOperatorType = 1
	// AllianceChangeLeaderTransfer 表示盟主主动转让。
	AllianceChangeLeaderTransfer AllianceChangeLeaderOperatorType = 2
	// AllianceChangeLeaderGM 表示由 GM 后台操作变更盟主。
	AllianceChangeLeaderGM AllianceChangeLeaderOperatorType = 3
)

// AllianceApplyOperatorType 联盟申请操作类型。
type AllianceApplyOperatorType int32

const (
	// AllianceApplySend 表示发起入盟申请。
	AllianceApplySend AllianceApplyOperatorType = 1
	// AllianceApplyAccept 表示同意入盟申请。
	AllianceApplyAccept AllianceApplyOperatorType = 2
	// AllianceApplyReject 表示拒绝入盟申请。
	AllianceApplyReject AllianceApplyOperatorType = 3
)

// AllianceInviteOperatorType 联盟邀请操作类型。
type AllianceInviteOperatorType int32

const (
	// AllianceInviteSend 表示发起联盟邀请。
	AllianceInviteSend AllianceInviteOperatorType = 1
	// AllianceInviteAccept 表示同意联盟邀请。
	AllianceInviteAccept AllianceInviteOperatorType = 2
	// AllianceInviteReject 表示拒绝联盟邀请。
	AllianceInviteReject AllianceInviteOperatorType = 3
)

// AllianceModifyOperatorType 联盟修改操作类型。
type AllianceModifyOperatorType int32

const (
	// AllianceModifyName 表示修改联盟名称。
	AllianceModifyName AllianceModifyOperatorType = 1
	// AllianceModifyCode 表示修改联盟简称。
	AllianceModifyCode AllianceModifyOperatorType = 2
	// AllianceModifyDeclaration 表示修改联盟宣言。
	AllianceModifyDeclaration AllianceModifyOperatorType = 3
	// AllianceModifyLanguage 表示修改联盟语言。
	AllianceModifyLanguage AllianceModifyOperatorType = 4
	// AllianceModifyNoticePublish 表示发布联盟公告。
	AllianceModifyNoticePublish AllianceModifyOperatorType = 5
	// AllianceModifyNoticeSchedule 表示预约发布联盟公告。
	AllianceModifyNoticeSchedule AllianceModifyOperatorType = 6
	// AllianceModifyNoticeScheduleEdit 表示编辑预约联盟公告。
	AllianceModifyNoticeScheduleEdit AllianceModifyOperatorType = 7
	// AllianceModifyNoticeScheduleCancel 表示取消预约联盟公告。
	AllianceModifyNoticeScheduleCancel AllianceModifyOperatorType = 8
)

// AllianceScienceResearchOperatorType 联盟科技探索操作类型。
type AllianceScienceResearchOperatorType int32

const (
	// AllianceScienceResearchStart 表示开始升级联盟科技。
	AllianceScienceResearchStart AllianceScienceResearchOperatorType = 1
	// AllianceScienceResearchComplete 表示联盟科技升级完成。
	AllianceScienceResearchComplete AllianceScienceResearchOperatorType = 2
)

// AllianceBuildingStatusOperatorType 联盟建筑状态操作类型。
type AllianceBuildingStatusOperatorType int32

const (
	// AllianceBuildingStatusBuild 表示建造联盟建筑。
	AllianceBuildingStatusBuild AllianceBuildingStatusOperatorType = 1
	// AllianceBuildingStatusRemove 表示移除联盟建筑。
	AllianceBuildingStatusRemove AllianceBuildingStatusOperatorType = 2
	// AllianceBuildingStatusGarrison 表示驻守联盟建筑。
	AllianceBuildingStatusGarrison AllianceBuildingStatusOperatorType = 3
)

// AllianceTabOperatorType 联盟标记操作类型。
type AllianceTabOperatorType int32

const (
	// AllianceTabAdd 表示添加联盟标记。
	AllianceTabAdd AllianceTabOperatorType = 1
	// AllianceTabDelete 表示删除联盟标记。
	AllianceTabDelete AllianceTabOperatorType = 2
	// AllianceTabModify 表示修改联盟标记。
	AllianceTabModify AllianceTabOperatorType = 3
)

// AllianceRelocationType 联盟迁城类型。
type AllianceRelocationType int32

const (
	// AllianceRelocationToAllianceArea 表示迁移到联盟领地内。
	AllianceRelocationToAllianceArea AllianceRelocationType = 1
	// AllianceRelocationToLeader 表示迁移到盟主附近。
	AllianceRelocationToLeader AllianceRelocationType = 2
	// AllianceRelocationAdvanced 表示高级迁城到指定位置。
	AllianceRelocationAdvanced AllianceRelocationType = 3
	// AllianceRelocationRandom 表示随机迁移到王国内其他位置。
	AllianceRelocationRandom AllianceRelocationType = 4
)

// AllianceOpts 联盟模块打点可选属性。
type AllianceOpts struct {
	FpID string
}

func (o *AllianceOpts) fpid() string {
	if o == nil {
		return ""
	}
	return o.FpID
}

func (o *AllianceOpts) toProperties(p map[string]interface{}) map[string]interface{} {
	return p
}

// TrackAllianceCreateDismiss 创建或解散联盟打点。
// roleID 为事件用户主键 roleid，按文档通常传盟主 roleid。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// ownerRoleID 为盟主角色 ID。
// allianceCode 为联盟简称。
// allianceName 为联盟名称。
// allianceContent 为联盟公告内容。
// operatorType 为创建、解散或不活跃自动解散方式。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceCreateDismiss(roleID uint64, serverID int32, allianceID uint64, ownerRoleID uint64,
	allianceCode, allianceName, allianceContent string, operatorType AllianceCreateDismissOperatorType, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"owner_role_id":    stringifyUint64(ownerRoleID),
		"server_id":        stringifyInt32(serverID),
		"alliance_id":      stringifyUint64(allianceID),
		"alliance_code":    allianceCode,
		"alliance_name":    allianceName,
		"alliance_content": allianceContent,
		"operator_type":    stringifyInt32(int32(operatorType)),
	}
	return trackAllianceEvent(roleID, EventAllianceCreateDismiss, p, opts)
}

// TrackAllianceLevelUp 联盟升级打点。
// roleID 为事件用户主键 roleid，按文档通常传盟主 roleid。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// ownerRoleID 为盟主角色 ID。
// allianceCode 为联盟简称。
// allianceName 为联盟名称。
// allianceContent 为联盟公告内容。
// allianceLevel 为升级后的联盟等级。
// memberNum 为当前联盟成员数量。
// allianceExp 为当前联盟经验值。
// alliancePower 为当前联盟战力。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceLevelUp(roleID uint64, serverID int32, allianceID uint64, ownerRoleID uint64,
	allianceCode, allianceName, allianceContent string, allianceLevel, memberNum int32, allianceExp, alliancePower int64, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"owner_role_id":    stringifyUint64(ownerRoleID),
		"server_id":        stringifyInt32(serverID),
		"alliance_id":      stringifyUint64(allianceID),
		"alliance_code":    allianceCode,
		"alliance_name":    allianceName,
		"alliance_content": allianceContent,
		"alliance_level":   stringifyInt32(allianceLevel),
		"member_num":       stringifyInt32(memberNum),
		"alliance_exp":     stringifyInt64(allianceExp),
		"alliance_power":   stringifyInt64(alliancePower),
	}
	return trackAllianceEvent(roleID, EventAllianceLevelUp, p, opts)
}

// TrackAllianceUserEvent 联盟人员变动打点。
// roleID 为事件用户主键 roleid，通常传发生变动的玩家 roleid。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// allianceCode 为联盟简称。
// allianceName 为联盟名称。
// allianceContent 为联盟公告内容。
// toMemberRoleID 为发生变动的玩家 roleid。
// operatorRoleID 为操作人的 roleid。
// 文档约定：operatorType=3 或 6 时传邀请人/操作者 roleid，其他情况默认传 0。
// operatorType 为成员变动类型，如加入、退出、踢出等。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceUserEvent(roleID uint64, serverID int32, allianceID uint64,
	allianceCode, allianceName, allianceContent string, toMemberRoleID, operatorRoleID uint64,
	operatorType AllianceUserEventOperatorType, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"server_id":        stringifyInt32(serverID),
		"alliance_id":      stringifyUint64(allianceID),
		"alliance_code":    allianceCode,
		"alliance_name":    allianceName,
		"alliance_content": allianceContent,
		"to_member_roleid": stringifyUint64(toMemberRoleID),
		"operator_roleid":  stringifyUint64(operatorRoleID),
		"operator_type":    stringifyInt32(int32(operatorType)),
	}
	return trackAllianceEvent(roleID, EventAllianceUserEvent, p, opts)
}

// TrackAllianceChangeLeaderEvent 更换盟主打点。
// roleID 为事件用户主键 roleid，按文档通常传新盟主 roleid。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// ownerRoleID 为新盟主角色 ID。
// allianceCode 为联盟简称。
// allianceName 为联盟名称。
// allianceContent 为联盟公告内容。
// ownerIDOld 为旧盟主角色 ID。
// ownerOldOffDuration 为旧盟主离线时长，单位秒。
// operatorType 为盟主变更方式，如不活跃替代、转让或 GM 操作。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceChangeLeaderEvent(roleID uint64, serverID int32, allianceID uint64, ownerRoleID uint64,
	allianceCode, allianceName, allianceContent string, ownerIDOld uint64, ownerOldOffDuration int64,
	operatorType AllianceChangeLeaderOperatorType, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"owner_role_id":          stringifyUint64(ownerRoleID),
		"server_id":              stringifyInt32(serverID),
		"alliance_id":            stringifyUint64(allianceID),
		"alliance_code":          allianceCode,
		"alliance_name":          allianceName,
		"alliance_content":       allianceContent,
		"owner_id_old":           stringifyUint64(ownerIDOld),
		"owner_old_off_duration": stringifyInt64(ownerOldOffDuration),
		"operator_type":          stringifyInt32(int32(operatorType)),
	}
	return trackAllianceEvent(roleID, EventAllianceChangeLeaderEvent, p, opts)
}

// TrackAllianceApplyLog 联盟申请打点。
// roleID 为事件用户主键 roleid。
// 文档约定：operatorType=1 时通常传申请玩家 toMemberRoleID；
// operatorType=2 或 3 时通常传审批人 operatorRoleID。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// toMemberRoleID 为申请加入联盟的玩家 roleid。
// operatorRoleID 为审批申请的操作人 roleid。
// 文档约定：仅 operatorType=2 或 3 时需要上报，其他情况默认传 0。
// operatorAllianceR 为操作人的联盟 R 级。
// 文档约定：仅 operatorType=2 或 3 时需要上报，其他情况传空字符串即可。
// operatorType 为申请操作类型，如发起申请、同意申请、拒绝申请。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceApplyLog(roleID uint64, serverID int32, allianceID uint64, toMemberRoleID uint64,
	operatorRoleID uint64, operatorAllianceR string, operatorType AllianceApplyOperatorType, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"to_member_roleid":    stringifyUint64(toMemberRoleID),
		"server_id":           stringifyInt32(serverID),
		"alliance_id":         stringifyUint64(allianceID),
		"operator_roleid":     stringifyUint64(operatorRoleID),
		"operator_alliance_r": operatorAllianceR,
		"operator_type":       stringifyInt32(int32(operatorType)),
	}
	return trackAllianceEvent(roleID, EventAllianceApplyLog, p, opts)
}

// TrackAllianceInviteLog 联盟邀请打点。
// roleID 为事件用户主键 roleid。
// 文档约定：operatorType=1 时通常传邀请人 operatorRoleID；
// operatorType=2 或 3 时通常传被邀请人 toMemberRoleID。
// serverID 为联盟所在服务器 ID。
// allianceID 为发起邀请的联盟 ID。
// toMemberRoleID 为被邀请玩家的 roleid。
// operatorRoleID 为发起邀请的操作人 roleid。
// 文档约定：仅 operatorType=1 时需要上报，其他情况默认传 0。
// operatorAllianceR 为邀请人的联盟 R 级。
// 文档约定：仅 operatorType=1 时需要上报，其他情况传空字符串即可。
// operatorType 为邀请操作类型，如发起邀请、同意邀请、拒绝邀请。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceInviteLog(roleID uint64, serverID int32, allianceID uint64, toMemberRoleID uint64,
	operatorRoleID uint64, operatorAllianceR string, operatorType AllianceInviteOperatorType, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"to_member_roleid":    stringifyUint64(toMemberRoleID),
		"server_id":           stringifyInt32(serverID),
		"alliance_id":         stringifyUint64(allianceID),
		"operator_roleid":     stringifyUint64(operatorRoleID),
		"operator_alliance_r": operatorAllianceR,
		"operator_type":       stringifyInt32(int32(operatorType)),
	}
	return trackAllianceEvent(roleID, EventAllianceInviteLog, p, opts)
}

// TrackAllianceRecommend 联盟推荐打点。
// roleID 为事件用户主键 roleid，按文档通常传盟主 roleid。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// ownerRoleID 为盟主角色 ID。
// allianceCode 为联盟简称。
// allianceName 为联盟名称。
// allianceContent 为联盟公告内容。
// allianceLevel 为联盟等级。
// memberNum 为当前联盟成员数量。
// memberNumMax 为联盟成员上限。
// alliancePower 为联盟战力。
// setLang 为联盟设置语言。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceRecommend(roleID uint64, serverID int32, allianceID uint64, ownerRoleID uint64,
	allianceCode, allianceName, allianceContent string, allianceLevel, memberNum, memberNumMax int32,
	alliancePower int64, setLang string, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"owner_role_id":    stringifyUint64(ownerRoleID),
		"server_id":        stringifyInt32(serverID),
		"alliance_id":      stringifyUint64(allianceID),
		"alliance_code":    allianceCode,
		"alliance_name":    allianceName,
		"alliance_content": allianceContent,
		"alliance_level":   stringifyInt32(allianceLevel),
		"member_num":       stringifyInt32(memberNum),
		"member_num_max":   stringifyInt32(memberNumMax),
		"alliance_power":   stringifyInt64(alliancePower),
		"set_lang":         setLang,
	}
	return trackAllianceEvent(roleID, EventAllianceRecommend, p, opts)
}

// TrackAllianceModifyEvent 联盟修改事件打点。
// roleID 为事件用户主键 roleid，通常传 operatorRoleID。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// operatorRoleID 为操作人的 roleid。
// operatorAllianceR 为操作人的联盟 R 级。
// fromGM 标记是否为 GM 后台操作。
// operatorType 为修改类型，如改名、改简称、发布公告、预约公告等。
// indexID 为公告 ID，仅 operatorType 为 5-8 时使用。
// startTime 为预约发布时间时间戳，仅 operatorType 为 5-7 时使用。
// content 为修改后的文本内容，仅 operatorType 为 1-3、5-7 时使用。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceModifyEvent(roleID uint64, serverID int32, allianceID uint64, operatorRoleID uint64, operatorAllianceR string,
	fromGM bool, operatorType AllianceModifyOperatorType, indexID uint64, startTime int64, content string, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"operator_roleid":     stringifyUint64(operatorRoleID),
		"operator_alliance_r": operatorAllianceR,
		"server_id":           stringifyInt32(serverID),
		"alliance_id":         stringifyUint64(allianceID),
		"from_gm":             xconv.String(fromGM),
		"operator_type":       stringifyInt32(int32(operatorType)),
		"index_id":            stringifyUint64(indexID),
		"start_time":          stringifyInt64(startTime),
		"content":             content,
	}
	return trackAllianceEvent(roleID, EventAllianceModifyEvent, p, opts)
}

// TrackAllianceMemberRankChangeEvent 联盟 R 级标签修改打点。
// roleID 为事件用户主键 roleid，按文档通常传被操作成员 toMemberRoleID。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// toMemberRoleID 为被修改权限的成员 roleid。
// rankNew 为新的联盟权限等级。
// rankOld 为旧的联盟权限等级。
// operatorRoleID 为操作人的 roleid。
// operatorAllianceR 为操作人的联盟 R 级。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceMemberRankChangeEvent(roleID uint64, serverID int32, allianceID uint64, toMemberRoleID uint64,
	rankNew, rankOld string, operatorRoleID uint64, operatorAllianceR string, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"to_member_roleid":    stringifyUint64(toMemberRoleID),
		"server_id":           stringifyInt32(serverID),
		"alliance_id":         stringifyUint64(allianceID),
		"rank_new":            rankNew,
		"rank_old":            rankOld,
		"operator_roleid":     stringifyUint64(operatorRoleID),
		"operator_alliance_r": operatorAllianceR,
	}
	return trackAllianceEvent(roleID, EventAllianceMemberRankChange, p, opts)
}

// TrackAllianceScienceResearch 联盟科技探索打点。
// roleID 为事件用户主键 roleid，通常传 operatorRoleID。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// operatorRoleID 为操作人的 roleid。
// operatorType 为科技操作类型，如开始升级或升级完成。
// scienceType 为科技类型。
// scienceID 为科技 ID。
// scienceLevel 为升级后的科技等级。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceScienceResearch(roleID uint64, serverID int32, allianceID uint64, operatorRoleID uint64,
	operatorType AllianceScienceResearchOperatorType, scienceType string, scienceID uint64, scienceLevel int32, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"operator_roleid": stringifyUint64(operatorRoleID),
		"server_id":       stringifyInt32(serverID),
		"alliance_id":     stringifyUint64(allianceID),
		"operator_type":   stringifyInt32(int32(operatorType)),
		"science_type":    scienceType,
		"science_id":      stringifyUint64(scienceID),
		"science_level":   stringifyInt32(scienceLevel),
	}
	return trackAllianceEvent(roleID, EventAllianceScienceResearch, p, opts)
}

// TrackAllianceScienceDonate 联盟科技捐赠打点。
// roleID 为事件用户主键 roleid，通常传 operatorRoleID。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// operatorRoleID 为操作人的 roleid。
// scienceType 为科技类型。
// scienceID 为捐赠的科技 ID。
// scienceLevel 为当前科技等级。
// donateTypeID 为捐赠材料类型 ID，由项目自定义。
// scienceExpAdd 为本次捐赠增加的进度值。
// scienceExpNew 为捐赠后的新进度值。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceScienceDonate(roleID uint64, serverID int32, allianceID uint64, operatorRoleID uint64,
	scienceType string, scienceID uint64, scienceLevel int32, donateTypeID uint64, scienceExpAdd, scienceExpNew int64, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"operator_roleid": stringifyUint64(operatorRoleID),
		"server_id":       stringifyInt32(serverID),
		"alliance_id":     stringifyUint64(allianceID),
		"science_type":    scienceType,
		"science_id":      stringifyUint64(scienceID),
		"science_level":   stringifyInt32(scienceLevel),
		"donate_type_id":  stringifyUint64(donateTypeID),
		"science_exp_add": stringifyInt64(scienceExpAdd),
		"science_exp_new": stringifyInt64(scienceExpNew),
	}
	return trackAllianceEvent(roleID, EventAllianceScienceDonate, p, opts)
}

// TrackAllianceBuildingStatusChange 联盟建筑状态变更打点。
// roleID 为事件用户主键 roleid，通常传 operatorRoleID。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// operatorRoleID 为操作人的 roleid。
// 文档约定：建造/移除时通常传盟主或管理员；驻守时可传任意联盟成员。
// x、y 为建筑所在坐标点。
// buildingID 为建筑物 ID，由项目自行定义。
// buildingName 为建筑物名称。
// buildingType 为建筑物类型。
// operatorType 为建筑操作方式，如建造、移除、驻守。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceBuildingStatusChange(roleID uint64, serverID int32, allianceID uint64, operatorRoleID uint64,
	x, y int32, buildingID uint64, buildingName, buildingType string, operatorType AllianceBuildingStatusOperatorType, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"operator_roleid": stringifyUint64(operatorRoleID),
		"server_id":       stringifyInt32(serverID),
		"alliance_id":     stringifyUint64(allianceID),
		"x":               stringifyInt32(x),
		"y":               stringifyInt32(y),
		"building_id":     stringifyUint64(buildingID),
		"building_name":   buildingName,
		"building_type":   buildingType,
		"operator_type":   stringifyInt32(int32(operatorType)),
	}
	return trackAllianceEvent(roleID, EventAllianceBuildingStatus, p, opts)
}

// TrackAllianceHelpRequest 联盟求助打点。
// roleID 为事件用户主键 roleid，通常传 operatorRoleID。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// operatorRoleID 为发起求助的玩家 roleid。
// requestHelpList 为求助详情 JSON 列表。
// 文档示例字段包括 help_id、help_type、help_max_ts。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceHelpRequest(roleID uint64, serverID int32, allianceID uint64, operatorRoleID uint64,
	requestHelpList string, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"operator_roleid":   stringifyUint64(operatorRoleID),
		"server_id":         stringifyInt32(serverID),
		"alliance_id":       stringifyUint64(allianceID),
		"request_help_list": requestHelpList,
	}
	return trackAllianceEvent(roleID, EventAllianceHelpRequest, p, opts)
}

// TrackAllianceHelpAlly 联盟帮助打点。
// roleID 为事件用户主键 roleid，通常传 operatorRoleID。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// operatorRoleID 为执行帮助的玩家 roleid。
// helpList 为帮助详情 JSON 列表。
// 文档示例字段包括 help_id、help_type、target_roleid、help_ts。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceHelpAlly(roleID uint64, serverID int32, allianceID uint64, operatorRoleID uint64,
	helpList string, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"operator_roleid": stringifyUint64(operatorRoleID),
		"server_id":       stringifyInt32(serverID),
		"alliance_id":     stringifyUint64(allianceID),
		"help_list":       helpList,
	}
	return trackAllianceEvent(roleID, EventAllianceHelpAlly, p, opts)
}

// TrackAllianceKingdomAppointment 任命王国官职打点。
// roleID 为事件用户主键 roleid，按文档通常传盟主 ownerRoleID。
// serverID 为联盟所在服务器 ID，文档备注该字段等同王国 ID。
// allianceID 为联盟 ID。
// ownerRoleID 为盟主 roleid。
// appointmentList 为任命详情 JSON 列表。
// 文档示例字段包括 job_id、target_roleid。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceKingdomAppointment(roleID uint64, serverID int32, allianceID uint64, ownerRoleID uint64,
	appointmentList string, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"owner_role_id":    stringifyUint64(ownerRoleID),
		"server_id":        stringifyInt32(serverID),
		"alliance_id":      stringifyUint64(allianceID),
		"appointment_list": appointmentList,
	}
	return trackAllianceEvent(roleID, EventAllianceKingdomAppointment, p, opts)
}

// TrackAllianceTab 联盟标记打点。
// roleID 为事件用户主键 roleid，通常传 operatorRoleID。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// operatorRoleID 为操作人的 roleid。
// operatorType 为标记操作方式，如添加、删除、修改。
// tabName 为标记名称。
// tabID 为标记 ID。
// tabType 为标记类型。
// x、y 为标记坐标点。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceTab(roleID uint64, serverID int32, allianceID uint64, operatorRoleID uint64,
	operatorType AllianceTabOperatorType, tabName string, tabID uint64, tabType string, x, y int32, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"operator_roleid": stringifyUint64(operatorRoleID),
		"server_id":       stringifyInt32(serverID),
		"alliance_id":     stringifyUint64(allianceID),
		"operator_type":   stringifyInt32(int32(operatorType)),
		"tab_name":        tabName,
		"tab_id":          stringifyUint64(tabID),
		"tab_type":        tabType,
		"x":               stringifyInt32(x),
		"y":               stringifyInt32(y),
	}
	return trackAllianceEvent(roleID, EventAllianceTab, p, opts)
}

// TrackAllianceRelocation 联盟迁城打点。
// roleID 为事件用户主键 roleid，通常传 operatorRoleID。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// operatorRoleID 为执行迁城的玩家 roleid。
// relocationType 为迁城类型，如联盟迁城、盟主迁城、高级迁城、随机迁城。
// oldPosX、oldPosY 为迁城前坐标。
// newPosX、newPosY 为迁城后坐标。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceRelocation(roleID uint64, serverID int32, allianceID uint64, operatorRoleID uint64,
	relocationType AllianceRelocationType, oldPosX, oldPosY, newPosX, newPosY int32, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"operator_roleid": stringifyUint64(operatorRoleID),
		"server_id":       stringifyInt32(serverID),
		"alliance_id":     stringifyUint64(allianceID),
		"relocation_type": stringifyInt32(int32(relocationType)),
		"old_pos_x":       stringifyInt32(oldPosX),
		"old_pos_y":       stringifyInt32(oldPosY),
		"new_pos_x":       stringifyInt32(newPosX),
		"new_pos_y":       stringifyInt32(newPosY),
	}
	return trackAllianceEvent(roleID, EventAllianceRelocation, p, opts)
}

// TrackAllianceChest 联盟宝箱打点。
// roleID 为事件用户主键 roleid，按文档通常传宝箱发起者 toMemberRoleID。
// serverID 为联盟所在服务器 ID。
// allianceID 为联盟 ID。
// toMemberRoleID 为联盟宝箱发起者的 roleid。
// chestType 为宝箱类型。
// chestID 为宝箱 ID。
// opts 为公共可选参数，当前仅支持 FpID。
func TrackAllianceChest(roleID uint64, serverID int32, allianceID uint64, toMemberRoleID uint64,
	chestType string, chestID uint64, opts *AllianceOpts) error {
	p := map[string]interface{}{
		"to_member_roleid": stringifyUint64(toMemberRoleID),
		"server_id":        stringifyInt32(serverID),
		"alliance_id":      stringifyUint64(allianceID),
		"chest_type":       chestType,
		"chest_id":         stringifyUint64(chestID),
	}
	return trackAllianceEvent(roleID, EventAllianceChest, p, opts)
}

func trackAllianceEvent(roleID uint64, event string, properties map[string]interface{}, opts *AllianceOpts) error {
	return track(xconv.String(roleID), event, opts.fpid(), opts.toProperties(properties))
}

func stringifyUint64(v uint64) string {
	if v == 0 {
		return emptyValue
	}
	return xconv.String(v)
}

func stringifyInt32(v int32) string {
	if v == 0 {
		return emptyValue
	}
	return xconv.String(v)
}

func stringifyInt64(v int64) string {
	if v == 0 {
		return emptyValue
	}
	return xconv.String(v)
}
