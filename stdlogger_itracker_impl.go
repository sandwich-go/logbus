package logbus

import (
	"github.com/sandwich-go/logbus/utils"

	"github.com/sandwich-go/logbus/thinkingdata"

	"go.uber.org/zap"
)

func (s *StdLogger) PrintThingkingData(data thinkingdata.Data) {
	bytes, err := data.MarshalAsJsonV2()
	if err != nil {
		s.ErrorWithChannel(Setting.DefaultChannel, zap.String("PrintThingkingData", err.Error()))
	}
	s.InfoWithChannel(THINKINGDATA, zap.ByteString(MsgBody, bytes))
}

func (s *StdLogger) PrintBigQuery(tableName zap.Field, fields ...zap.Field) {
	bytes, err := utils.Zap2Json(fields)
	if err != nil {
		s.ErrorWithChannel(Setting.DefaultChannel, zap.String("PrintBigQuery", err.Error()))
	}
	fields = append([]zap.Field{tableName, zap.ByteString(MsgBody, bytes)})
	s.InfoWithChannel(BIGQUERY, fields...)
}
