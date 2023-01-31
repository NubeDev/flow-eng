package conversions

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/str"
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

func GetStringPointerOk(in interface{}) (val *string, ok bool) {
	switch i := in.(type) {
	case bool:
		val = str.New(fmt.Sprintf("%v", i))
	case int:
		val = str.New(fmt.Sprintf("%v", i))
	case float64:
		val = str.New(fmt.Sprintf("%v", i))
	case float32:
		val = str.New(fmt.Sprintf("%v", i))
	case int64:
		val = str.New(fmt.Sprintf("%v", i))
	case uint64:
		val = str.New(fmt.Sprintf("%v", i))
	case string:
		val = str.New(i)
	case interface{}:
		val = str.New(fmt.Sprint(in))
	default:
		return str.New(""), false
	}
	return val, true
}

func ConvertInterfaceToStringMultiple(values []interface{}) []*string {
	var output []*string
	for _, value := range values {
		v, ok := GetStringPointerOk(value)
		if ok {
			output = append(output, v)
		}
	}
	return output
}
