package thinkingdata

import "errors"

// https://doc.thinkingdata.cn/tdamanual/installation/pre_installation/data_format.html

func User(accountId, distinctId, dataType, appid string, properties map[string]interface{}) (Data, error) {
	// 公开入口的属性未经提取路径过滤，必须保留属性名校验及非法 key 报错。
	return user(accountId, distinctId, dataType, appid, properties, true)
}

func user(accountId, distinctId, dataType, appid string, properties map[string]interface{}, checkKeys bool) (Data, error) {
	if properties == nil && dataType != USER_DEL {
		return emptyData, errors.New("invalid params for " + dataType + ": properties is nil")
	}
	return add(accountId, distinctId, dataType, "", "", appid, properties, checkKeys)
}

func Track(accountId, distinctId, eventName, eventID, appid string, properties map[string]interface{}) (Data, error) {
	return TrackWithType(TRACK, accountId, distinctId, eventName, eventID, appid, properties)
}

func TrackWithType(dataType, accountId, distinctId, eventName, eventID, appid string, properties map[string]interface{}) (Data, error) {
	// Track 也经由此入口；直接传入的 properties 必须完整校验。
	return trackWithType(dataType, accountId, distinctId, eventName, eventID, appid, properties, true)
}

func trackWithType(dataType, accountId, distinctId, eventName, eventID, appid string, properties map[string]interface{}, checkKeys bool) (Data, error) {
	if len(eventName) == 0 {
		return emptyData, errors.New("the event name must be provided")
	}
	return add(accountId, distinctId, dataType, eventName, eventID, appid, properties, checkKeys)
}

// add 提取协议字段并处理属性值。checkKeys 只控制属性名正则校验，
// 仅 ExtractEncoder / extractScalarFields 已校验或过滤顶层 key 后可传 false。
// 账号检查、时间/IP/UUID 提取、事件名和 USER_ADD 校验仍按原顺序执行。
func add(accountId, distinctId, dataType, eventName, eventID, appid string, properties map[string]interface{}, checkKeys bool) (Data, error) {
	if len(accountId) == 0 && len(distinctId) == 0 {
		return emptyData, errors.New("invalid parameters: account_id and distinct_id cannot be empty at the same time")
	}

	// 获取 properties 中 #ip 值, 如不存在则返回 ""
	ip := extractStringProperty(properties, "#ip")

	// 获取 properties 中 #time 值, 如不存在则返回当前时间
	eventTime := extractTime(properties)

	// 如果上传#uuid， 只支持UUID标准格式xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx的string类型
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
		Appid:      appid,
		Properties: properties,
	}

	// 检查数据格式, 并将时间类型数据转为符合格式要求的字符串
	return formatProperties(data, checkKeys)
}
