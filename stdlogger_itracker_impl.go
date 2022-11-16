package logbus

import (
	"github.com/sandwich-go/logbus/utils"

	"github.com/sandwich-go/logbus/thinkingdata"

	"go.uber.org/zap"
)

/*func (s *StdLogger) Track(fields ...zap.Field) error {
	for _, tag := range s.tags {
		switch tag {
		case THINKINGDATA:
			memoryEncoder := zapcore.NewMapObjectEncoder()
			for _, v := range fields {
				v.AddTo(memoryEncoder)
			}
			data, err := thinkingdata.ExtractEncoder(memoryEncoder)
			if err != nil {
				return err
			}
			s.PrintThingkingData(data)
		case BIGQUERY:
			tableName, bigFields, err := bigquery.ExtractEncoder(fields)
			if err != nil {
				return err
			}
			s.PrintBigQuery(tableName, bigFields...)
		default:
			return fmt.Errorf("tag %s not implement", tag)
		}
	}
	return nil
}*/

func (s *StdLogger) PrintThingkingData(data thinkingdata.Data) {
	bytes, err := data.MarshalAsJson()
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
