package boolean

import (
	"fmt"
	"strconv"
)

func New(value bool) *bool {
	b := value
	return &b
}

func NewTrue() *bool {
	b := true
	return &b
}

func NewFalse() *bool {
	b := false
	return &b
}

func IsTrue(b *bool) bool {
	if b == nil {
		return false
	} else {
		return *b
	}
}

func IsFalse(b *bool) bool {
	return !IsTrue(b)
}

func NonNil(b *bool) bool {
	if b == nil {
		return false
	} else {
		if *b == true {
			return true
		} else {
			return false
		}

	}
}

func IsNil(b *bool) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

func ConvertInterfaceToBool(value interface{}) *bool {
	var valAsBoolPtr *bool
	if value != nil {
		valAsBool, _ := strconv.ParseBool(fmt.Sprintf("%v", value))
		valAsBoolPtr = New(valAsBool)
	}
	return valAsBoolPtr
}

func ConvertInterfaceToBoolMultiple(values []interface{}) []*bool {
	var output []*bool
	for _, value := range values {
		output = append(output, ConvertInterfaceToBool(value))
	}
	return output
}
