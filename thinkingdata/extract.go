package thinkingdata

import (
	"errors"
	"unicode"

	"go.uber.org/zap/zapcore"
)

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
	userType, ok1 := memoryEncoder.Fields[TYPE]
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

	delete(memoryEncoder.Fields, ACCOUNT)
	delete(memoryEncoder.Fields, DISTINCT)
	delete(memoryEncoder.Fields, TYPE)
	delete(memoryEncoder.Fields, EVENT)
	delete(memoryEncoder.Fields, EVENT_ID)
	if hasEvent {
		return Track(accountId.(string), distinctId.(string), eventName.(string), strEventID, memoryEncoder.Fields)
	}
	if ok1 {
		if userType.(string) == TRACK {
			return emptyData, errors.New("the event name must be provided")
		}
		return User(accountId.(string), distinctId.(string), userType.(string), memoryEncoder.Fields)
	}
	return emptyData, errors.New("no #type or #event_name")
}
