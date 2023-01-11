package node

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/helpers/integer"
	"reflect"
	"strconv"
	"time"
)

// InputUpdated if true means that the node input value has updated
func (n *Spec) InputUpdated(name InputName) (updated bool, boolCOV bool) {

	input := n.GetInput(name)
	if input == nil {
		return false, false
	}
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

	isBool, val := conversions.IsBoolWithValue(input.GetValue())
	if input.updated && isBool && val {
		boolCOV = true
	} else {
		boolCOV = false
	}

	return input.updated, boolCOV

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

func isNil(i interface{}) bool {
	return i == nil || reflect.ValueOf(i).IsNil()
}

// InputHasConnectionOrValue true if the node input has a connection, or there is a manual input
func (n *Spec) InputHasConnectionOrValue(name InputName) bool {
	if n.InputHasConnection(name) {
		return true
	}
	_, null := n.ReadPinAsString(name)
	if !null {
		return true
	}
	return false
}

func (n *Spec) ReadPinOrSettings(name InputName) interface{} {
	input := n.GetInput(name)
	connection := n.InputHasConnection(name)
	if !connection {
		if n.Settings != nil {
			val := reflect.ValueOf(n.Settings)
			if val.Kind() == reflect.Map {
				for _, e := range val.MapKeys() {
					v := val.MapIndex(e)
					if e.String() == string(name) { // add in case for string
						f, ok := conversions.GetFloatOk(v)
						if ok {
							return f
						}
						i, ok := conversions.GetIntOk(v)
						if ok {
							return i
						}
					}
				}
			}
		}
	}
	return input.GetValue()
}

func (n *Spec) ReadPin(name InputName) interface{} {
	input := n.GetInput(name)
	if input == nil {
		return nil
	}
	return input.GetValue()
}

func (n *Spec) ReadPinAsFloat(name InputName) (value float64, null bool) {
	r := n.ReadPin(name)
	if r == nil {
		return 0, true
	}
	return conversions.GetFloat(r), false
}

func (n *Spec) ReadInputPriority(name InputName) (value float64, null bool) {
	r := n.ReadPin(name)
	if r == nil {
		return 0, true
	}
	return conversions.GetFloat(r), false
}

func (n *Spec) readPinAsFloat(name InputName) (value float64) {
	r := n.ReadPin(name)
	out := conversions.GetFloat(r)
	return out
}

func (n *Spec) ReadPinAsDuration(name InputName) (value time.Duration, null bool) {
	r := n.ReadPin(name)
	if r == nil {
		return 0, true
	}
	fmt.Println(r)
	return time.Duration(conversions.GetInt(r)), false
}

func (n *Spec) ReadPinAsUint64(name InputName) (value uint64, null bool) {
	r := n.ReadPin(name)
	if r == nil {
		return 0, true
	}
	out := integer.GetUnit64(r)
	return out, false
}

func (n *Spec) ReadPinAsString(name InputName) (value string, null bool) {
	r := n.ReadPin(name)
	if r == nil {
		return "", true
	}
	return fmt.Sprintf("%v", r), false
}

func (n *Spec) readPinBool(name InputName) bool {
	r := n.ReadPin(name)
	result, _ := strconv.ParseBool(fmt.Sprintf("%v", r))
	return result
}

func (n *Spec) ReadPinAsBool(name InputName) (value bool, null bool) {
	r := n.ReadPin(name)
	if r == nil {
		return false, true
	}
	result, _ := strconv.ParseBool(fmt.Sprintf("%v", r))
	return result, false
}

func (n *Spec) ReadPinAsInt(name InputName) (value int, null bool) {
	r := n.ReadPin(name)
	if r == nil {
		return 0, true
	}
	out := conversions.GetInt(r)
	return out, false

}

func (n *Spec) ReadMultipleFloat(count int) []float64 {
	var out []float64
	for i, input := range n.GetInputs() {
		if i < count {
			out = append(out, n.readPinAsFloat(input.Name))
		}
	}
	return out
}

func (n *Spec) ReadPinAsFloatPointer(name InputName) *float64 {
	r := n.ReadPin(name)
	return conversions.GetFloatPointer(r)
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
