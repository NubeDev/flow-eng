package broker

import (
	"fmt"
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/driver"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeDev/flow-eng/services/clients/ffclient"
	log "github.com/sirupsen/logrus"
)

type Broker struct {
	*node.Spec
	firstLoop      bool
	loopCount      uint64
	networkUUID    string
	connectionUUID string
	connection     *db.Connection
	client         *ffclient.Client
	pool           driver.Driver
}

func NewBroker(body *node.Spec, pool driver.Driver) (node.Node, error) {
	//var err error
	body = node.Defaults(body, mqttBroker, category)
	connectionName := node.BuildInput(node.Connection, node.TypeString, nil, body.Inputs)
	name := node.BuildInput(node.Name, node.TypeString, nil, body.Inputs)
	inputs := node.BuildInputs(connectionName, name)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body.IsParent = true
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	n, _ := body.ReadPinAsString(node.UUID)
	network := &Broker{body, false, 0, n, "", nil, nil, pool}
	body.SetSchema(network.buildSchema())
	network.Spec = body
	return network, nil
}

func (inst *Broker) getBroker() driver.Driver {
	return inst.pool
}

func (inst *Broker) getInst() *Broker {
	return inst
}

func (inst *Broker) setConnection() {
	fmt.Println("ADD NETWORK", inst.pool)

	connection, err := inst.GetDB().GetConnection("con_db7de7598bba")
	if err != nil {
		inst.firstLoop = false // if fail try again
		log.Error("add flow network failed to find connection")
		return
	}
	inst.connectionUUID = connection.UUID
	inst.client = ffclient.New(&ffclient.Connection{
		Ip:   connection.Host,
		Port: connection.Port,
	})
	//t, _ := inst.ReadPinAsString(node.Topic)
	inst.firstLoop = true

}

func (inst *Broker) ping(loop uint64) {
	rePing := loop % 10
	if rePing == 0 {
		err := inst.client.Ping()
		if err != nil {

		}
	}

}

func (inst *Broker) GetSchema() *schemas.Schema {
	return inst.buildSchema()
}

func (inst *Broker) Process() {
	inst.loopCount++
	if !inst.firstLoop {
		inst.setConnection()
	}
	inst.ping(inst.loopCount)

}
