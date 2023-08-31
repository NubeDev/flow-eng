package str

import (
	"fmt"
)

func New(value string) *string {
	b := value
	return &b
}
func NonNil(b *string) string {
	if b == nil {
		return ""
	}
	return *b
}

func IsNil(b *string) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

func ConvertInterfaceToString(value interface{}) *string {
	var valAsBoolPtr *string
	if value != nil {
		valAsBoolPtr = New(fmt.Sprint(value))
	}
	return valAsBoolPtr
}

func ConvertInterfaceToStringMultiple(values []interface{}) []*string {
	var output []*string
	for _, value := range values {
		output = append(output, ConvertInterfaceToString(value))
	}
	return output
}
