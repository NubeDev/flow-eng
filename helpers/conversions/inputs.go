package conversions

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
)

type number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

func BoolToNum(x bool) float64 {
	if x {
		return 0
	}
	return 1
}

func NumToBool[T number](x T) bool {
	if x >= 1 {
		return true
	}
	return false
}

func ToString(in interface{}) string {
	return fmt.Sprintf("%v", in)
}

func GetFloatPointer(in interface{}) (val *float64) {
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
		return nil
	}
	return val
}

func GetFloatPointerOk(in interface{}) (val *float64, ok bool) {
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

func GetFloat(in interface{}) (val float64) {
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
		return 0
	}
	return val
}

func GetFloatOk(in interface{}) (val float64, ok bool) {
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

func GetInt(in interface{}) (val int) {
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
		return 0
	}
	return val
}

func GetIntOk(in interface{}) (val int, ok bool) {
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
