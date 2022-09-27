package flow

import (
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/driver"
)

type Device struct {
	*node.Spec
}

func NewDevice(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowDevice, category)
	networkName := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs)
	value := node.BuildInput(node.In, node.TypeString, nil, body.Inputs)
	inputs := node.BuildInputs(networkName, value)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	if network == nil {
		network = driver.New(&driver.Network{
			Name:        string(names.FlowFramework),
			Application: names.FlowFramework,
			Storage:     body.GetDB(),
		})
	}
	return &Device{body}, nil
}

func (inst *Device) Process() {

}

func (inst *Device) Cleanup() {}
