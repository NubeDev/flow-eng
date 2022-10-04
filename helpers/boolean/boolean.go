package boolean

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
	if value == nil {
		return nil
	}
	output, ok := value.(bool)
	if ok {
		return &output
	}
	return nil
}

func ConvertInterfaceToBoolMultiple(values []interface{}) []*bool {
	var output []*bool
	for _, value := range values {
		output = append(output, ConvertInterfaceToBool(value))
	}
	return output
}
