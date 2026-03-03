package bi

import (
	"github.com/sandwich-go/zapgen/zapencoder"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sandwich-go/logbus/utils"
)

var emptyData = Data{}

const (
	fieldAppID = "app_id"
	fieldTs    = "ts"
	fieldRole  = "role_id"
	fieldEvent = "event"
	fieldFpID  = "fpid"
	fieldProps = "properties"
)

// Data BI 打点数据结构
// 输出格式: {"app_id":"...", "ts":1756457851660, "role_id":"...", "event":"...", "fpid":"...", "properties":{...}}
type Data struct {
	AppID      string                 `json:"app_id"`     // 必传参数 项目 BI app ID
	Ts         int64                  `json:"ts"`         // 必传参数 毫秒时间戳
	RoleID     string                 `json:"role_id"`    // role ID
	Event      string                 `json:"event"`      // 必传参数 不同事件对应的 event 名称
	FpID       string                 `json:"fpid"`       // fpid
	Properties map[string]interface{} `json:"properties"` // 自定义属性
}

// MarshalLogObject 实现 zapcore.ObjectMarshaler，用于 zap 日志
func (d Data) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString(fieldAppID, d.AppID)
	enc.AddInt64(fieldTs, d.Ts)
	enc.AddString(fieldRole, d.RoleID)
	enc.AddString(fieldEvent, d.Event)
	enc.AddString(fieldFpID, d.FpID)
	return enc.AddObject(fieldProps, zapencoder.StringInterfaceMap(d.Properties))
}

// MarshalAsJson 使用 zap 序列化为 JSON
func (d Data) MarshalAsJson() ([]byte, error) {
	return utils.Zap2Json(d.toFields())
}

// toFields 将 Data 转换为 zap.Field 切片
func (d Data) toFields() []zap.Field {
	if d.Properties == nil {
		d.Properties = make(map[string]interface{})
	}
	return []zap.Field{
		zap.String(fieldAppID, d.AppID),
		zap.Int64(fieldTs, d.Ts),
		zap.String(fieldRole, d.RoleID),
		zap.String(fieldEvent, d.Event),
		zap.String(fieldFpID, d.FpID),
		zap.Object(fieldProps, zapencoder.StringInterfaceMap(d.Properties)),
	}
}
