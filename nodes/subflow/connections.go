package subflow

import "github.com/NubeDev/flow-eng/node"

type InputFloat struct {
	*node.Spec
}

func NewSubFlowInputFloat(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, inputFloat, category)
	in := node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil)
	in.FolderExport = true
	out := node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &InputFloat{body}, nil
}

func (inst *InputFloat) Process() {
	v, _ := inst.ReadPinAsFloat(node.Inp)
	inst.WritePinFloat(node.Outp, v)
}

type InputBool struct {
	*node.Spec
}

func NewSubFlowInputBool(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, inputBool, category)
	in := node.BuildInput(node.Inp, node.TypeBool, nil, body.Inputs, nil)
	in.FolderExport = true
	out := node.BuildOutput(node.Outp, node.TypeBool, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &InputBool{body}, nil
}

func (inst *InputBool) Process() {
	v, _ := inst.ReadPinAsBool(node.Inp)
	inst.WritePinBool(node.Outp, v)
}

type InputString struct {
	*node.Spec
}

func NewSubFlowInputString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, inputString, category)
	in := node.BuildInput(node.Inp, node.TypeString, nil, body.Inputs, nil)
	in.FolderExport = true
	out := node.BuildOutput(node.Outp, node.TypeString, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &InputString{body}, nil
}

func (inst *InputString) Process() {
	inst.WritePin(node.Outp, inst.ReadPin(node.Inp))
}

type OutputFloat struct {
	*node.Spec
}

func NewSubFlowOutputFloat(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, outputFloat, category)
	in := node.BuildInput(node.Inp, node.TypeFloat, nil, body.Inputs, nil)
	in.FolderExport = true
	out := node.BuildOutput(node.Outp, node.TypeFloat, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &OutputFloat{body}, nil
}

func (inst *OutputFloat) Process() {
	v, _ := inst.ReadPinAsFloat(node.Inp)
	inst.WritePinFloat(node.Outp, v)
}

type OutputBool struct {
	*node.Spec
}

func NewSubFlowOutputBool(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, outputBool, category)
	in := node.BuildInput(node.Inp, node.TypeBool, nil, body.Inputs, nil)
	in.FolderExport = true
	out := node.BuildOutput(node.Outp, node.TypeBool, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &OutputBool{body}, nil
}

func (inst *OutputBool) Process() {
	v, _ := inst.ReadPinAsBool(node.Inp)
	inst.WritePinBool(node.Outp, v)
}

type OutputString struct {
	*node.Spec
}

func NewSubFlowOutputString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, outputString, category)
	in := node.BuildInput(node.Inp, node.TypeString, nil, body.Inputs, nil)
	in.FolderExport = true
	out := node.BuildOutput(node.Outp, node.TypeString, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &OutputString{body}, nil
}

func (inst *OutputString) Process() {
	inst.WritePin(node.Outp, inst.ReadPin(node.Inp))
}
