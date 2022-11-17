package stdl

/*func (s *StdLogger) Track(fields ...zap.Field) error {
	for _, v := range s.tags {
		switch v {
		case config.THINKINGDATA:
			memoryEncoder := zapcore.NewMapObjectEncoder()
			for _, v := range fields {
				v.AddTo(memoryEncoder)
			}
			data, err := thinkingdata.ExtractEncoder(memoryEncoder)
			if err != nil {
				return err
			}
			s.PrintThingkingData(data)
		case config.BIGQUERY:
			tableName, bigFields, err := bigquery.ExtractEncoder(fields)
			if err != nil {
				return err
			}
			s.PrintBigQuery(tableName, bigFields...)
		default:
			return fmt.Errorf("tag %s not implement", v)
		}
	}
	return nil
}

func (s *StdLogger) PrintThingkingData(data thinkingdata.Data) {
	bytes, err := data.MarshalAsJson()
	if err != nil {
		s.ErrorWithChannel(basics.Setting.DefaultChannel, zap.String("PrintThingkingData", err.Error()))
	}
	s.InfoWithChannel(config.THINKINGDATA, zap.ByteString(config.MsgBody, bytes))
}

func (s *StdLogger) PrintBigQuery(tableName zap.Field, fields ...zap.Field) {
	bytes, err := utils.Zap2Json(fields)
	if err != nil {
		s.ErrorWithChannel(basics.Setting.DefaultChannel, zap.String("PrintBigQuery", err.Error()))
	}
	fields = append([]zap.Field{tableName, zap.ByteString(config.MsgBody, bytes)})
	s.InfoWithChannel(config.BIGQUERY, fields...)
}*/
