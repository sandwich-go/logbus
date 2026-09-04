package thinkingdata

import (
	"errors"
	"math"
	"time"
	"unicode"

	"go.uber.org/zap/zapcore"
)

// ExtractFields converts ordinary scalar zap fields without allocating an
// intermediate MapObjectEncoder. Fields that need zap's marshaling behavior
// continue through ExtractEncoder so their existing semantics are preserved.
func ExtractFields(fields []zapcore.Field) (Data, error) {
	data, err, ok := extractScalarFields(fields)
	if ok {
		return data, err
	}

	memoryEncoder := zapcore.NewMapObjectEncoder()
	for _, field := range fields {
		field.AddTo(memoryEncoder)
	}
	return ExtractEncoder(memoryEncoder)
}

func ExtractEncoder(memoryEncoder *zapcore.MapObjectEncoder) (Data, error) {
	for k := range memoryEncoder.Fields {
		if !KeyPattern.MatchString(k) {
			delete(memoryEncoder.Fields, k)
		}
	}
	accountId, ok1 := memoryEncoder.Fields[ACCOUNT]
	distinctId, ok2 := memoryEncoder.Fields[DISTINCT]
	if !ok1 {
		accountId = ""
	}
	if !ok2 {
		distinctId = ""
	}
	if !ok1 && !ok2 {
		return emptyData, errors.New("#account_id and #distinct_id not exist")
	}
	dataType, ok1 := memoryEncoder.Fields[TYPE]
	eventName, hasEvent := memoryEncoder.Fields[EVENT]

	// event_id
	eventID, _ := memoryEncoder.Fields[EVENT_ID]
	if eventID == nil {
		eventID = ""
	}
	strEventID, ok := eventID.(string)
	if !ok {
		return emptyData, errors.New("#event_id must be string")
	}
	if len(strEventID) > 0 { // #event_id如果存在，必须以字母或数字开头
		firstChar := rune(eventID.(string)[0])
		if !unicode.IsLetter(firstChar) && !unicode.IsDigit(firstChar) {
			return emptyData, errors.New("the event name must start with a letter or number")
		}
	}
	appid, ok3 := memoryEncoder.Fields[APPID]
	if !ok3 {
		appid = ""
	}
	delete(memoryEncoder.Fields, ACCOUNT)
	delete(memoryEncoder.Fields, DISTINCT)
	delete(memoryEncoder.Fields, TYPE)
	delete(memoryEncoder.Fields, EVENT)
	delete(memoryEncoder.Fields, EVENT_ID)
	delete(memoryEncoder.Fields, APPID)
	if hasEvent {
		if !ok1 {
			dataType = TRACK // 没传TYPE默认TRACK
		}
		return TrackWithType(dataType.(string), accountId.(string), distinctId.(string), eventName.(string), strEventID, appid.(string), memoryEncoder.Fields)
	}
	if ok1 {
		if dataType.(string) == TRACK {
			return emptyData, errors.New("the event name must be provided")
		}
		return User(accountId.(string), distinctId.(string), dataType.(string), appid.(string), memoryEncoder.Fields)
	}
	return emptyData, errors.New("no #type or #event_name")
}

