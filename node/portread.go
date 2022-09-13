package node

import (
	"errors"
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

type RedMultiplePins struct {
	Value *float64
	Real  float64
	Found bool
}

func (n *BaseNode) ReadPinsNum(name ...InputName) []*RedMultiplePins {
	var out []*RedMultiplePins
	var resp *RedMultiplePins
	for _, portName := range name {
		v, r, f := n.ReadPinNum(portName)
		resp.Value = v
		resp.Real = r
		resp.Found = f
		out = append(out, resp)
	}
	return out
}

func (n *BaseNode) OverrideInputValue(name InputName, value interface{}) error {
	in := n.GetInput(name)
	if in == nil {
		return errors.New(fmt.Sprintf("failed to find port%s", name))
	}
	if in.Connection != nil {
		in.Connection.OverrideValue = value
	} else {
		return errors.New(fmt.Sprintf("this node has no inputs"))
	}
	return nil

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
