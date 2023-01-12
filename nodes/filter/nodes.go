package filter

import (
	"fmt"

	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
)

type PreventNull struct {
	*node.Spec
	lastValue float64
}

type PreventEqualFloat struct {
	*node.Spec
}

type PreventEqualString struct {
	*node.Spec
}

type OnlyBetween struct {
	*node.Spec
	lastValue float64
}

type OnlyGreater struct {
	*node.Spec
}

type OnlyLower struct {
	*node.Spec
}

type PreventDuplicates struct {
	*node.Spec
	lastValue *float64
}

func NewPreventNull(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventNull, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil))
	outputs := node.BuildOutputs(node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘null’. Output the last known value if the input became null. lastValue is defaulted to be 0."))
	return &PreventNull{body, 0}, nil
}

func NewPreventEqualFloat(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventEqualFloat, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil), node.BuildInput(node.Match, node.TypeFloat, nil, body.Inputs, nil))
	outputs := node.BuildOutputs(node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘input’ values which are equal to ‘match’. Output the last known input float value if two inputs equal to each other. Out put nil if any of the inputs are nil."))
	return &PreventEqualFloat{body}, nil
}

func NewPreventEqualString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventEqualString, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeString, nil, body.Inputs, nil), node.BuildInput(node.Match, node.TypeString, nil, body.Inputs, nil))
	outputs := node.BuildOutputs(node.BuildOutput(node.Outp, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘input’ values which are equal to ‘match’. Output the last known input string value if two inputs equal to each other. Out put nil if any of the inputs are nil."))
	return &PreventEqualString{body}, nil
}

func NewOnlyBetween(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyBetween, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil), node.BuildInput(node.MinInput, node.TypeFloat, nil, body.Inputs, nil), node.BuildInput(node.MaxInput, node.TypeFloat, nil, body.Inputs, nil))
	outputs := node.BuildOutputs(node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only Numeric ‘input’ values between ‘min’ and ‘max’ are passed to ‘output’. Output the last known input float value if the input is not in range. Out put nil if any of the inputs are nil."))
	return &OnlyBetween{body, 0}, nil
}

func NewOnlyGreater(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyGreater, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil), node.BuildInput(node.Threshold, node.TypeFloat, nil, body.Inputs, nil))
	outputs := node.BuildOutputs(node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only Numeric ‘input’ values greater than ‘OutMax’ are passed to ‘output’. Output the last known input float value if the input is less than 'OutMax'. Out put nil if any of the inputs are nil."))
	return &OnlyGreater{body}, nil
}

func NewOnlyLower(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyLower, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil), node.BuildInput(node.Threshold, node.TypeFloat, nil, body.Inputs, nil))
	outputs := node.BuildOutputs(node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only Numeric ‘input’ values less than 'InMax' are passed to ‘output’. Output the last known input float value if the input is greater than 'OutMin'. Out put nil if any of the inputs are nil."))
	return &OnlyLower{body}, nil
}

func NewPreventDuplicates(body *node.Spec) (node.Node, error) {
	var init float64 = 0
	body = node.Defaults(body, preventDuplicates, category)
	inputs := node.BuildInputs(node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil))
	outputs := node.BuildOutputs(node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘input’ values which are equal to the previous ‘input’ value. Output the last known input float value if the subsequent input value is the same. Out put nil if the input is nil."))
	return &PreventDuplicates{body, &init}, nil
}

// preventNull, output the last known value if the input became null.
// lastValue is defaulted to be 0
func (inst *PreventNull) Process() {
	input, null := inst.ReadPinAsFloat(node.Inp)

	if null {
		inst.WritePinFloat(node.Outp, inst.lastValue)
	} else {
		inst.WritePinFloat(node.Outp, input)
		inst.lastValue = input
	}
}

