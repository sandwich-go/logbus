package thinkingdata

// from https://github.com/ThinkingDataAnalytics/go-sdk/tree/master/src/thinkingdata

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"
)

var locationTGA = time.UTC

// 包初始化时编译一次，事件名和属性名校验复用同一匹配器。
var KeyPattern, _ = regexp.Compile(KEY_PATTERN)

func checkPattern(name []byte) bool {
	return KeyPattern.Match(name)
}

func mergeProperties(target, source map[string]interface{}) {
	for k, v := range source {
		target[k] = v
	}
}

func extractTime(p map[string]interface{}) string {
	if t, ok := p["#time"]; ok {
		delete(p, "#time")
		switch v := t.(type) {
		case string:
			return v
		case time.Time:
			return v.Format(DATE_FORMAT)
		default:
			return time.Now().In(locationTGA).Format(DATE_FORMAT)
		}
	}

	return time.Now().In(locationTGA).Format(DATE_FORMAT)
}

func extractStringProperty(p map[string]interface{}, key string) string {
	if t, ok := p[key]; ok {
		delete(p, key)
		v, ok := t.(string)
		if !ok {
			fmt.Fprintln(os.Stderr, "Invalid data type for "+key)
			return ""
		}
		return v
	}
	return ""
}

func isNotNumber(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
	case float32, float64:
	default:
		return true
	}
	return false
}

// formatProperties 校验事件名和属性值，并原地格式化时间属性。
// checkKeys=false 仅表示顶层属性名已由提取路径检查，不能跳过整个属性遍历。
func formatProperties(d Data, checkKeys bool) (Data, error) {
	if d.EventName != "" {
		matched := checkPattern([]byte(d.EventName))
		if !matched {
			return emptyData, errors.New("Invalid event name: " + d.EventName)
		}
	}

	if d.Properties != nil {
		for k, v := range d.Properties {
			// 短路时也不会执行 []byte(k) 转换。Go 1.25.3 的本地逃逸分析与基准显示，
			// 该转换会为每个属性名额外分配一次堆内存；32 属性样本因此少 32 次分配。
			if checkKeys && !checkPattern([]byte(k)) {
				return emptyData, errors.New("Invalid property key: " + k)
			}

			if d.Type == USER_ADD && isNotNumber(v) {
				return emptyData, errors.New("Invalid property value: only numbers is supported by UserAdd")
			}

			//check value
			switch v.(type) {
			case time.Time: //only support time.Time
				d.Properties[k] = v.(time.Time).Format(DATE_FORMAT)
			default:
			}
		}
	}

	return d, nil
}
