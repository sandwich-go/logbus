package logbus

import (
	"github.com/sandwich-go/logbus/bi"
	"github.com/sandwich-go/logbus/thinkingdata"
	"github.com/sandwich-go/logbus/utils"
	"go.uber.org/zap"
)

func (s *StdLogger) PrintThingkingData(data thinkingdata.Data) {
	err := data.WithJSONV2(func(bytes []byte) {
		// TrackWithChannel 返回前会完成 ByteString 编码；不能把这段借用的切片
		// 交给 goroutine 或其他异步消费者。
		s.TrackWithChannel(THINKINGDATA, zap.ByteString(MsgBody, bytes))
	})
	if err != nil {
		s.ErrorWithChannel(s.config().DefaultChannel, zap.String("PrintThingkingData", err.Error()))
		s.TrackWithChannel(THINKINGDATA, zap.ByteString(MsgBody, nil))
	}
}

func (s *StdLogger) PrintBIData(data bi.Data) {
	bytes, err := data.MarshalAsJson()
	if err != nil {
		s.ErrorWithChannel(s.config().DefaultChannel, zap.String("PrintBIData", err.Error()))
	}
	s.TrackWithChannel(BI, zap.ByteString(MsgBody, bytes))
}

func (s *StdLogger) PrintBigQuery(tableName zap.Field, fields ...zap.Field) {
	bytes, err := utils.Zap2Json(fields)
	if err != nil {
		s.ErrorWithChannel(s.config().DefaultChannel, zap.String("PrintBigQuery", err.Error()))
	}
	fields = append([]zap.Field{tableName, zap.ByteString(MsgBody, bytes)})
	s.TrackWithChannel(BIGQUERY, fields...)
}
