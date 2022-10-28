package flow

import (
	"fmt"
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	log "github.com/sirupsen/logrus"
)

type Network struct {
	*node.Spec
	firstLoop     bool
	loopCount     uint64
	connection    *db.Connection
	mqttClient    *mqttclient.Client
	mqttConnected bool
	points        []*point
}

var mqttQOS = mqttclient.AtMostOnce
var mqttRetain = false

func NewNetwork(body *node.Spec) (node.Node, error) {
	//var err error
	body = node.Defaults(body, flowNetwork, category)
	inputs := node.BuildInputs()
	outputs := node.BuildOutputs(node.BuildOutput(node.Connected, node.TypeBool, nil, body.Outputs))
	body.IsParent = true
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	network := &Network{body, false, 0, nil, nil, false, nil}
	return network, nil
}

func (inst *Network) setConnection() {
	settings, err := getSettings(inst.GetSettings())
	if err != nil {
		errMes := fmt.Sprintf("flow-network add mqtt broker failed to get settings err:%s", err.Error())
		log.Errorf(errMes)
		inst.SetStatusError(errMes)
		return
	}
	if settings == nil {
		return
	}
	connection, err := inst.GetDB().GetConnection(settings.Conn)
	if err != nil {
		errMes := fmt.Sprintf("flow-network error in getting connection err:%s", err.Error())
		log.Errorf(errMes)
		inst.SetStatusError(errMes)
		return
	}
	if connection == nil {
		errMes := fmt.Sprintf("no flow-network mqtt connection, please select a connection")
		log.Errorf(errMes)
		inst.SetStatusError(errMes)
		return
	}
	inst.connection = connection
	mqttClient, err := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{fmt.Sprintf("tcp://%s:%d", connection.Host, connection.Port)},
	})
	err = mqttClient.Connect()
	if err != nil {
		errMes := fmt.Sprintf("flow-network mqtt connect err: %s", err.Error())
		log.Error(errMes)
		inst.SetStatusError(errMes)
	} else {
		inst.mqttClient = mqttClient
		inst.mqttConnected = true
	}

}

func (inst *Network) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}

func (inst *Network) Process() {
	loopCount, firstLoop := inst.Loop()
	if firstLoop {
		go inst.setConnection()
	}
	if loopCount == 2 {
		go inst.subscribe()
		go inst.pointsList()
		go inst.publish(loopCount)
	}
	if inst.mqttConnected {
		inst.WritePinTrue(node.Connected)
	} else {
		inst.WritePinFalse(node.Connected)
	}
	if loopCount > 1 {
		go inst.publish(loopCount)
	}
	if loopCount%50 == 0 { // get the points every 50 loops
		inst.fetchPointsList()
	}

}
