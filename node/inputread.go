package node

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/conversions"
)

// InputUpdated if true means that the node input value has updated
func (n *Spec) InputUpdated(name InputName) bool {
	input := n.GetInput(name)
	if input != nil {
		return input.updated
	}
	return false
}

func (n *Spec) ReadPin(name InputName) interface{} {
	input := n.GetInput(name)
	if input == nil {
		return nil
	}
	if input.Connection.OverrideValue != nil { // this would be that the user wrote a value to the input directly
		return input.Connection.OverrideValue
	}
	if input.Connection.FallbackValue != nil { // this would be that the user wrote a value to the input directly
		return input.Connection.FallbackValue
	}

	return input.GetValue()
}

func (n *Spec) ReadPinAsFloat(name InputName) float64 {
	r := n.ReadPin(name)
	out := conversions.GetFloat(r)
	return out
}

func (n *Spec) ReadPinAsString(name InputName) string {
	r := n.ReadPin(name)
	return fmt.Sprintf("%v", r)
}

func (n *Spec) ReadPinAsInt(name InputName) int {
	r := n.ReadPin(name)
	out := conversions.GetInt(r)
	return out
}

func (n *Spec) ReadMultiple(count int) []interface{} {
	var out []interface{}
	for i, input := range n.GetInputs() {
		if i < count {
			out = append(out, n.ReadPin(input.Name))
		}
	}
	return out
}
