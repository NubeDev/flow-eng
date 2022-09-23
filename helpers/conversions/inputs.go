package conversions

import "github.com/NubeDev/flow-eng/helpers/float"

func GetFloatPointer(in interface{}) (val *float64, ok bool) {
	switch i := in.(type) {
	case int:
		val = float.New(float64(i))
	case float64:
		val = float.New(i)
	case float32:
		val = float.New(float64(i))
	case int64:
		val = float.New(float64(i))
	default:
		return nil, false
	}
	return val, true
}

func GetFloat(in interface{}) (val float64, ok bool) {
	switch i := in.(type) {
	case int:
		val = float64(i)
	case float64:
		val = i
	case float32:
		val = float64(i)
	case int64:
		val = float64(i)
	default:
		return 0, false
	}
	return val, true
}

func GetInt(in interface{}) (val int, ok bool) {
	switch i := in.(type) {
	case int:
		val = i
	case float64:
		val = int(i)
	case float32:
		val = int(i)
	case int64:
		val = int(i)
	default:
		return 0, false
	}
	return val, true
}
