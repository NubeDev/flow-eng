package subflow

import "github.com/NubeDev/flow-eng/node"

type InputFloat struct {
	*node.Spec
}

func NewSubFlowInputFloat(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, inputFloat, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, nil)
	in.FolderExport = true
	in.HideInput = true
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &InputFloat{body}, nil
}

func (inst *InputFloat) Process() {
	inst.WritePin(node.Out, inst.ReadPin(node.In))
}

type InputBool struct {
	*node.Spec
}

func NewSubFlowInputBool(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, inputBool, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs, nil)
	in.FolderExport = true
	in.HideInput = true
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &InputBool{body}, nil
}

func (inst *InputBool) Process() {
	inst.WritePin(node.Out, inst.ReadPin(node.In))
}

type InputString struct {
	*node.Spec
}

func NewSubFlowInputString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, inputString, category)
	in := node.BuildInput(node.In, node.TypeString, nil, body.Inputs, nil)
	in.FolderExport = true
	in.HideInput = true
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

func NewSubFlowOutputFloat(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, outputFloat, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, nil)
	in.FolderExport = true
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	out.HideOutput = true
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &OutputFloat{body}, nil
}

func (inst *OutputFloat) Process() {
	inst.WritePin(node.Out, inst.ReadPin(node.In))
}

type OutputBool struct {
	*node.Spec
}

func NewSubFlowOutputBool(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, outputBool, category)
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs, nil)
	in.FolderExport = true
	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	out.HideOutput = true
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &OutputBool{body}, nil
}

func (inst *OutputBool) Process() {
	inst.WritePin(node.Out, inst.ReadPin(node.In))
}

type OutputString struct {
	*node.Spec
}

func NewSubFlowOutputString(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, outputString, category)
	in := node.BuildInput(node.In, node.TypeString, nil, body.Inputs, nil)
	in.FolderExport = true
	out := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs)
	out.HideOutput = true
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &OutputString{body}, nil
}

func (inst *OutputString) Process() {
	inst.WritePin(node.Out, inst.ReadPin(node.In))
}
