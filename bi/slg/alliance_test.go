package slg

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTrackAllianceEventsRequiredFields(t *testing.T) {
	type tc struct {
		name   string
		event  string
		fields []string
		fn     func() error
	}
	cases := []tc{
		{name: "创建或解散联盟", event: EventAllianceCreateDismiss, fields: []string{"owner_role_id", "server_id", "alliance_id", "alliance_code", "alliance_name", "alliance_content", "operator_type"}, fn: func() error {
			return TrackAllianceCreateDismiss(11, 1001, 2002, 101, "AC", "A-Team", "hello", AllianceCreateDismissCreate, nil)
		}},
		{name: "联盟升级", event: EventAllianceLevelUp, fields: []string{"owner_role_id", "server_id", "alliance_id", "alliance_code", "alliance_name", "alliance_content", "alliance_level", "member_num", "alliance_exp", "alliance_power"}, fn: func() error {
			return TrackAllianceLevelUp(11, 1001, 2002, 101, "AC", "A-Team", "hello", 2, 30, 100, 999, nil)
		}},
		{name: "联盟人员变动", event: EventAllianceUserEvent, fields: []string{"server_id", "alliance_id", "alliance_code", "alliance_name", "alliance_content", "to_member_roleid", "operator_roleid", "operator_type"}, fn: func() error {
			return TrackAllianceUserEvent(11, 1001, 2002, "AC", "A-Team", "hello", 10001, 10002, AllianceUserEventInviteJoin, nil)
		}},
		{name: "更换盟主", event: EventAllianceChangeLeaderEvent, fields: []string{"owner_role_id", "server_id", "alliance_id", "alliance_code", "alliance_name", "alliance_content", "owner_id_old", "owner_old_off_duration", "operator_type"}, fn: func() error {
			return TrackAllianceChangeLeaderEvent(11, 1001, 2002, 101, "AC", "A-Team", "hello", 100, 120, AllianceChangeLeaderTransfer, nil)
		}},
		{name: "联盟申请", event: EventAllianceApplyLog, fields: []string{"to_member_roleid", "server_id", "alliance_id", "operator_roleid", "operator_alliance_r", "operator_type"}, fn: func() error {
			return TrackAllianceApplyLog(11, 1001, 2002, 10001, 10002, "R4", AllianceApplySend, nil)
		}},
		{name: "联盟邀请", event: EventAllianceInviteLog, fields: []string{"to_member_roleid", "server_id", "alliance_id", "operator_roleid", "operator_alliance_r", "operator_type"}, fn: func() error {
			return TrackAllianceInviteLog(11, 1001, 2002, 10001, 10002, "R4", AllianceInviteSend, nil)
		}},
		{name: "联盟推荐", event: EventAllianceRecommend, fields: []string{"owner_role_id", "server_id", "alliance_id", "alliance_code", "alliance_name", "alliance_content", "alliance_level", "member_num", "member_num_max", "alliance_power", "set_lang"}, fn: func() error {
			return TrackAllianceRecommend(11, 1001, 2002, 101, "AC", "A-Team", "hello", 2, 30, 50, 999, "en", nil)
		}},
		{name: "联盟修改事件", event: EventAllianceModifyEvent, fields: []string{"operator_roleid", "operator_alliance_r", "server_id", "alliance_id", "from_gm", "operator_type", "index_id", "start_time", "content"}, fn: func() error {
			return TrackAllianceModifyEvent(11, 1001, 2002, 10001, "R4", false, AllianceModifyName, 10, 1756457851660, "new", nil)
		}},
		{name: "联盟R级标签修改", event: EventAllianceMemberRankChange, fields: []string{"to_member_roleid", "server_id", "alliance_id", "rank_new", "rank_old", "operator_roleid", "operator_alliance_r"}, fn: func() error {
			return TrackAllianceMemberRankChangeEvent(11, 1001, 2002, 10001, "R3", "R4", 10002, "R5", nil)
		}},
		{name: "联盟科技探索", event: EventAllianceScienceResearch, fields: []string{"operator_roleid", "server_id", "alliance_id", "operator_type", "science_type", "science_id", "science_level"}, fn: func() error {
			return TrackAllianceScienceResearch(11, 1001, 2002, 10001, AllianceScienceResearchStart, "military", 5001, 3, nil)
		}},
		{name: "联盟科技捐赠", event: EventAllianceScienceDonate, fields: []string{"operator_roleid", "server_id", "alliance_id", "science_type", "science_id", "science_level", "donate_type_id", "science_exp_add", "science_exp_new"}, fn: func() error {
			return TrackAllianceScienceDonate(11, 1001, 2002, 10001, "military", 5001, 3, 7001, 10, 100, nil)
		}},
		{name: "联盟建筑状态变更", event: EventAllianceBuildingStatus, fields: []string{"operator_roleid", "server_id", "alliance_id", "x", "y", "building_id", "building_name", "building_type", "operator_type"}, fn: func() error {
			return TrackAllianceBuildingStatusChange(11, 1001, 2002, 10001, 1, 2, 9001, "tower", "def", AllianceBuildingStatusBuild, nil)
		}},
		{name: "联盟求助", event: EventAllianceHelpRequest, fields: []string{"operator_roleid", "server_id", "alliance_id", "request_help_list"}, fn: func() error {
			return TrackAllianceHelpRequest(11, 1001, 2002, 10001, "[{\"help_id\":325}]", nil)
		}},
		{name: "联盟帮助", event: EventAllianceHelpAlly, fields: []string{"operator_roleid", "server_id", "alliance_id", "help_list"}, fn: func() error {
			return TrackAllianceHelpAlly(11, 1001, 2002, 10001, "[{\"help_id\":325}]", nil)
		}},
		{name: "任命王国官职", event: EventAllianceKingdomAppointment, fields: []string{"owner_role_id", "server_id", "alliance_id", "appointment_list"}, fn: func() error {
			return TrackAllianceKingdomAppointment(11, 1001, 2002, 101, "[{\"job_id\":1}]", nil)
		}},
		{name: "联盟标记", event: EventAllianceTab, fields: []string{"operator_roleid", "server_id", "alliance_id", "operator_type", "tab_name", "tab_id", "tab_type", "x", "y"}, fn: func() error {
			return TrackAllianceTab(11, 1001, 2002, 10001, AllianceTabAdd, "n1", 8001, "type1", 1, 2, nil)
		}},
		{name: "联盟迁城", event: EventAllianceRelocation, fields: []string{"operator_roleid", "server_id", "alliance_id", "relocation_type", "old_pos_x", "old_pos_y", "new_pos_x", "new_pos_y"}, fn: func() error {
			return TrackAllianceRelocation(11, 1001, 2002, 10001, AllianceRelocationAdvanced, 11, 22, 33, 44, nil)
		}},
		{name: "联盟宝箱", event: EventAllianceChest, fields: []string{"to_member_roleid", "server_id", "alliance_id", "chest_type", "chest_id"}, fn: func() error {
			return TrackAllianceChest(11, 1001, 2002, 10001, "rare", 6001, nil)
		}},
	}

	for _, c := range cases {
		c := c
		Convey(c.name, t, func() {
			ct := &captureTracker{}
			MustInitialize(ct)
			So(c.fn(), ShouldBeNil)
			data := ct.getData()
			So(data, ShouldNotBeNil)
			So(data.Event, ShouldEqual, c.event)
			for _, k := range c.fields {
				v, ok := data.Properties[k]
				So(ok, ShouldBeTrue)
				So(v, ShouldHaveSameTypeAs, "")
			}
		})
	}
}

