package flow

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/driver"
)

type Device struct {
	*node.Spec
	firstLoop   bool
	networkUUID string
	pool        driver.Driver
}

func NewDevice(body *node.Spec, pool driver.Driver) (node.Node, error) {
	body = node.Defaults(body, flowDevice, category)
	name := node.BuildInput(node.Name, node.TypeString, nil, body.Inputs)
	networkUUID := node.BuildInput(node.UUID, node.TypeString, nil, body.Inputs)
	inputs := node.BuildInputs(name, networkUUID)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body.IsParent = true
	body = node.BuildNode(body, inputs, outputs, nil)
	n, _ := body.ReadPinAsString(node.UUID)
	return &Device{body, false, n, pool}, nil
}

func (inst *Device) setConnection() {
	net := inst.pool.GetNetwork(inst.networkUUID)
	if net != nil {
		n, _ := inst.ReadPinAsString(node.Name)
		inst.pool.AddDevice(inst.networkUUID, &driver.Device{
			UUID: n,
			Name: n,
		})
		inst.firstLoop = true
	} else {
	}

}

func (inst *Device) Process() {
	if !inst.firstLoop {
		inst.setConnection()
	}
}
func (inst *Device) Cleanup() {}
