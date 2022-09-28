package flow

import (
	"github.com/NubeDev/flow-eng/node"
)

type Point struct {
	*node.Spec
	firstLoop bool
}

func NewPoint(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowPoint, category)
	networkName := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs)
	value := node.BuildInput(node.In, node.TypeString, nil, body.Inputs)
	inputs := node.BuildInputs(networkName, value)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Point{body, false}, nil
}

func (inst *Point) Process() {

}

func (inst *Point) Cleanup() {}
