package float

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/str"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func StringFloatErr(value *string) (*float64, float64, error) {
	if str.IsNil(value) {
		return nil, 0, nil
	}
	v, err := StrToFloatErr(str.NonNil(value))
	if err != nil {
		return nil, 0, err
	} else {
		return New(v), v, nil
	}
}

func StringFloat(value *string) (*float64, float64) {
	v, nonNil, _ := StringFloatErr(value)
	return v, nonNil
}

func StrToFloat(value string) float64 {
	v, err := StrToFloatErr(value)
	if err != nil {
		return 0
	} else {
		return v
	}
}

func PtrToStringPtr(value *float64) *string {
	v := PointerToStr(value)
	return str.New(v)
}

func PointerToStr(value *float64) string {
	return fmt.Sprintf("%f", NonNil(value))
}

func ToStrPtr(value float64) *string {
	v := fmt.Sprintf("%f", value)
	return str.New(v)
}

func ToStr(value float64) string {
	return fmt.Sprintf("%f", value)
}

func StrToFloatErr(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

func FirstNotNil(values ...*float64) *float64 {
	for _, n := range values {
		if n != nil {
			return n
		}
	}
	return nil
}

func ComparePtrValues(value1 *float64, value2 *float64) bool {
	return NonNilMax(value1) == NonNilMax(value2)
}

func NonNilMax(value *float64) float64 {
	if value == nil {
		return math.MaxFloat64
	}
	return *value
}

func New(value float64) *float64 {
	return &value
}

func NonNil(b *float64) float64 {
	if b == nil {
		return 0
	}
	return *b
}

func NonNil32(b *float32) float32 {
	if b == nil {
		return 0
	}
	return *b
}

func NotNil(b *float64) bool {
	if b == nil {
		return false
	} else {
		return true
	}
}

func IsNil(b *float64) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

func IsNil32(b *float32) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

func EvalAsBool(b *float64) *float64 {
	if b == nil {
		return New(0)
	} else if *b == 0 {
		return New(0)
	} else {
		return New(1)
	}
}

func EvalAsBoolOnlyOneIsTrue(b *float64) *float64 {
	if b == nil {
		return New(0)
	} else if *b == 1 {
		return New(1)
	} else {
		return New(0)
	}
}

func Copy(b *float64) *float64 {
	if b == nil {
		return nil
	}
	out := *b
	return &out
}

// RandFloat returns a random float64 within the specified range.
func RandFloat(range1, range2 float64) float64 {
	if range1 == range2 {
		return range1
	}
	var min, max float64
	if range1 > range2 {
		max = range1
		min = range2
	} else {
		max = range2
		min = range1
	}
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

func ConvertInterfaceToFloat(value interface{}) *float64 {
	if value == nil {
		return nil
	}
	output, ok := value.(float64)
	if ok {
		return &output
	}
	return nil
}

func ConvertInterfaceToFloatMultiple(values []interface{}) []*float64 {
	var output []*float64
	for _, value := range values {
		output = append(output, ConvertInterfaceToFloat(value))
	}
	return output
}
