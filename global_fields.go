package logbus

import (
	"os"

	"github.com/rs/xid"
	"github.com/sandwich-go/boost/xip"
	"github.com/sandwich-go/boost/xtime"
	"go.uber.org/zap"
)

var globalFields []zap.Field
var cacheUserDefineFields []zap.Field

// DynamicGlobalField 动态全局字段类型
type DynamicGlobalField func() Field

var dynamicGlobalFields []DynamicGlobalField

func init() {
	for _, key := range []string{"sys_env_name", "sys_stage", "sys_cd_service"} {
		val := os.Getenv(key)
		if val != "" {
			ReservedGlobalFields = append(ReservedGlobalFields, String(key, val))
		}
	}
	ReservedGlobalFields = append(ReservedGlobalFields, String("server_id", xid.New().String()))
	ReservedGlobalFields = append(ReservedGlobalFields, String("server_ip", xip.GetLocalIP()))
	ReservedGlobalFields = append(ReservedGlobalFields, Int64("server_birth", xtime.Unix()))
	if hostName, err := os.Hostname(); err == nil {
		ReservedGlobalFields = append(ReservedGlobalFields, String("host_name", hostName))
	} else {
		ReservedGlobalFields = append(ReservedGlobalFields, String("host_name", "-"))
	}
	AppendGlobalFields()
}

// ReservedGlobalFields 预留的全局字段，可以通过显式这只为空清除
var ReservedGlobalFields []Field

func GetGlobalFields() []Field { return globalFields }

func SetGlobalFields(fields []Field) {
	cacheUserDefineFields = fields
	freshGlobal()
}

func AppendGlobalFields(fields ...Field) {
	cacheUserDefineFields = append(cacheUserDefineFields, fields...)
	freshGlobal()
}

// SetDynamicGlobalFields 设置动态全局字段列表，会替换现有的动态字段
// 注意：非线程安全，必须在初始化时调用, 考虑性能问题，不推荐生产环境使用
func SetDynamicGlobalFields(fields ...DynamicGlobalField) {
	dynamicGlobalFields = fields
}

// AppendDynamicGlobalFields 追加动态全局字段
// 注意：非线程安全，必须在初始化时调用， 考虑性能问题，不推荐生产环境使用
func AppendDynamicGlobalFields(fields ...DynamicGlobalField) {
	dynamicGlobalFields = append(dynamicGlobalFields, fields...)
}

// GetDynamicGlobalFields 获取动态全局字段的值（调用函数获取当前值） unsafely
func GetDynamicGlobalFields() []Field {
	if len(dynamicGlobalFields) == 0 {
		return nil
	}
	var fields []Field
	for _, fn := range dynamicGlobalFields {
		if fn != nil {
			if field := fn(); field.Type != zap.Skip().Type {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

func freshGlobal() {
	globalFields = nil
	for _, v := range ReservedGlobalFields {
		overrideByUser := false
		for _, vt := range cacheUserDefineFields {
			if v.Key == vt.Key {
				overrideByUser = true
			}
		}
		if !overrideByUser {
			// 用户层没有覆盖的字段则使用默认字段
			globalFields = append(globalFields, v)
		}
	}
	globalFields = append(globalFields, cacheUserDefineFields...)
}
