package filter

import (
	"fmt"
	"reflect"

	"github.com/NubeDev/flow-eng/helpers/float"
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
	lastValue *float64
}

func NewPreventNull(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventNull, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘null’. Output the last known value if the input became null. lastValue is defaulted to be 0."))
	return &PreventNull{body, 0}, nil
}

func NewPreventEqualFloat(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventEqualFloat, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs), node.BuildInput(node.InputMatch, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘input’ values which are equal to ‘match’. Output the last known input float value if two inputs equal to each other. Out put nil if any of the inputs are nil."))
	return &PreventEqualFloat{body, 0}, nil
}

func NewPreventEqualString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventEqualString, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeString, nil, body.Inputs), node.BuildInput(node.InputMatch, node.TypeString, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘input’ values which are equal to ‘match’. Output the last known input string value if two inputs equal to each other. Out put nil if any of the inputs are nil."))
	return &PreventEqualString{body, "defaultString"}, nil
}

func NewOnlyBetween(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyBetween, category)
	inputs := node.BuildInputs(node.BuildInput(node.InputValue, node.TypeFloat, nil, body.Inputs), node.BuildInput(node.Min, node.TypeFloat, nil, body.Inputs), node.BuildInput(node.Max, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only Numeric ‘input’ values between ‘min’ and ‘max’ are passed to ‘output’. Output the last known input float value if the input is not in range. Out put nil if any of the inputs are nil."))
	return &OnlyBetween{body, 0}, nil
}

func NewOnlyGreater(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyGreater, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs), node.BuildInput(node.InputThresh, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only Numeric ‘input’ values greater than ‘OutMax’ are passed to ‘output’. Output the last known input float value if the input is less than 'OutMax'. Out put nil if any of the inputs are nil."))
	return &OnlyGreater{body, 0}, nil
}

func NewOnlyLower(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyLower, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs), node.BuildInput(node.InputThresh, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only Numeric ‘input’ values less than 'InMax' are passed to ‘output’. Output the last known input float value if the input is greater than 'OutMin'. Out put nil if any of the inputs are nil."))
	return &OnlyLower{body, 0}, nil
}

func NewPreventDuplicates(body *node.Spec) (node.Node, error) {
	var init float64 = 0
	body = node.Defaults(body, preventDuplicates, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘input’ values which are equal to the previous ‘input’ value. Output the last known input float value if the subsequent input value is the same. Out put nil if the input is nil."))
	return &PreventDuplicates{body, &init}, nil
}

// preventNull, output the last known value if the input became null.
// lastValue is defaulted to be 0
func (inst *PreventNull) Process() {
	input, null := inst.ReadPinAsFloat(node.Inp)

	if null {
		inst.WritePinFloat(node.Out, inst.lastValue)
	} else {
		inst.WritePinFloat(node.Out, input)
		inst.lastValue = input
	}
}

// preventEqualFloat, output the last known value if two inputs equal to each other.
// filter outputs lastValue if match is null
// lastValue is defaulted to be 0
func (inst *PreventEqualFloat) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.Inp)
	match, matchNull := inst.ReadPinAsFloat(node.InputMatch)

	if matchNull {
		inst.WritePinFloat(node.Out, inst.lastValue)
	} else {
		if inputNull {
			inst.WritePinNull(node.Out)
		} else {
			if float.RoundTo(input, 4) != float.RoundTo(match, 4) {
				inst.WritePinFloat(node.Out, input)
				inst.lastValue = input
			} else {
				inst.WritePinFloat(node.Out, inst.lastValue)
			}
		}
	}
}

