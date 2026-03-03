package bi

import (
	"errors"
	"time"
)

// Track 创建 BI 打点事件数据
// 参数均为必传：appID 项目 BI app ID, roleID 角色 ID, event 事件名, fpid
// properties 为可自定义属性，可为 nil 表示无额外属性
// https://centurygames.feishu.cn/wiki/DKNHwwusVi6q2tkeGKncaN6Tnjb
func Track(appID, roleID, event, fpid string, properties map[string]interface{}) (Data, error) {
	if event == "" {
		return emptyData, errors.New("event is required")
	}

	if properties == nil {
		properties = make(map[string]interface{})
	}

	return Data{
		AppID:      appID,
		Ts:         time.Now().UnixMilli(),
		RoleID:     roleID,
		Event:      event,
		FpID:       fpid,
		Properties: properties,
	}, nil
}
