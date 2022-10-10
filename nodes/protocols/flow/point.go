package flow

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/driver"
)

type Point struct {
	*node.Spec
	firstLoop  bool
	deviceUUID string
	pool       driver.Driver
}

func NewPoint(body *node.Spec, pool driver.Driver) (node.Node, error) {
	body = node.Defaults(body, flowPoint, category)
	name := node.BuildInput(node.Name, node.TypeString, nil, body.Inputs)
	deviceUUID := node.BuildInput(node.UUID, node.TypeString, nil, body.Inputs)
	inputs := node.BuildInputs(name, deviceUUID)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	n, _ := body.ReadPinAsString(node.UUID)
	return &Point{body, false, n, pool}, nil
}

func (inst *Point) setConnection() {
	dev := inst.pool.GetDevice(inst.deviceUUID)
	if dev != nil {
		n, _ := inst.ReadPinAsString(node.Name)
		pnt := inst.pool.AddPoint(inst.deviceUUID, &driver.Point{
			UUID: inst.GetID(),
			Name: n,
		})
		inst.firstLoop = true
		pprint.Print(pnt)
	} else {
	}

}

func (inst *Point) Process() {
	if !inst.firstLoop {
		inst.setConnection()

	}
}
