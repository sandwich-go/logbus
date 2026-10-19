package slg

import "github.com/sandwich-go/boost/xconv"

const EventWechatRiskRank = "wechat_risk_rank"

// WechatRiskRankOpts 微信风控打点可选属性，业务按需传入。
type WechatRiskRankOpts struct {
	FpID              string // 玩家唯一标识符。
	IsTest            bool   // 是否为测试调用，默认为 false，表示正式调用。
	ErrCode           int64  // 微信返回的错误码。
	ErrMsg            string // 微信返回的错误信息，未获取到时留空。
	CumulativeUSDPaid string // 累计充值金额，单位为 USD，选填，由调用方转换为字符串，不进行汇率换算。
}

func (o *WechatRiskRankOpts) fpid() string {
	if o == nil {
		return emptyValue
	}
	return o.FpID
}

func (o *WechatRiskRankOpts) toProperties(p map[string]interface{}) map[string]interface{} {
	p["is_test"] = "false"
	p["errcode"] = emptyValue
	p["errmsg"] = emptyValue
	p["cumulative_usdpayed"] = emptyValue
	if o == nil {
		return p
	}
	p["is_test"] = xconv.String(o.IsTest)
	p["errcode"] = xconv.String(o.ErrCode)
	p["errmsg"] = o.ErrMsg
	p["cumulative_usdpayed"] = o.CumulativeUSDPaid
	return p
}

// TrackWechatRiskRank 微信风控打点，查询对应场景并拿到结果后调用。
// 三个场景分别请求、分别打点；每次调用仅上报一个场景。
// scene 为场景值：0 表示注册，1 表示营销作弊，2 表示用户生成内容。
// roleID 为玩家角色 ID，requestID 为本次请求的唯一标识，对应微信响应中的 unoin_id，而非玩家的微信身份标识。
// riskRank 为风险等级，取值为 0 至 4，数字越大风险越高，上报时转换为字符串。
// opts.CumulativeUSDPaid 为选填的累计充值金额字符串，单位为 USD，不填时上报空字符串。
// 风险等级 0 上报为 "0"，不会被转为空字符串。
// opts 可为 nil，此时仍保留全部属性，is_test 为 "false"，错误码、错误信息和金额为空字符串。
func TrackWechatRiskRank(roleID uint64, scene int32, requestID string, riskRank int32, opts *WechatRiskRankOpts) error {
	p := map[string]interface{}{
		"scene":      xconv.String(scene),
		"request_id": requestID,
		"risk_rank":  xconv.String(riskRank),
	}
	return track(xconv.String(roleID), EventWechatRiskRank, opts.fpid(), opts.toProperties(p))
}
