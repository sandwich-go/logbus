package slg

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/sandwich-go/logbus/bi"
)

func TestTrackWechatRiskRank(t *testing.T) {
	cases := []struct {
		name           string
		scene          int32
		requestID      string
		riskRank       int32
		opts           *WechatRiskRankOpts
		wantFpID       string
		wantProperties map[string]interface{}
	}{
		{
			name:  "不传可选参数时金额留空",
			scene: 0, requestID: "request-no-opts", riskRank: 0,
			wantProperties: map[string]interface{}{
				"scene": "0", "is_test": "false", "errcode": "", "errmsg": "",
				"request_id": "request-no-opts", "risk_rank": "0", "cumulative_usdpayed": "",
			},
		},
		{
			name:  "注册场景保留零风险和零金额",
			scene: 0, requestID: "request-register", riskRank: 0,
			opts: &WechatRiskRankOpts{CumulativeUSDPaid: "0"},
			wantProperties: map[string]interface{}{
				"scene": "0", "is_test": "false", "errcode": "0", "errmsg": "",
				"request_id": "request-register", "risk_rank": "0", "cumulative_usdpayed": "0",
			},
		},
		{
			name:  "营销作弊场景透传测试标记和金额",
			scene: 1, requestID: "request-marketing", riskRank: 4,
			opts:     &WechatRiskRankOpts{FpID: "fp-123", IsTest: true, ErrCode: 0, ErrMsg: "ok", CumulativeUSDPaid: "1234567890123456.78"},
			wantFpID: "fp-123",
			wantProperties: map[string]interface{}{
				"scene": "1", "is_test": "true", "errcode": "0", "errmsg": "ok",
				"request_id": "request-marketing", "risk_rank": "4", "cumulative_usdpayed": "1234567890123456.78",
			},
		},
		{
			name:  "用户内容场景默认正式调用",
			scene: 2, requestID: "request-ugc", riskRank: 2,
			opts: &WechatRiskRankOpts{CumulativeUSDPaid: "10.50"},
			wantProperties: map[string]interface{}{
				"scene": "2", "is_test": "false", "errcode": "0", "errmsg": "",
				"request_id": "request-ugc", "risk_rank": "2", "cumulative_usdpayed": "10.50",
			},
		},
		{
			name:  "查询失败透传负数错误码",
			scene: 0, requestID: "request-failed",
			opts: &WechatRiskRankOpts{ErrCode: -1, ErrMsg: "系统繁忙"},
			wantProperties: map[string]interface{}{
				"scene": "0", "is_test": "false", "errcode": "-1", "errmsg": "系统繁忙",
				"request_id": "request-failed", "risk_rank": "0", "cumulative_usdpayed": "",
			},
		},
		{
			name:  "完整保留六十四位整数错误码",
			scene: 0, requestID: "request-large-error",
			opts: &WechatRiskRankOpts{ErrCode: 9223372036854775807, ErrMsg: "大错误码"},
			wantProperties: map[string]interface{}{
				"scene": "0", "is_test": "false", "errcode": "9223372036854775807", "errmsg": "大错误码",
				"request_id": "request-large-error", "risk_rank": "0", "cumulative_usdpayed": "",
			},
		},
	}

	previousTracker := tracker
	t.Cleanup(func() { tracker = previousTracker })
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ct := &captureTracker{}
			MustInitialize(ct)
			if err := TrackWechatRiskRank(123, tc.scene, tc.requestID, tc.riskRank, tc.opts); err != nil {
				t.Fatalf("上报微信风控点位失败：%v", err)
			}
			data := ct.getData()
			if data == nil {
				t.Fatal("没有收到微信风控点位")
			}
			if data.Event != "wechat_risk_rank" || data.RoleID != "123" || data.FpID != tc.wantFpID || data.Ts <= 0 {
				t.Fatalf("事件公共字段不正确：%+v", data)
			}
			if !reflect.DeepEqual(data.Properties, tc.wantProperties) {
				t.Fatalf("点位属性不正确：实际 %#v，预期 %#v", data.Properties, tc.wantProperties)
			}
			encoded, err := data.MarshalAsJson()
			if err != nil {
				t.Fatalf("序列化微信风控点位失败：%v", err)
			}
			var decoded bi.Data
			if err := json.Unmarshal(encoded, &decoded); err != nil {
				t.Fatalf("解析微信风控点位失败：%v", err)
			}
			if !reflect.DeepEqual(decoded.Properties, tc.wantProperties) {
				t.Fatalf("序列化后属性不正确：实际 %#v，预期 %#v", decoded.Properties, tc.wantProperties)
			}
		})
	}
}

type failingWechatRiskTracker struct {
	captureTracker
	err error
}

func (c *failingWechatRiskTracker) TrackWithBIData(bi.Data) error {
	return c.err
}

func TestTrackWechatRiskRankReturnsTrackerError(t *testing.T) {
	previousTracker := tracker
	t.Cleanup(func() { tracker = previousTracker })
	wantErr := errors.New("上报失败")
	MustInitialize(&failingWechatRiskTracker{err: wantErr})
	if err := TrackWechatRiskRank(123, 0, "request-error", 0, nil); !errors.Is(err, wantErr) {
		t.Fatalf("上报错误未透传：实际 %v，预期 %v", err, wantErr)
	}
}
