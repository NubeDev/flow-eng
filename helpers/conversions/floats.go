package conversions

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	"math"
	"strconv"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func FloatToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

type number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

func BoolToNum(x bool) float64 {
	if x {
		return 0
	}
	return 1
}

func IsBool(value interface{}) bool {
	switch value.(type) {
	case bool:
		return true
	default:
		return false
	}
}

func IsBoolWithValue(value interface{}) (isBool, val bool) {
	switch v := value.(type) {
	case bool:
		return true, v
	default:
		return false, false
	}
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
	case uint64:
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
	case uint64:
		val = float.New(float64(i))
	default:
		return nil, false
	}
	return val, true
}

func GetFloat(in interface{}) (val float64) {
	switch i := in.(type) {
	case bool:
		if i {
			val = 1
		} else {
			val = 0
		}
	case int:
		val = float64(i)
	case float64:
		val = i
	case float32:
		val = float64(i)
	case int64:
		val = float64(i)
	case uint64:
		val = float64(i)
	case string:
		if s, err := strconv.ParseFloat(fmt.Sprintf("%v", in), 64); err == nil {
			val = s
		} else {
			val = 0
		}
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
	case uint64:
		val = float64(i)
	case string:
		if s, err := strconv.ParseFloat(fmt.Sprintf("%v", in), 64); err == nil {
			val = s
		} else {
			val = 0
			ok = false
		}
	default:
		return 0, false
	}
	return val, true
}

func GetInt(in interface{}) (val int) {
	switch i := in.(type) {
	case bool:
		if i {
			val = 1
		} else {
			val = 0
		}
	case int:
		val = i
	case float64:
		val = int(i)
	case float32:
		val = int(i)
	case int64:
		val = int(i)
	case uint64:
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
	case uint64:
		val = int(i)
	default:
		return 0, false
	}
	return val, true
}
