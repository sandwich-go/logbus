package thinkingdata

// from https://github.com/ThinkingDataAnalytics/go-sdk/blob/master/thinkingdata/utils.go

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"time"
)

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
			return time.Now().Format(DATE_FORMAT)
		}
	}

	return time.Now().Format(DATE_FORMAT)
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

func formatProperties(d Data) (Data, error) {
	if d.EventName != "" {
		matched := checkPattern([]byte(d.EventName))
		if !matched {
			return emptyData, errors.New("Invalid event name: " + d.EventName)
		}
	}

	if d.Properties != nil {
		for k, v := range d.Properties {
			isMatch := checkPattern([]byte(k))
			if !isMatch {
				return emptyData, errors.New("Invalid property key: " + k)
			}

			if d.Type == USER_ADD && isNotNumber(v) {
				return emptyData, errors.New("Invalid property value: only numbers is supported by UserAdd")
			}

			//check value
			switch v.(type) {
			case bool:
			case string:
			case []string, []interface{}:
			case time.Time: //only support time.Time
				d.Properties[k] = v.(time.Time).Format(DATE_FORMAT)
			default:
				if isNotNumber(v) {
					return emptyData, fmt.Errorf("Invalid property value %v. Supported types: numbers, string, time.Time, bool, []string", reflect.TypeOf(v))
				}
			}
		}
	}

	return d, nil
}
