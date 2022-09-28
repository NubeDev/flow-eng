package flow

import (
	"fmt"
	"github.com/NubeDev/flow-eng/db"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/driver"
	"github.com/NubeDev/flow-eng/services/clients/ffclient"
)

type Network struct {
	*node.Spec
	firstLoop      bool
	loopCount      uint64
	networkUUID    string
	connectionUUID string
	connection     *db.Connection
	client         *ffclient.Client
	pool           driver.Driver
}

// user will select an existing connection
// user will select a ff network/device and then point by names

func NewNetwork(body *node.Spec, pool driver.Driver) (node.Node, error) {
	//var err error
	body = node.Defaults(body, flowNetwork, category)
	connectionName := node.BuildInput(node.Connection, node.TypeString, nil, body.Inputs)
	name := node.BuildInput(node.Name, node.TypeString, nil, body.Inputs)
	networkUUID := node.BuildInput(node.UUID, node.TypeString, nil, body.Inputs)
	inputs := node.BuildInputs(connectionName, name, networkUUID)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Network{body, false, 0, body.ReadPinAsString(node.UUID), "", nil, nil, pool}, nil
}

func (inst *Network) getNetwork() driver.Driver {
	return inst.pool
}

func (inst *Network) getInst() *Network {
	return inst
}

func (inst *Network) setConnection() {
	fmt.Println("ADD NETWORK", inst.pool)
	connection, err := inst.GetDB().GetConnection("con_1b8b9c8bd63f")
	if err != nil {
		inst.firstLoop = false // if fail try again
		return
	}
	inst.connectionUUID = connection.UUID
	inst.client = ffclient.New(&ffclient.Connection{
		Ip:   connection.Host,
		Port: connection.Port,
	})

	net := inst.pool.AddNetwork(&driver.Network{
		UUID:           inst.networkUUID,
		Name:           inst.ReadPinAsString(node.Topic),
		ConnectionUUID: connection.UUID,
	})

	pprint.PrintJOSN(net)
	inst.firstLoop = true

}

func (inst *Network) ping(loop uint64) {
	rePing := loop % 10
	if rePing == 0 {
		err := inst.client.Ping()
		if err != nil {

		}
	}

}

func (inst *Network) Process() {
	inst.loopCount++
	if !inst.firstLoop {
		inst.setConnection()

	}
	inst.ping(inst.loopCount)

	inst.runner()

}

func (inst *Network) Cleanup() {}
