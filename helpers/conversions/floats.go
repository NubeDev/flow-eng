package conversions

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/nmath"
	"math"
	"strconv"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func FloatToFixed(num float64, precision int) float64 {
	if precision < 0 {
		precision = 0
	}
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func TruncateFloat(num float64, precision int) float64 {
	if precision < 0 {
		precision = 0
	}
	return math.Trunc(num*math.Pow(10, float64(precision))) / math.Pow(10, float64(precision))
}

type number interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

func BoolToNum(x bool) float64 {
	if x {
		return 1
	}
	return 0
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
	case bool:
		val = float.New(BoolToNum(i))
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
	case bool:
		val = float.New(BoolToNum(i))
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
	case interface{}:
		s := fmt.Sprint(in)
		if s, err := strconv.ParseFloat(s, 64); err == nil {
			val = float.New(s)
		}
	default:
		return nil, false
	}
	return val, true
}

func GetFloat(in interface{}) (val float64) {
	switch i := in.(type) {
	case bool:
		val = BoolToNum(i)
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
	case interface{}:
		s := fmt.Sprint(in)
		if s, err := strconv.ParseFloat(s, 64); err == nil {
			val = s
		}
	default:
		return 0
	}
	return val
}

func GetFloatOk(in interface{}) (val float64, ok bool) {
	switch i := in.(type) {
	case bool:
		val = BoolToNum(i)
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
	case interface{}:
		s := fmt.Sprint(in)
		if s, err := strconv.ParseFloat(s, 64); err == nil {
			val = s
		}
	default:
		return 0, false
	}
	return val, true
}

func GetInt(in interface{}) (val int) {
	switch i := in.(type) {
	case bool:
		val = int(BoolToNum(i))
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
	case bool:
		val = int(BoolToNum(i))
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
	case interface{}:
		s := fmt.Sprint(in)
		if i, err := strconv.Atoi(s); err == nil {
			val = i
		}
	default:
		return 0, false
	}
	return val, true
}

func ConvertInterfaceToFloatMultiple(values []interface{}) []*float64 {
	var output []*float64
	for _, value := range values {
		v, ok := GetFloatPointerOk(value)
		if ok {
			output = append(output, v)
		}
	}
	return output
}

func PointScale(presentValue, scaleInMin, scaleInMax, scaleOutMin, scaleOutMax float64) (value float64) {
	if scaleInMin == 0 && scaleInMax == 0 && scaleOutMin == 0 && scaleOutMax == 0 {
		return presentValue
	}
	out := nmath.Scale(presentValue, scaleInMin, scaleInMax, scaleOutMin, scaleOutMax)
	return out
}

func PointRange(presentValue, limitMin, limitMax float64) (value float64) {
	if limitMin == 0 && limitMax == 0 {
		return presentValue
	}
	out := nmath.LimitToRange(presentValue, limitMin, limitMax)
	return out
}

func ValueTransformOnRead(originalValue float64, scaleEnable bool, factor, scaleInMin, scaleInMax, scaleOutMin,
	scaleOutMax, offset float64) (transformedValue float64) {
	ov := originalValue

	// perform factor operation
	factored := ov
	if factor != 0 {
		factored = ov * factor
	}

	// perform scaling and limit operations
	scaledAndLimited := factored
	if scaleEnable {
		if scaleOutMin != 0 || scaleOutMax != 0 {
			if scaleInMin != 0 || scaleInMax != 0 { // scale with all 4 configs
				scaledAndLimited = PointScale(factored, scaleInMin, scaleInMax, scaleOutMin,
					scaleOutMax)
			} else { // do limit with only scaleOutMin and scaleOutMin
				scaledAndLimited = PointRange(factored, scaleOutMin, scaleOutMax)
			}
		}
	}
	// perform offset operation
	offsetted := scaledAndLimited + offset
	return offsetted
}

func ValueTransformOnWrite(originalValue float64, scaleEnable bool, factor, scaleInMin, scaleInMax, scaleOutMin,
	scaleOutMax, offset float64) (transformedValue float64) {
	ov := originalValue

	// reverse offset operation
	unoffsetted := ov - offset

	// reverse scaling and limit operations
	unscaledAndUnlimited := unoffsetted
	if scaleEnable {
		if scaleOutMin != 0 || scaleOutMax != 0 {
			if scaleInMin != 0 || scaleInMax != 0 { // scale with all 4 configs
				unscaledAndUnlimited = PointScale(unoffsetted, scaleOutMin, scaleOutMax,
					scaleInMin, scaleInMax)
			} else { // do limit with only scaleOutMin and scaleOutMin
				unscaledAndUnlimited = PointRange(unoffsetted, scaleOutMin, scaleOutMax)
			}
		}
	}

	// reverse factoring operation
	unfactored := unscaledAndUnlimited
	if factor != 0 {
		unfactored = unscaledAndUnlimited / factor
	}

	return unfactored
}
