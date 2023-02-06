package filter

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/str"

	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
)

type PreventNull struct {
	*node.Spec
	lastValue float64
}

type PreventEqualFloat struct {
	*node.Spec
	lastValue *float64
}

type PreventEqualString struct {
	*node.Spec
	lastValue *string
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
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘null’. Output the last known value if the input became null. lastValue is defaulted to be 0."))
	return &PreventNull{body, 0}, nil
}

func NewPreventEqualFloat(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventEqualFloat, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false), node.BuildInput(node.Match, node.TypeFloat, nil, body.Inputs, false))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘input’ values which are equal to ‘match’. Output the last known input float value if two inputs equal to each other. Out put nil if any of the inputs are nil."))
	return &PreventEqualFloat{body, nil}, nil
}

func NewPreventEqualString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventEqualString, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeString, nil, body.Inputs, false), node.BuildInput(node.Match, node.TypeString, nil, body.Inputs, false))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘input’ values which are equal to ‘match’. Output the last known input string value if two inputs equal to each other. Out put nil if any of the inputs are nil."))
	return &PreventEqualString{body, str.New("")}, nil
}

func NewOnlyBetween(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyBetween, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false), node.BuildInput(node.MinInput, node.TypeFloat, nil, body.Inputs, false), node.BuildInput(node.MaxInput, node.TypeFloat, nil, body.Inputs, false))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only Numeric ‘input’ values between ‘min’ and ‘max’ are passed to ‘output’. Output the last known input float value if the input is not in range. Out put nil if any of the inputs are nil."))
	return &OnlyBetween{body, 0}, nil
}

func NewOnlyGreater(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyGreater, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false), node.BuildInput(node.Threshold, node.TypeFloat, nil, body.Inputs, false))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only Numeric ‘input’ values greater than ‘OutMax’ are passed to ‘output’. Output the last known input float value if the input is less than 'OutMax'. Out put nil if any of the inputs are nil."))
	return &OnlyGreater{body, 0}, nil
}

func NewOnlyLower(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, onlyLower, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false), node.BuildInput(node.Threshold, node.TypeFloat, nil, body.Inputs, false))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. Only Numeric ‘input’ values less than 'InMax' are passed to ‘output’. Output the last known input float value if the input is greater than 'OutMin'. Out put nil if any of the inputs are nil."))
	return &OnlyLower{body, 0}, nil
}

func NewPreventDuplicates(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, preventDuplicates, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetHelp(fmt.Sprintln("This node filters ‘input’ values. All ‘input’ values are passed to ‘output’ EXCEPT ‘input’ values which are equal to the previous ‘input’ value. Output the last known input float value if the subsequent input value is the same. Out put nil if the input is nil."))
	return &PreventDuplicates{body, nil}, nil
}

// PreventNull, output the last known value if the input became null.
// lastValue is defaulted to be 0
func (inst *PreventNull) Process() {
	input, null := inst.ReadPinAsFloat(node.In)

	if null {
		inst.WritePinFloat(node.Out, inst.lastValue)
	} else {
		inst.WritePinFloat(node.Out, input)
		inst.lastValue = input
	}
}

// PreventEqualFloat, output the last known value if two inputs equal to each other.
// filter outputs lastValue if match is null
// lastValue is defaulted to be null
func (inst *PreventEqualFloat) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.In)
	match, matchNull := inst.ReadPinAsFloat(node.Match)

	if matchNull && !inputNull { // if the `match` value is null, and the input value isn't, then the value passes through.
		inst.WritePinFloat(node.Out, input)
		inst.lastValue = float.New(input)
	} else if !matchNull && inputNull { // if the `match` value isn't null, and the input value is null, then the null passes through.
		inst.WritePinNull(node.Out)
		inst.lastValue = nil
	} else if !matchNull && !inputNull && input != match { // if the `match` value and the input value aren't null, and the values don't match, then the input value passes through.
		inst.WritePinFloat(node.Out, input)
		inst.lastValue = float.New(input)
	} else { // Otherwise the values must match and so the output value doesn't change
		if inst.lastValue == nil {
			inst.WritePinNull(node.Out)
		} else {
			inst.WritePinFloat(node.Out, *inst.lastValue)
		}
	}
}

