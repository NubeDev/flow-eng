package subflow

import "github.com/NubeDev/flow-eng/node"

type InputFloat struct {
	*node.Spec
}

func NewSubFlowInputFloat(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, inputFloat, Category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false)
	in.FolderExport = true
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &InputFloat{body}, nil
}

func (inst *InputFloat) Process() {
	v, null := inst.ReadPinAsFloat(node.In)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePinFloat(node.Out, v)
	}
}

type InputBool struct {
	*node.Spec
}

func NewSubFlowInputBool(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, inputBool, Category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs, false, false)
	in.FolderExport = true
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &InputBool{body}, nil
}

func (inst *InputBool) Process() {
	v, null := inst.ReadPinAsBool(node.In)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePinBool(node.Out, v)
	}
}

type InputString struct {
	*node.Spec
}

func NewSubFlowInputString(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, inputString, Category)
	in := node.BuildInput(node.In, node.TypeString, nil, body.Inputs, false, false)
	in.FolderExport = true
	out := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &InputString{body}, nil
}

func (inst *InputString) Process() {
	inst.WritePin(node.Out, inst.ReadPin(node.In))
}

type OutputFloat struct {
	*node.Spec
}

func NewSubFlowOutputFloat(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, outputFloat, Category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, false, false)
	in.FolderExport = true
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &OutputFloat{body}, nil
}

func (inst *OutputFloat) Process() {
	v, null := inst.ReadPinAsFloat(node.In)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePinFloat(node.Out, v)
	}
}

type OutputBool struct {
	*node.Spec
}

func NewSubFlowOutputBool(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, outputBool, Category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs, false, false)
	in.FolderExport = true
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &OutputBool{body}, nil
}

func (inst *OutputBool) Process() {
	v, null := inst.ReadPinAsBool(node.In)
	if null {
		inst.WritePinNull(node.Out)
	} else {
		inst.WritePinBool(node.Out, v)
	}
}

type OutputString struct {
	*node.Spec
}

func NewSubFlowOutputString(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, outputString, Category)
	in := node.BuildInput(node.In, node.TypeString, nil, body.Inputs, false, false)
	in.FolderExport = true
	out := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &OutputString{body}, nil
}

func (inst *OutputString) Process() {
	inst.WritePin(node.Out, inst.ReadPin(node.In))
}
