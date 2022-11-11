package thinkingdata

import (
	"bitbucket.org/funplus/sandwich/pkg/logbus/utils"
	"github.com/sandwich-go/zapgen/zapencoder"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
		enc.AddString("#account_id", d.AccountId)
	}
	if d.DistinctId != "" {
		enc.AddString("#distinct_id", d.DistinctId)
	}
	enc.AddString("#type", d.Type)
	enc.AddString("#time", d.Time)
	if d.EventName != "" {
		enc.AddString("#event_name", d.EventName)
	}
	if d.EventId != "" {
		enc.AddString("#event_id", d.EventId)
	}
	if d.FirstCheckId != "" {
		enc.AddString("#first_check_id", d.FirstCheckId)
	}
	if d.Ip != "" {
		enc.AddString("#ip", d.Ip)
	}
	if d.UUID != "" {
		enc.AddString("#uuid", d.UUID)
	}
	return enc.AddObject("properties", zapencoder.StringInterfaceMap(d.Properties))
}

func (d Data) MarshalAsJson() ([]byte, error) {
	var fields = make([]zap.Field, 0, 10)
	if d.AccountId != "" {
		fields = append(fields, zap.String("#account_id", d.AccountId))
	}
	if d.DistinctId != "" {
		fields = append(fields, zap.String("#distinct_id", d.DistinctId))
	}
	fields = append(fields, zap.String("#type", d.Type))
	fields = append(fields, zap.String("#time", d.Time))
	if d.EventName != "" {
		fields = append(fields, zap.String("#event_name", d.EventName))
	}
	if d.EventId != "" {
		fields = append(fields, zap.String("#event_id", d.EventId))
	}
	if d.FirstCheckId != "" {
		fields = append(fields, zap.String("#first_check_id", d.FirstCheckId))
	}
	if d.Ip != "" {
		fields = append(fields, zap.String("#ip", d.Ip))
	}
	if d.UUID != "" {
		fields = append(fields, zap.String("#uuid", d.UUID))
	}
	fields = append(fields, zap.Object("properties", zapencoder.StringInterfaceMap(d.Properties)))
	return utils.Zap2Json(fields)
}
