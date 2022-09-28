package flow

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/driver"
	"github.com/NubeDev/flow-eng/services/clients/ffclient"
)

type Network struct {
	*node.Spec
	network driver.Networks
}

var network driver.Networks

// user will add a new connection from maybe the UI
// user will add the network node and then select the connection by name/uuid
// the network node will be a container node, so once they add the network they can then add the device then point

func NewNetwork(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowNetwork, category)
	networkName := node.BuildInput(node.Topic, node.TypeString, nil, body.Inputs)
	value := node.BuildInput(node.In, node.TypeString, nil, body.Inputs)
	inputs := node.BuildInputs(networkName, value)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	network = driver.New(&driver.Network{
		Name:        string(names.FlowFramework),
		Application: names.FlowFramework,
		Storage:     body.GetDB(),
	})

	if body.GetDB() != nil {
		connection, err := body.GetDB().GetConnectionByName("flow-framework")
		fmt.Println(connection, err)
	} else {
		fmt.Println("NO DB BODY")
	}

	fmt.Println(4444, body.ReadPinAsString(node.Topic))

	return &Network{body, network}, nil
}

func (inst *Network) getInst() *Network {
	return inst
}

//type netDetails struct {
//	ip    string
//	port  int
//	token string
//}

func (inst *Network) connection() (*ffclient.Client, error) {
	connection, err := inst.GetDB().GetConnectionByName("flow")
	if err != nil {
		return nil, err
	}

	return ffclient.New(&ffclient.Connection{
		Ip:   connection.Host,
		Port: connection.Port,
	}), nil
}

func (inst *Network) Process() {

}

func (inst *Network) Cleanup() {}