// preventEqualFloat, output the last known value if two inputs equal to each other.
// filter outputs lastValue if match is null
// lastValue is defaulted to be null
func (inst *PreventEqualFloat) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.Inp)
	match, matchNull := inst.ReadPinAsFloat(node.Match)

	if matchNull && !inputNull { // if the `match` value is null, and the input value isn't, then the value passes through.
		inst.WritePinFloat(node.Outp, input)
	} else if !matchNull && inputNull { // if the `match` value isn't null, and the input value is null, then the null passes through.
		inst.WritePinNull(node.Outp)
	} else if !matchNull && !inputNull && float.RoundTo(input, 4) != float.RoundTo(match, 4) { // if the `match` value and the input value aren't null, and the values don't match, then the input value passes through.
		inst.WritePinFloat(node.Outp, input)
	} // Otherwise the values must match and so the output value doesn't change

}

// preventEqualString, output the last known value if two inputs equal to each other.
// filter outputs lastValue if match is null
// lastValue is defaulted to be "defaultString"
func (inst *PreventEqualString) Process() {
	input, inputNull := inst.ReadPinAsString(node.Inp)
	match, matchNull := inst.ReadPinAsString(node.Match)

	// TODO: need to confirm if there is a null value for strings (ie. and empty string is treated as null)
	if matchNull && !inputNull { // if the `match` value is null, and the input value isn't, then the value passes through.
		inst.WritePin(node.Outp, input)
	} else if !matchNull && inputNull { // if the `match` value isn't null, and the input value is null, then the null passes through.
		inst.WritePinNull(node.Outp)
	} else if !matchNull && !inputNull && input != match { // if the `match` value and the input value aren't null, and the values don't match, then the input value passes through.
		inst.WritePin(node.Outp, input)
	} // Otherwise the values must match and so the output value doesn't change
}

// onlyBetween, Output the last known input float value if the input is not in range.
// filter behaves like onlyLower or onlyGreater if one of the Min and Max value is null.
// all values except for null are allow to pass through the filter if both Min and Max are null.
// lastValue is defaulted to be 0
func (inst *OnlyBetween) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.Inp)
	min, minNull := inst.ReadPinAsFloat(node.MinInput)
	max, maxNull := inst.ReadPinAsFloat(node.MaxInput)

	if inputNull {
		inst.WritePinFloat(node.Outp, inst.lastValue)
	} else {
		if minNull && maxNull {
			inst.WritePinFloat(node.Outp, input)
			inst.lastValue = input
		} else if !minNull && maxNull {
			if input >= min {
				inst.WritePinFloat(node.Outp, input)
				inst.lastValue = input
			}
		} else if minNull && !maxNull {
			if input <= max {
				inst.WritePinFloat(node.Outp, input)
				inst.lastValue = input
			}
		} else {
			if input >= min && input <= max {
				inst.WritePinFloat(node.Outp, input)
				inst.lastValue = input
			}
		}
	}
}

// onlyLower, Output the last known input float value if the input is greater than 'threshold'.
// filter outputs lastValue if min(threshold) is null
// lastValue is defaulted to be null
func (inst *OnlyGreater) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.Inp)
	min, minNull := inst.ReadPinAsFloat(node.Threshold)

	if inputNull {
		inst.WritePinNull(node.Outp)
	} else {
		if minNull || input > min {
			inst.WritePinFloat(node.Outp, input)
		}
	}
}

// onlyLower, Output the last known input float value if the input is lower than 'threshold'.
// filter outputs lastValue if max(threshold) is null
// lastValue is defaulted to be null
func (inst *OnlyLower) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.Inp)
	max, maxNull := inst.ReadPinAsFloat(node.Threshold)

	if inputNull {
		inst.WritePinNull(node.Outp)
	} else {
		if maxNull || input < max {
			inst.WritePinFloat(node.Outp, input)
		}
	}
}

// preventDuplicates, don't do anything if the subsequent input value is the same. Out put nil if the input is nil.
// lastValue is defaulted to be the address of value 0
func (inst *PreventDuplicates) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.Inp)

	if inputNull && inst.lastValue != nil {
		inst.WritePinNull(node.Outp)
		inst.lastValue = nil
	} else if (inst.lastValue == nil && !inputNull) || (inst.lastValue != nil && *inst.lastValue != input) {
		// if last value is null and input is not null
		// if last value is not null and last value is not the same as current input
		inst.WritePin(node.Outp, input)
		inst.lastValue = float.New(input)
	}
}
