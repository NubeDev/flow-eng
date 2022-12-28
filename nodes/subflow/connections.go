package subflow

import "github.com/NubeDev/flow-eng/node"

type InputFloat struct {
	*node.Spec
}

func NewSubFlowInputFloat(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, inputFloat, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs)
	in.FolderExport = true
	out := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &InputFloat{body}, nil
}

func NewSubFlowOutputFloat(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, outputFloat, category)
	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs)
	in.FolderExport = true
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
	in := node.BuildInput(node.In, node.TypeBool, nil, body.Inputs)
	in.FolderExport = true
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
	in := node.BuildInput(node.In, node.TypeString, nil, body.Inputs)
	in.FolderExport = true
	out := node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs)
	body = node.BuildNode(body, node.BuildInputs(in), node.BuildOutputs(out), nil)
	return &InputString{body}, nil
}

func (inst *InputString) Process() {
	inst.WritePin(node.Out, inst.ReadPin(node.In))
}
