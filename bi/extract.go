package bi

import (
	"go.uber.org/zap/zapcore"
)

const (
	KeyAppID   = "app_id"
	KeyRoleID  = "role_id"
	KeyEvent   = "event"
	KeyFpid    = "fpid"
)

// ExtractEncoder 从 zap MapObjectEncoder 中提取 BI 打点数据
// 需包含 app_id、role_id、event、fpid 必传字段，其余字段作为 properties
func ExtractEncoder(memoryEncoder *zapcore.MapObjectEncoder) (Data, error) {
	fields := memoryEncoder.Fields

	appID := extractString(fields, KeyAppID)
	roleID := extractString(fields, KeyRoleID)
	event := extractString(fields, KeyEvent)
	fpid := extractString(fields, KeyFpid)

	delete(fields, KeyAppID)
	delete(fields, KeyRoleID)
	delete(fields, KeyEvent)
	delete(fields, KeyFpid)

	properties := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		properties[k] = v
	}

	return Track(appID, roleID, event, fpid, properties)
}

func extractString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