// preventEqualString, output the last known value if two inputs equal to each other.
// filter outputs lastValue if match is null
// lastValue is defaulted to be "defaultString"
func (inst *PreventEqualString) Process() {
	input, inputNull := inst.ReadPinAsString(node.Inp)
	match, matchNull := inst.ReadPinAsString(node.InputMatch)

	if matchNull {
		inst.WritePin(node.Out, inst.lastValue)
	} else {
		if inputNull {
			inst.WritePinNull(node.Out)
		} else {
			if input != match {
				inst.WritePin(node.Out, input)
				inst.lastValue = input
			} else {
				inst.WritePin(node.Out, inst.lastValue)
			}
		}
	}
}

// onlyBetween, Output the last known input float value if the input is not in range.
// filter behaves like onlyLower or onlyGreater if one of the Min and Max value is null.
// all values except for null are allow to pass through the filter if both Min and Max are null.
// lastValue is defaulted to be 0
func (inst *OnlyBetween) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.InputValue)
	min, minNull := inst.ReadPinAsFloat(node.Min)
	max, maxNull := inst.ReadPinAsFloat(node.Max)

	if inputNull {
		inst.WritePinFloat(node.Out, inst.lastValue)
	} else {
		if minNull && maxNull {
			inst.WritePinFloat(node.Out, input)
			inst.lastValue = input
		} else if !minNull && maxNull {
			if float.RoundTo(min, 4) < float.RoundTo(input, 4) {
				inst.WritePinFloat(node.Out, input)
				inst.lastValue = input
			} else {
				inst.WritePinFloat(node.Out, inst.lastValue)
			}
		} else if minNull && !maxNull {
			if float.RoundTo(max, 4) > float.RoundTo(input, 4) {
				inst.WritePinFloat(node.Out, input)
				inst.lastValue = input
			} else {
				inst.WritePinFloat(node.Out, inst.lastValue)
			}
		} else {
			if float.RoundTo(min, 4) < float.RoundTo(input, 4) && float.RoundTo(input, 4) < float.RoundTo(max, 4) {
				inst.WritePinFloat(node.Out, input)
				inst.lastValue = input
			} else {
				inst.WritePinFloat(node.Out, inst.lastValue)
			}
		}
	}
}

// onlyLower, Output the last known input float value if the input is greater than 'threshold'.
// filter outputs lastValue if min(threshold) is null
// lastValue is defaulted to be 0
func (inst *OnlyGreater) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.Inp)
	min, minNull := inst.ReadPinAsFloat(node.InputThresh)

	if minNull {
		inst.WritePinFloat(node.Out, inst.lastValue)
	} else {
		if inputNull {
			inst.WritePinNull(node.Out)
		} else {
			if float.RoundTo(input, 4) > float.RoundTo(min, 4) {
				inst.WritePinFloat(node.Out, input)
				inst.lastValue = input
			} else {
				inst.WritePinFloat(node.Out, inst.lastValue)
			}
		}
	}
}

// onlyLower, Output the last known input float value if the input is lower than 'threshold'.
// filter outputs lastValue if max(threshold) is null
// lastValue is defaulted to be 0
func (inst *OnlyLower) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.Inp)
	max, maxNull := inst.ReadPinAsFloat(node.InputThresh)

	if maxNull {
		inst.WritePinFloat(node.Out, inst.lastValue)
	} else {
		if inputNull {
			inst.WritePinNull(node.Out)
		} else {
			if float.RoundTo(input, 4) < float.RoundTo(max, 4) {
				inst.WritePinFloat(node.Out, input)
				inst.lastValue = input
			} else {
				inst.WritePinFloat(node.Out, inst.lastValue)
			}
		}
	}
}

// preventDuplicates, don't do anything if the subsequent input value is the same. Out put nil if the input is nil.
// lastValue is defaulted to be the address of value 0
func (inst *PreventDuplicates) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.Inp)

	if inputNull {
		inst.WritePinNull(node.Out)
		inst.lastValue = nil
	} else {
		if !reflect.ValueOf(inst.lastValue).IsNil() {
			if float.RoundTo(input, 4) != *inst.lastValue {
				inst.WritePinFloat(node.Out, input)
				*inst.lastValue = input
			}
		}
	}
}
