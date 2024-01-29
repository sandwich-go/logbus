package thinkingdata

import (
	"github.com/sandwich-go/zapgen/zapencoder"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/sandwich-go/logbus/utils"
)

var emptyData = Data{}

// Data 数据信息
type Data struct {
	AccountId    string                 `json:"#account_id,omitempty"`
	DistinctId   string                 `json:"#distinct_id,omitempty"`
	Type         string                 `json:"#type"`
	Time         string                 `json:"#time"`
	EventName    string                 `json:"#event_name,omitempty"`
	EventId      string                 `json:"#event_id,omitempty"`
	FirstCheckId string                 `json:"#first_check_id,omitempty"`
	Ip           string                 `json:"#ip,omitempty"`
	UUID         string                 `json:"#uuid,omitempty"`
	Properties   map[string]interface{} `json:"properties"`
}

func (d Data) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if d.AccountId != "" {
		enc.AddString(ACCOUNT, d.AccountId)
	}
	if d.DistinctId != "" {
		enc.AddString(DISTINCT, d.DistinctId)
	}
	enc.AddString(TYPE, d.Type)
	enc.AddString(TIME, d.Time)
	if d.EventName != "" {
		enc.AddString(EVENT, d.EventName)
	}
	if d.EventId != "" {
		enc.AddString(EVENT_ID, d.EventId)
	}
	if d.FirstCheckId != "" {
		enc.AddString(FIRST_CHECK_ID, d.FirstCheckId)
	}
	if d.Ip != "" {
		enc.AddString(IP, d.Ip)
	}
	if d.UUID != "" {
		enc.AddString(UUID, d.UUID)
	}
	return enc.AddObject("properties", zapencoder.StringInterfaceMap(d.Properties))
}

func (d Data) MarshalAsJson() ([]byte, error) {
	var fields = make([]zap.Field, 0, 10)
	if d.AccountId != "" {
		fields = append(fields, zap.String(ACCOUNT, d.AccountId))
	}
	if d.DistinctId != "" {
		fields = append(fields, zap.String(DISTINCT, d.DistinctId))
	}
	fields = append(fields, zap.String(TYPE, d.Type))
	fields = append(fields, zap.String(TIME, d.Time))
	if d.EventName != "" {
		fields = append(fields, zap.String(EVENT, d.EventName))
	}
	if d.EventId != "" {
		fields = append(fields, zap.String(EVENT_ID, d.EventId))
	}
	if d.FirstCheckId != "" {
		fields = append(fields, zap.String(FIRST_CHECK_ID, d.FirstCheckId))
	}
	if d.Ip != "" {
		fields = append(fields, zap.String(IP, d.Ip))
	}
	if d.UUID != "" {
		fields = append(fields, zap.String(UUID, d.UUID))
	}
	fields = append(fields, zap.Object("properties", zapencoder.StringInterfaceMap(d.Properties)))
	return utils.Zap2Json(fields)
}
