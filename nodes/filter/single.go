package filter

import (
	"fmt"
	"math"

	"github.com/NubeDev/flow-eng/node"
)

// type OnlyTrue struct {
// 	*node.Spec
// }

// type OnlyFalse struct {
// 	*node.Spec
// }

type PreventNull struct {
	*node.Spec
	lastValue bool
}

type PreventEqualFloat struct {
	*node.Spec
	lastValue float64
}

type PreventEqualString struct {
	*node.Spec
	lastValue string
}

type OnlyEqualFloat struct {
	*node.Spec
	lastValue float64
}

type OnlyEqualString struct {
	*node.Spec
	lastValue string
}

type OnlyBetween struct {
	*node.Spec
	lastValue float64
}

type OnlyGreater struct {
	*node.Spec
	lastValue float64
}

type OnlyLower struct {
	*node.Spec
	lastValue float64
}

type PreventDuplicates struct {
	*node.Spec
	lastValue float64
}

// func NewOnlyTrue(body *node.Spec) (node.Node, error) {
// 	body = node.Defaults(body, onlyTrue, category)
// 	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeBool, nil, body.Inputs))
// 	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
// 	body = node.BuildNode(body, inputs, outputs, nil)
// 	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only ‘true’ values are passed to ‘output’."))
// 	return &OnlyTrue{body}, nil
// }

// func NewOnlyFalse(body *node.Spec) (node.Node, error) {
// 	body = node.Defaults(body, onlyFalse, category)
// 	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeBool, nil, body.Inputs))
// 	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
// 	body = node.BuildNode(body, inputs, outputs, nil)
// 	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only ‘false’ values are passed to ‘output’."))
// 	return &OnlyFalse{body}, nil
// }

func NewPreventNull(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventNull, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeBool, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘null’. Output the last known value if the input became null."))
	return &PreventNull{body, false}, nil
}

func NewPreventEqualFloat(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventEqualFloat, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs), node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘input’ values which are equal to ‘match’. Output the last known input float value if two inputs equal to each other. Out put nil if any of the inputs are nil."))
	return &PreventEqualFloat{body, 0}, nil
}

func NewPreventEqualString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventEqualString, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeString, nil, body.Inputs), node.BuildInput(node.In1, node.TypeString, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘input’ values which are equal to ‘match’. Output the last known input string value if two inputs equal to each other. Out put nil if any of the inputs are nil."))
	return &PreventEqualString{body, "defaultString"}, nil
}

func NewOnlyEqualFloat(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyEqualFloat, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs), node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only ‘input’ values equal to ‘match’ are passed to ‘output’. Output the last known input float value if two inputs do not equal to each other. Out put nil if any of the inputs are nil."))
	return &OnlyEqualFloat{body, 0}, nil
}

func NewOnlyEqualString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyEqualString, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeString, nil, body.Inputs), node.BuildInput(node.In1, node.TypeString, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only ‘input’ values equal to ‘match’ are passed to ‘output’. Output the last known input string value if two inputs do not equal to each other. Out put nil if any of the inputs are nil."))
	return &OnlyEqualString{body, "defaultString"}, nil
}

func NewOnlyBetween(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyBetween, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs), node.BuildInput(node.InMin, node.TypeFloat, nil, body.Inputs), node.BuildInput(node.InMax, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only Numeric ‘input’ values between ‘min’ and ‘max’ are passed to ‘output’. Output the last known input float value if the input is not in range. Out put nil if any of the inputs are nil."))
	return &OnlyBetween{body, 0}, nil
}

func NewOnlyGreater(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyGreater, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs), node.BuildInput(node.InMin, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only Numeric ‘input’ values greater than ‘OutMax’ are passed to ‘output’. Output the last known input float value if the input is less than 'OutMax'. Out put nil if any of the inputs are nil."))
	return &OnlyGreater{body, 0}, nil
}

func NewOnlyLower(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyLower, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs), node.BuildInput(node.InMax, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only Numeric ‘input’ values less than 'InMax' are passed to ‘output’. Output the last known input float value if the input is greater than 'OutMin'. Out put nil if any of the inputs are nil."))
	return &OnlyLower{body, 0}, nil
}

func NewPreventDuplicates(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventDuplicates, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘input’ values which are equal to the previous ‘input’ value. Output the last known input float value if the subsequent input value is the same. Out put nil if the input is nil."))
	return &PreventDuplicates{body, 0}, nil
}

// func (inst *OnlyTrue) Process() {
// 	v, _ := inst.ReadPinAsBool(node.In)
// 	if v {
// 		inst.WritePinTrue(node.Out)
// 	} else {
// 		inst.WritePinFalse(node.Out)
// 	}
// }

// func (inst *OnlyFalse) Process() {
// 	v, null := inst.ReadPinAsBool(node.In)
// 	if !v || null {
// 		inst.WritePinFalse(node.Out)
// 	} else {
// 		inst.WritePinTrue(node.Out)
// 	}
// }

// preventNull, output the last known value if the input became null.
// last value is defaulted to be false
func (inst *PreventNull) Process() {
	v, null := inst.ReadPinAsBool(node.In)

	if null {
		inst.WritePinBool(node.Out, inst.lastValue)
	} else {
		inst.WritePinBool(node.Out, v)
		inst.lastValue = v
	}
}