// PreventEqualString, output the last known value if two inputs equal to each other.
// filter outputs lastValue if match is null
// lastValue is defaulted to be "defaultString"
func (inst *PreventEqualString) Process() {
	input, inputNull := inst.ReadPinAsString(node.In)
	match, matchNull := inst.ReadPinAsString(node.Match)

	// TODO: need to confirm if there is a null value for strings (ie. and empty string is treated as null)
	if matchNull && !inputNull { // if the `match` value is null, and the input value isn't, then the value passes through.
		inst.WritePin(node.Out, input)
		inst.lastValue = str.New(input)
	} else if !matchNull && inputNull { // if the `match` value isn't null, and the input value is null, then the null passes through.
		inst.WritePinNull(node.Out)
		inst.lastValue = nil
	} else if !matchNull && !inputNull && input != match { // if the `match` value and the input value aren't null, and the values don't match, then the input value passes through.
		inst.WritePin(node.Out, input)
		inst.lastValue = str.New(input)
	} else { // Otherwise the values must match and so the output value doesn't change
		if inst.lastValue == nil {
			inst.WritePinNull(node.Out)
		} else {
			inst.WritePin(node.Out, *inst.lastValue)
		}
	}
}

// OnlyBetween, Output the last known input float value if the input is not in range.
// filter behaves like onlyLower or onlyGreater if one of the Min and Max value is null.
// all values except for null are allow to pass through the filter if both Min and Max are null.
// lastValue is defaulted to be 0
func (inst *OnlyBetween) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.In)
	min, minNull := inst.ReadPinAsFloat(node.MinInput)
	max, maxNull := inst.ReadPinAsFloat(node.MaxInput)

	if inputNull {
		inst.WritePinFloat(node.Out, inst.lastValue)
	} else {
		if minNull && maxNull {
			inst.WritePinFloat(node.Out, input)
			inst.lastValue = input
		} else if !minNull && maxNull {
			if input >= min {
				inst.WritePinFloat(node.Out, input)
				inst.lastValue = input
			}
		} else if minNull && !maxNull {
			if input <= max {
				inst.WritePinFloat(node.Out, input)
				inst.lastValue = input
			}
		} else {
			if input >= min && input <= max {
				inst.WritePinFloat(node.Out, input)
				inst.lastValue = input
			}
		}
		inst.WritePinFloat(node.Out, inst.lastValue)
	}
}

// OnlyGreater, Output the last known input float value if the input is greater than 'threshold'.
// filter outputs lastValue if min(threshold) is null
// lastValue is defaulted to be null
func (inst *OnlyGreater) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.In)
	min, minNull := inst.ReadPinAsFloat(node.Threshold)

	if inputNull {
		inst.WritePinNull(node.Out)
	} else {
		if minNull || input > min {
			inst.WritePinFloat(node.Out, input)
		}
		inst.WritePinFloat(node.Out, inst.lastValue)
	}
}

// OnlyLower, Output the last known input float value if the input is lower than 'threshold'.
// filter outputs lastValue if max(threshold) is null
// lastValue is defaulted to be null
func (inst *OnlyLower) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.In)
	max, maxNull := inst.ReadPinAsFloat(node.Threshold)

	if inputNull {
		inst.WritePinNull(node.Out)
	} else {
		if maxNull || input < max {
			inst.WritePinFloat(node.Out, input)
		}
		inst.WritePinFloat(node.Out, inst.lastValue)
	}
}

// PreventDuplicates, don't do anything if the subsequent input value is the same. Out put nil if the input is nil.
// lastValue is defaulted to be the address of value 0
func (inst *PreventDuplicates) Process() {
	input, inputNull := inst.ReadPinAsFloat(node.In)

	if inputNull && inst.lastValue != nil {
		inst.WritePinNull(node.Out)
		inst.lastValue = nil
	} else if (inst.lastValue == nil && !inputNull) || (inst.lastValue != nil && *inst.lastValue != input) {
		// if last value is null and input is not null
		// if last value is not null and last value is not the same as current input
		inst.WritePinFloat(node.Out, input)
		inst.lastValue = float.New(input)
	}
}
