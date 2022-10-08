package node

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/boolean"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"strconv"
	"time"
)

// InputUpdated if true means that the node input value has updated
func (n *Spec) InputUpdated(name InputName) (updated bool, boolCOV bool) {
	input := n.GetInput(name)

	if input.values.Length() > 1 { // work out if the input has updated
		input.values.RemoveFirst()
		input.values.Add(input.value)
	} else {
		input.values.Add(input.value)
	}
	if input.values.GetByIndex(0) != input.values.GetByIndex(1) {
		input.updated = true
	} else {
		input.updated = false
	}
	if input != nil {

		isBool, val := conversions.IsBoolWithValue(input.GetValue())
		if input.updated && isBool && val {
			boolCOV = true
		} else {
			boolCOV = false
		}

		return input.updated, boolCOV
	}
	return false, false
}

// InputHasConnection true if the node input has a connection
func (n *Spec) InputHasConnection(name InputName) bool {
	input := n.GetInput(name)
	if input == nil {
		return false
	}
	if input.Connection.NodeID != "" {
		return true
	}
	return false
}

func (n *Spec) ReadPin(name InputName) interface{} {
	input := n.GetInput(name)
	if input == nil {
		return nil
	}
	return input.GetValue()
}

func (n *Spec) ReadPinAsFloatPointer(name InputName) *float64 {
	r := n.ReadPin(name)
	return conversions.GetFloatPointer(r)
}

func (n *Spec) ReadPinAsFloatOk(name InputName) (value float64, null bool) {
	r := n.ReadPin(name)
	if r != nil {
		return conversions.GetFloat(r), false
	}
	return 0, true
}

func (n *Spec) ReadPinAsFloat(name InputName) (value float64) {
	r := n.ReadPin(name)
	out := conversions.GetFloat(r)
	return out
}

func (n *Spec) ReadPinAsDuration(name InputName) time.Duration {
	r := n.ReadPin(name)
	return time.Duration(conversions.GetInt(r))
}

func (n *Spec) ReadPinAsUint64(name InputName) uint64 {
	r := n.ReadPin(name)
	out := uint64(conversions.GetFloat(r))
	return out
}

func (n *Spec) ReadPinAsString(name InputName) string {
	r := n.ReadPin(name)
	return fmt.Sprintf("%v", r)
}

func (n *Spec) ReadPinBool(name InputName) bool {
	r := n.ReadPin(name)
	result, _ := strconv.ParseBool(fmt.Sprintf("%v", r))
	return result
}

func (n *Spec) ReadPinBoolOk(name InputName) (value bool, null bool) {
	r := n.ReadPin(name)
	if r != nil {
		result, _ := strconv.ParseBool(fmt.Sprintf("%v", r))
		return result, false
	}
	return false, true
}

func (n *Spec) ReadPinBoolPointer(name InputName) *bool {
	r := n.ReadPin(name)
	result, _ := strconv.ParseBool(fmt.Sprintf("%v", r))
	return boolean.New(result)
}

func (n *Spec) ReadPinAsInt(name InputName) int {
	r := n.ReadPin(name)
	out := conversions.GetInt(r)
	return out
}

func (n *Spec) ReadMultipleFloat(count int) []float64 {
	var out []float64
	for i, input := range n.GetInputs() {
		if i < count {
			out = append(out, n.ReadPinAsFloat(input.Name))
		}
	}
	return out
}

func (n *Spec) ReadMultipleFloatPointer(count int) []*float64 {
	var out []*float64
	for i, input := range n.GetInputs() {
		if i < count {
			out = append(out, n.ReadPinAsFloatPointer(input.Name))
		}
	}
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
