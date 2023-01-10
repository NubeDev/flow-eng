package filter

import (
	"fmt"

	"github.com/NubeDev/flow-eng/node"
)

type PreventNull struct {
	*node.Spec
	lastValue float64
}

type PreventEqualFloat struct {
	*node.Spec
	lastValue float64
}

type PreventEqualString struct {
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

func NewPreventNull(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventNull, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘null’. Output the last known value if the input became null. lastValue is defaulted to be 0."))
	return &PreventNull{body, 0}, nil
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
// lastValue is defaulted to be 0
func (inst *PreventNull) Process() {
	v, null := inst.ReadPinAsFloat(node.In)

	if null {
		inst.WritePinFloat(node.Out, inst.lastValue)
	} else {
		inst.WritePinFloat(node.Out, v)
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

// preventDuplicates, don't do anything if the subsequent input value is the same. Out put nil if the input is nil.
// lastValue is defaulted to be 0
func (inst *PreventDuplicates) Process() {
	v, vNull := inst.ReadPinAsFloat(node.In)

	if vNull {
		inst.WritePinNull(node.Out)
	} else {
		if toFixed(v, 2) != inst.lastValue {
			inst.WritePinFloat(node.Out, v)
			inst.lastValue = v
		}
	}
}
