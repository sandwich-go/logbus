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

func init() {
	ReservedGlobalFields = append(ReservedGlobalFields, zap.String("server_id", xid.New().String()))
	ReservedGlobalFields = append(ReservedGlobalFields, zap.String("server_ip", xip.GetLocalIP()))
	ReservedGlobalFields = append(ReservedGlobalFields, zap.Int64("server_birth", xtime.Unix()))
	if hostName, err := os.Hostname(); err == nil {
		ReservedGlobalFields = append(ReservedGlobalFields, zap.String("host_name", hostName))
	} else {
		ReservedGlobalFields = append(ReservedGlobalFields, zap.String("host_name", "-"))
	}
	appendGlobalFields()
}

// ReservedGlobalFields 预留的全局字段，可以通过显式这只为空清除
var ReservedGlobalFields []zap.Field

func getGlobalFields() []zap.Field { return globalFields }

func setGlobalFields(fields []zap.Field) {
	cacheUserDefineFields = fields
	freshGlobal()
}

func appendGlobalFields(fields ...zap.Field) {
	cacheUserDefineFields = append(cacheUserDefineFields, fields...)
	freshGlobal()
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