func TestTrackAllianceOptsAndExtra(t *testing.T) {
	Convey("FpID + Extra 透传", t, func() {
		ct := &captureTracker{}
		MustInitialize(ct)

		opts := &AllianceOpts{
			FpID: "fp_1",
			Extra: map[string]string{
				"custom_field": "custom_value",
			},
		}

		So(TrackAllianceCreateDismiss(11, 1001, 2002, 101, "AC", "A-Team", "hello", AllianceCreateDismissDismiss, opts), ShouldBeNil)
		data := ct.getData()
		So(data, ShouldNotBeNil)
		So(data.FpID, ShouldEqual, "fp_1")
		So(data.Properties["owner_role_id"], ShouldEqual, "101")
		So(data.Properties["operator_type"], ShouldEqual, "2")
		So(data.Properties["custom_field"], ShouldEqual, "custom_value")
	})
}

func TestTrackAllianceEmptyNumericToEmptyString(t *testing.T) {
	Convey("未获取到 numeric 字段时保留空字符串", t, func() {
		ct := &captureTracker{}
		MustInitialize(ct)

		So(TrackAllianceCreateDismiss(11, 0, 0, 0, "", "", "", 0, nil), ShouldBeNil)
		data := ct.getData()
		So(data.Properties["owner_role_id"], ShouldEqual, "")
		So(data.Properties["server_id"], ShouldEqual, "")
		So(data.Properties["alliance_id"], ShouldEqual, "")
		So(data.Properties["operator_type"], ShouldEqual, "")
	})
}
