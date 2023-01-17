package conversions

import (
	"fmt"
)

func FloatToString(f float64) string {
	f = FloatToFixed(f, 2)
	return fmt.Sprintf("%.2f", f)
}

func GetStringOk(in interface{}) (val string, ok bool) {
	switch i := in.(type) {
	case bool:
		val = fmt.Sprintf("%v", i)
	case int:
		val = fmt.Sprintf("%v", i)
	case float64:
		val = fmt.Sprintf("%v", i)
	case float32:
		val = fmt.Sprintf("%v", i)
	case int64:
		val = fmt.Sprintf("%v", i)
	case uint64:
		val = fmt.Sprintf("%v", i)
	case string:
		val = i
	case interface{}:
		val = fmt.Sprint(in)
	default:
		return "", false
	}
	return val, true
}
