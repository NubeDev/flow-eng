package node

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/str"
)

func (n *BaseNode) ReadMultipleNums(count int) []float64 {
	var out []float64
	inputs := n.ReadMultiple(count)
	for _, input := range inputs {
		val, realValue, _ := n.ReadPinNum(input.Name)
		if val != nil {
			out = append(out, realValue)
		}
	}
	return out
}

func (n *BaseNode) ReadMultiple(count int) []*Input {
	var out []*Input
	for i, input := range n.GetInputs() {
		if i < count {
			out = append(out, input)
		}
	}
	return out
}

func (n *BaseNode) ReadPinNum(name InputName) (value *float64, real float64, hasValue bool) {
	pinValPointer, _ := n.ReadPin(name)
	valPointer, val, err := float.StringFloatErr(pinValPointer)
	if err != nil {
		return nil, 0, hasValue
	}
	return valPointer, val, float.NotNil(valPointer)
}

func (n *BaseNode) ReadPin(name InputName) (*string, string) {
	in := n.GetInput(name)
	if in == nil {
		return nil, ""
	}
	if name == in.Name {
		if in.Connection.OverrideValue != nil { // this would be that the user wrote a value to the input directly
			toStr := fmt.Sprintf("%v", in.Connection.OverrideValue)
			return str.New(toStr), str.NonNil(str.New(toStr))
		}
		if in.Connection.FallbackValue != nil { // this would be that the user wrote a value to the input directly
			toStr := fmt.Sprintf("%v", in.Connection.FallbackValue)
			return str.New(toStr), str.NonNil(str.New(toStr))
		}
		val := fmt.Sprintf("%v", in.Value)
		return str.New(val), val
	}

	return nil, ""
}
