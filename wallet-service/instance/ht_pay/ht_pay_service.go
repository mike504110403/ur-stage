package htpay

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// 轉請求字串符
func ToQueryString(req interface{}) (string, error) {
	val := reflect.ValueOf(req)
	typ := reflect.TypeOf(req)

	var kvPairs []string

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("json")

		// 獲取字段值
		value := val.Field(i).Interface()
		var valueStr string

		// 空或零值檢查
		if val.Field(i).IsZero() {
			continue
		}

		// 根據字段類型進行處理
		switch v := value.(type) {
		case string:
			valueStr = v
		case float64:
			valueStr = fmt.Sprintf("%.2f", v)
		case int:
			valueStr = fmt.Sprintf("%d", v) // 將 int 轉換為字符串
		default:
			return "", fmt.Errorf("unsupported field type: %T", v)
		}

		kvPairs = append(kvPairs, fmt.Sprintf("%s=%s", tag, valueStr))
	}
	// 排序
	sort.Strings(kvPairs)

	// 組字符串
	return strings.Join(kvPairs, "&"), nil
}
