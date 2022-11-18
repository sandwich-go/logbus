package logbus

import "go.uber.org/zap"

// NameSpace 记录层级关系
func NameSpace(name string) Field {
	return zap.Namespace(name)
}
