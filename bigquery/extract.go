package bigquery

import (
	"errors"
	"regexp"
	"strings"

	"bitbucket.org/funplus/sandwich/pkg/logbus/utils"

	"go.uber.org/zap"
)

// $abc abc abc_1 abc1
var (
	KeyPattern       = "^[a-zA-Z" + ColumnPrefix + "][A-Za-z0-9_]{0,49}$"
	ColumnPrefix     = "$"
	ColumnProperties = "data"
	TableNameKey     = ColumnPrefix + "tablename"
	UseRecord        = false // https://cloud.google.com/bigquery/docs/nested-repeated?hl=zh-cn
)

var KeyPatternReg, _ = regexp.Compile(KeyPattern)

func ExtractEncoder(fields []zap.Field) (tableName zap.Field, res []zap.Field, err error) {
	var data []zap.Field
	for _, v := range fields {
		if v.Key == TableNameKey {
			tableName = v
			continue
		}
		if !KeyPatternReg.MatchString(v.Key) {
			continue
		}
		if strings.HasPrefix(v.Key, ColumnPrefix) {
			v.Key = strings.Replace(v.Key, ColumnPrefix, "", 1)
			res = append(res, v)
			continue
		}
		data = append(data, v)
	}
	if tableName.String == "" {
		err = errors.New("empty table name")
		return
	}
	if len(data) > 0 {
		if UseRecord {
			res = append(res, zap.Namespace(ColumnProperties))
			res = append(res, data...)
		} else {
			bytes, err1 := utils.Zap2Json(data)
			if err1 != nil {
				err = err1
				return
			}
			res = append(res, zap.ByteString(ColumnProperties, bytes))
		}
	}
	return
}
