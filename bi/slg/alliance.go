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
	AllianceCreateDismissCreate      AllianceCreateDismissOperatorType = 1 // 创建
	AllianceCreateDismissDismiss     AllianceCreateDismissOperatorType = 2 // 解散
	AllianceCreateDismissAutoDismiss AllianceCreateDismissOperatorType = 3 // 不活跃自动解散
)

// AllianceUserEventOperatorType 联盟成员变动类型。
type AllianceUserEventOperatorType int32

const (
	AllianceUserEventCreateAndJoin  AllianceUserEventOperatorType = 1
	AllianceUserEventApplyJoin      AllianceUserEventOperatorType = 2
	AllianceUserEventInviteJoin     AllianceUserEventOperatorType = 3
	AllianceUserEventDirectJoin     AllianceUserEventOperatorType = 4
	AllianceUserEventLeave          AllianceUserEventOperatorType = 5
	AllianceUserEventKicked         AllianceUserEventOperatorType = 6
	AllianceUserEventDismissedLeave AllianceUserEventOperatorType = 7
	AllianceUserEventWhitelistJoin  AllianceUserEventOperatorType = 8
)

// AllianceChangeLeaderOperatorType 盟主变更方式。
type AllianceChangeLeaderOperatorType int32

const (
	AllianceChangeLeaderInactiveReplace AllianceChangeLeaderOperatorType = 1
	AllianceChangeLeaderTransfer        AllianceChangeLeaderOperatorType = 2
	AllianceChangeLeaderGM              AllianceChangeLeaderOperatorType = 3
)

// AllianceApplyOperatorType 联盟申请操作类型。
type AllianceApplyOperatorType int32

const (
	AllianceApplySend   AllianceApplyOperatorType = 1
	AllianceApplyAccept AllianceApplyOperatorType = 2
	AllianceApplyReject AllianceApplyOperatorType = 3
)

// AllianceInviteOperatorType 联盟邀请操作类型。
type AllianceInviteOperatorType int32

const (
	AllianceInviteSend   AllianceInviteOperatorType = 1
	AllianceInviteAccept AllianceInviteOperatorType = 2
	AllianceInviteReject AllianceInviteOperatorType = 3
)

// AllianceModifyOperatorType 联盟修改操作类型。
type AllianceModifyOperatorType int32

const (
	AllianceModifyName                 AllianceModifyOperatorType = 1
	AllianceModifyCode                 AllianceModifyOperatorType = 2
	AllianceModifyDeclaration          AllianceModifyOperatorType = 3
	AllianceModifyLanguage             AllianceModifyOperatorType = 4
	AllianceModifyNoticePublish        AllianceModifyOperatorType = 5
	AllianceModifyNoticeSchedule       AllianceModifyOperatorType = 6
	AllianceModifyNoticeScheduleEdit   AllianceModifyOperatorType = 7
	AllianceModifyNoticeScheduleCancel AllianceModifyOperatorType = 8
)

// AllianceScienceResearchOperatorType 联盟科技探索操作类型。
type AllianceScienceResearchOperatorType int32

const (
	AllianceScienceResearchStart    AllianceScienceResearchOperatorType = 1
	AllianceScienceResearchComplete AllianceScienceResearchOperatorType = 2
)

// AllianceBuildingStatusOperatorType 联盟建筑状态操作类型。
type AllianceBuildingStatusOperatorType int32

const (
	AllianceBuildingStatusBuild    AllianceBuildingStatusOperatorType = 1
	AllianceBuildingStatusRemove   AllianceBuildingStatusOperatorType = 2
	AllianceBuildingStatusGarrison AllianceBuildingStatusOperatorType = 3
)

// AllianceTabOperatorType 联盟标记操作类型。
type AllianceTabOperatorType int32

const (
	AllianceTabAdd    AllianceTabOperatorType = 1
	AllianceTabDelete AllianceTabOperatorType = 2
	AllianceTabModify AllianceTabOperatorType = 3
)

// AllianceRelocationType 联盟迁城类型。
type AllianceRelocationType int32

const (
	AllianceRelocationToAllianceArea AllianceRelocationType = 1
	AllianceRelocationToLeader       AllianceRelocationType = 2
	AllianceRelocationAdvanced       AllianceRelocationType = 3
	AllianceRelocationRandom         AllianceRelocationType = 4
)

// AllianceOpts 联盟模块打点可选属性，业务按需传入。
type AllianceOpts struct {
	FpID  string
	Extra map[string]string
}

func (o *AllianceOpts) fpid() string {
	if o == nil {
		return ""
	}
	return o.FpID
}

func (o *AllianceOpts) toProperties(p map[string]interface{}) map[string]interface{} {
	if o == nil {
		return p
	}
	for k, v := range o.Extra {
		p[k] = v
	}
	return p
}

// TrackAllianceCreateDismiss 创建或解散联盟打点。
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
