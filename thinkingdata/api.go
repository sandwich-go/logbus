package thinkingdata

import "errors"

//https://doc.thinkingdata.cn/tdamanual/installation/pre_installation/data_format.html

func User(accountId, distinctId, dataType string, properties map[string]interface{}) (Data, error) {
	if properties == nil && dataType != USER_DEL {
		return emptyData, errors.New("invalid params for " + dataType + ": properties is nil")
	}
	return add(accountId, distinctId, dataType, "", "", properties)
}

func Track(accountId, distinctId, eventName, eventID string, properties map[string]interface{}) (Data, error) {
	if len(eventName) == 0 {
		return emptyData, errors.New("the event name must be provided")
	}
	dataType := TRACK
	return add(accountId, distinctId, dataType, eventName, eventID, properties)
}

func add(accountId, distinctId, dataType, eventName, eventID string, properties map[string]interface{}) (Data, error) {
	if len(accountId) == 0 && len(distinctId) == 0 {
		return emptyData, errors.New("invalid parameters: account_id and distinct_id cannot be empty at the same time")
	}

	// 获取 properties 中 #ip 值, 如不存在则返回 ""
	ip := extractStringProperty(properties, "#ip")

	// 获取 properties 中 #time 值, 如不存在则返回当前时间
	eventTime := extractTime(properties)

	//如果上传#uuid， 只支持UUID标准格式xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx的string类型
	uuid := extractStringProperty(properties, "#uuid")

	data := Data{
		AccountId:  accountId,
		DistinctId: distinctId,
		Type:       dataType,
		Time:       eventTime,
		EventName:  eventName,
		EventId:    eventID,
		Ip:         ip,
		UUID:       uuid,
		Properties: properties,
	}

	// 检查数据格式, 并将时间类型数据转为符合格式要求的字符串
	return formatProperties(data)
}
