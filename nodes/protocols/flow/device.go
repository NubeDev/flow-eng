package flow

import (
	"github.com/NubeDev/flow-eng/node"
)

type Device struct {
	*node.Spec
	firstLoop bool
}

func NewDevice(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowDevice, category)
	networkName := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs)
	value := node.BuildInput(node.In, node.TypeString, nil, body.Inputs)
	inputs := node.BuildInputs(networkName, value)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Device{body, false}, nil
}

func (inst *Device) Process() {

}

func (inst *Device) Cleanup() {}
