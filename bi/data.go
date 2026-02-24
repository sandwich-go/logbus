package bi

import (
	"github.com/sandwich-go/zapgen/zapencoder"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sandwich-go/logbus/utils"
)

var emptyData = Data{}

// Data BI 打点数据结构
// 输出格式: {"app_id":"...", "ts":1756457851660, "role_id":"...", "event":"...", "fpid":"...", "properties":{...}}
type Data struct {
	AppID      string                 `json:"app_id"`  // 必传参数 项目 BI app ID
	Ts         int64                  `json:"ts"`      // 必传参数 毫秒时间戳
	RoleID     string                 `json:"role_id"` // 必传参数 role ID
	Event      string                 `json:"event"`   // 必传参数 不同事件对应的 event 名称
	Fpid       string                 `json:"fpid"`
	Properties map[string]interface{} `json:"properties"` // 自定义属性
}

// MarshalLogObject 实现 zapcore.ObjectMarshaler，用于 zap 日志
func (d Data) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("app_id", d.AppID)
	enc.AddInt64("ts", d.Ts)
	enc.AddString("role_id", d.RoleID)
	enc.AddString("event", d.Event)
	enc.AddString("fpid", d.Fpid)
	return enc.AddObject("properties", zapencoder.StringInterfaceMap(d.Properties))
}

// MarshalAsJson 使用 zap 序列化为 JSON
func (d Data) MarshalAsJson() ([]byte, error) {
	if d.Properties == nil {
		d.Properties = make(map[string]interface{})
	}
	fields := []zap.Field{
		zap.String("app_id", d.AppID),
		zap.Int64("ts", d.Ts),
		zap.String("role_id", d.RoleID),
		zap.String("event", d.Event),
		zap.String("fpid", d.Fpid),
		zap.Object("properties", zapencoder.StringInterfaceMap(d.Properties)),
	}
	return utils.Zap2Json(fields)
}