func extractScalarFields(fields []zapcore.Field) (Data, error, bool) {
	properties := make(map[string]interface{}, len(fields))
	var accountID, distinctID, dataType, eventName, eventID, appid string
	var hasAccountID, hasDistinctID, hasDataType, hasEventName, hasEventID bool

	for _, field := range fields {
		if field.Type == zapcore.SkipType {
			continue
		}
		if !KeyPattern.MatchString(field.Key) {
			return emptyData, nil, false
		}

		switch field.Key {
		case ACCOUNT:
			value, ok := stringFieldValue(field, false)
			if !ok {
				return emptyData, nil, false
			}
			accountID, hasAccountID = value, true
		case DISTINCT:
			value, ok := stringFieldValue(field, false)
			if !ok {
				return emptyData, nil, false
			}
			distinctID, hasDistinctID = value, true
		case TYPE:
			value, ok := stringFieldValue(field, false)
			if !ok {
				return emptyData, nil, false
			}
			dataType, hasDataType = value, true
		case EVENT:
			value, ok := stringFieldValue(field, false)
			if !ok {
				return emptyData, nil, false
			}
			eventName, hasEventName = value, true
		case EVENT_ID:
			value, ok := stringFieldValue(field, true)
			if !ok {
				return emptyData, nil, false
			}
			eventID, hasEventID = value, true
		case APPID:
			value, ok := stringFieldValue(field, false)
			if !ok {
				return emptyData, nil, false
			}
			appid = value
		default:
			value, skip, ok := fieldValue(field)
			if !ok {
				return emptyData, nil, false
			}
			if skip {
				continue
			}
			properties[field.Key] = value
		}
	}

	if !hasAccountID && !hasDistinctID {
		return emptyData, errors.New("#account_id and #distinct_id not exist"), true
	}

	if hasEventID && len(eventID) > 0 {
		firstChar := rune(eventID[0])
		if !unicode.IsLetter(firstChar) && !unicode.IsDigit(firstChar) {
			return emptyData, errors.New("the event name must start with a letter or number"), true
		}
	}

	if hasEventName {
		if !hasDataType {
			dataType = TRACK
		}
		data, err := TrackWithType(dataType, accountID, distinctID, eventName, eventID, appid, properties)
		return data, err, true
	}
	if hasDataType {
		if dataType == TRACK {
			return emptyData, errors.New("the event name must be provided"), true
		}
		data, err := User(accountID, distinctID, dataType, appid, properties)
		return data, err, true
	}
	return emptyData, errors.New("no #type or #event_name"), true
}

func stringFieldValue(field zapcore.Field, allowNil bool) (string, bool) {
	switch field.Type {
	case zapcore.StringType:
		return field.String, true
	case zapcore.ByteStringType:
		return string(field.Interface.([]byte)), true
	case zapcore.ReflectType:
		if field.Interface == nil && allowNil {
			return "", true
		}
		value, ok := field.Interface.(string)
		return value, ok
	default:
		return "", false
	}
}

func fieldValue(field zapcore.Field) (value interface{}, skip, ok bool) {
	switch field.Type {
	case zapcore.BinaryType:
		return field.Interface.([]byte), false, true
	case zapcore.BoolType:
		return field.Integer == 1, false, true
	case zapcore.ByteStringType:
		return string(field.Interface.([]byte)), false, true
	case zapcore.Complex128Type:
		return field.Interface.(complex128), false, true
	case zapcore.Complex64Type:
		return field.Interface.(complex64), false, true
	case zapcore.DurationType:
		return time.Duration(field.Integer), false, true
	case zapcore.Float64Type:
		return math.Float64frombits(uint64(field.Integer)), false, true
	case zapcore.Float32Type:
		return math.Float32frombits(uint32(field.Integer)), false, true
	case zapcore.Int64Type:
		return field.Integer, false, true
	case zapcore.Int32Type:
		return int32(field.Integer), false, true
	case zapcore.Int16Type:
		return int16(field.Integer), false, true
	case zapcore.Int8Type:
		return int8(field.Integer), false, true
	case zapcore.StringType:
		return field.String, false, true
	case zapcore.TimeType:
		if field.Interface != nil {
			return time.Unix(0, field.Integer).In(field.Interface.(*time.Location)), false, true
		}
		return time.Unix(0, field.Integer), false, true
	case zapcore.TimeFullType:
		return field.Interface.(time.Time), false, true
	case zapcore.Uint64Type:
		return uint64(field.Integer), false, true
	case zapcore.Uint32Type:
		return uint32(field.Integer), false, true
	case zapcore.Uint16Type:
		return uint16(field.Integer), false, true
	case zapcore.Uint8Type:
		return uint8(field.Integer), false, true
	case zapcore.UintptrType:
		return uintptr(field.Integer), false, true
	case zapcore.ReflectType:
		return field.Interface, false, true
	case zapcore.SkipType:
		return nil, true, true
	default:
		return nil, false, false
	}
}