// preventEqualFloat, output the last known value if two inputs equal to each other. Out put nil if any of the inputs are nil.
// lastValue is defaulted to be 0
func (inst *PreventEqualFloat) Process() {
	v, vNull := inst.ReadPinAsFloat(node.In)
	m, mNull := inst.ReadPinAsFloat(node.In1)

	if vNull || mNull {
		inst.WritePinNull(node.Out)
	} else {
		if toFixed(v, 2) != toFixed(m, 2) {
			inst.WritePinFloat(node.Out, v)
			inst.lastValue = v
		} else {
			inst.WritePinFloat(node.Out, inst.lastValue)
		}
	}
}

// preventEqualString, output the last known value if two inputs equal to each other. Out put nil if any of the inputs are nil.
// lastValue is defaulted to be "defaultString"
func (inst *PreventEqualString) Process() {
	v, vNull := inst.ReadPinAsString(node.In)
	m, mNull := inst.ReadPinAsString(node.In1)

	if vNull || mNull {
		inst.WritePinNull(node.Out)
	} else {
		if v != m {
			inst.WritePin(node.Out, v)
			inst.lastValue = v
		} else {
			inst.WritePin(node.Out, inst.lastValue)
		}
	}
}

// onlyEqualFloat, output the last known value if two inputs do not equal to each other. Out put nil if any of the inputs are nil.
// lastValue is defaulted to be 0
func (inst *OnlyEqualFloat) Process() {
	v, vNull := inst.ReadPinAsFloat(node.In)
	m, mNull := inst.ReadPinAsFloat(node.In1)

	if vNull || mNull {
		inst.WritePinNull(node.Out)
	} else {
		if toFixed(v, 2) == toFixed(m, 2) {
			inst.WritePinFloat(node.Out, v)
			inst.lastValue = v
		} else {
			inst.WritePinFloat(node.Out, inst.lastValue)
		}
	}
}

// onlyEqualString, output the last known value if two inputs do not equal to each other. Out put nil if any of the inputs are nil.
// lastValue is defaulted to be "defaultString"
func (inst *OnlyEqualString) Process() {
	v, vNull := inst.ReadPinAsString(node.In)
	m, mNull := inst.ReadPinAsString(node.In1)

	if vNull || mNull {
		inst.WritePinNull(node.Out)
	} else {
		if v == m {
			inst.WritePin(node.Out, v)
			inst.lastValue = v
		} else {
			inst.WritePin(node.Out, inst.lastValue)
		}
	}
}

// onlyBetween, Output the last known input float value if the input is not in range. Out put nil if any of the inputs are nil.
// lastValue is defaulted to be 0
func (inst *OnlyBetween) Process() {
	v, vNull := inst.ReadPinAsFloat(node.In)
	min, minNull := inst.ReadPinAsFloat(node.InMin)
	max, maxNull := inst.ReadPinAsFloat(node.InMax)

	if vNull || minNull || maxNull {
		inst.WritePinNull(node.Out)
	} else {
		if toFixed(min, 2) <= toFixed(v, 2) && toFixed(v, 2) <= toFixed(max, 2) {
			inst.WritePinFloat(node.Out, v)
			inst.lastValue = v
		} else {
			inst.WritePinFloat(node.Out, inst.lastValue)
		}
	}
}

// onlyGreater, Output the last known input float value if the input is less than 'OutMax'. Out put nil if any of the inputs are nil.
// lastValue is defaulted to be 0
func (inst *OnlyGreater) Process() {
	v, vNull := inst.ReadPinAsFloat(node.In)
	min, minNull := inst.ReadPinAsFloat(node.InMin)

	if vNull || minNull {
		inst.WritePinNull(node.Out)
	} else {
		if toFixed(v, 2) > toFixed(min, 2) {
			inst.WritePinFloat(node.Out, v)
			inst.lastValue = v
		} else {
			inst.WritePinFloat(node.Out, inst.lastValue)
		}
	}
}

// onlyLower, Output the last known input float value if the input is greater than 'OutMin'. Out put nil if any of the inputs are nil.
// lastValue is defaulted to be 0
func (inst *OnlyLower) Process() {
	v, vNull := inst.ReadPinAsFloat(node.In)
	max, maxNull := inst.ReadPinAsFloat(node.InMax)

	if vNull || maxNull {
		inst.WritePinNull(node.Out)
	} else {
		if toFixed(v, 2) < toFixed(max, 2) {
			inst.WritePinFloat(node.Out, v)
			inst.lastValue = v
		} else {
			inst.WritePinFloat(node.Out, inst.lastValue)
		}
	}
}

// preventDuplicates, Output the last known input float value if the subsequent input value is the same. Out put nil if the input is nil.
// lastValue is defaulted to be 0
func (inst *PreventDuplicates) Process() {
	v, vNull := inst.ReadPinAsFloat(node.In)

	if vNull {
		inst.WritePinNull(node.Out)
	} else {
		if toFixed(v, 2) == inst.lastValue {
			inst.WritePinFloat(node.Out, inst.lastValue)
		} else {
			inst.WritePinFloat(node.Out, v)
			inst.lastValue = v
		}
	}
}

// helper functions that round the input number up
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
