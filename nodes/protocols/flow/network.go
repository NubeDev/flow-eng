package flow

import (
	"fmt"

	"github.com/NubeDev/flow-eng/connections"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	"github.com/enescakir/emoji"
	log "github.com/sirupsen/logrus"
)

type Network struct {
	*node.Spec
	firstLoop                    bool
	loopCount                    uint64
	connection                   *connections.Connection
	mqttClient                   *mqttclient.Client
	mqttConnected                bool
	points                       []*point
	pointsCount                  int
	errorCode                    errorCode
	error                        bool
	fetchPointResponseCount      int
	subscribeFailedPoints        bool
	subscribeFailedPointsList    bool
	subscribeFailedSchedulesList bool
	subscribeFailedMissingPoints bool
}

var mqttQOS = mqttclient.AtMostOnce
var mqttRetain = false

func NewNetwork(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, flowNetwork, Category)
	inputs := node.BuildInputs()
	outputs := node.BuildOutputs(node.BuildOutput(node.Connected, node.TypeBool, nil, body.Outputs))
	body.IsParent = true
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	network := &Network{body, false, 0, nil, nil, false, nil, 0, "", false, 0, false, false, false, false}
	return network, nil
}

func (inst *Network) setConnection() {
	settings, err := getSettings(inst.GetSettings())
	if err != nil {
		errMes := fmt.Sprintf("flow-network, add mqtt broker failed to get settings err:%s", err.Error())
		log.Errorf(errMes)
		inst.setError(errMes, false, false)
		return
	}
	var connection *connections.Connection
	var connectionName = "flow framework integration over MQTT (dont edit/delete)" // this name is set in rubix-edge-wires

	connection, err = inst.Connections().GetConnection(settings.Conn)
	if err != nil {
		errMes := fmt.Sprintf("flow-network error in getting connection: %+v. err:%s", settings.Conn, err.Error())
		log.Errorf(errMes)
		inst.setError(errMes, false, false)
	}

	if connection == nil {
		connection, err = inst.Connections().GetConnectionByName(connectionName)
		if err != nil {
			errMes := fmt.Sprintf("flow-network error in getting connection name: %+v. err:%s", connectionName, err.Error())
			log.Errorf(errMes)
			inst.setError(errMes, false, false)
		}
	}

	if connection == nil {
		errMes := fmt.Sprintf("no flow-network mqtt connection, please select a connection")
		log.Errorf(errMes)
		inst.setError(errMes, false, false)
		return
	}
	inst.connection = connection
	mqttClient, err := mqttclient.NewClient(mqttclient.ClientOptions{
		Servers: []string{fmt.Sprintf("tcp://%s:%d", connection.Host, connection.Port)},
	})
	log.Infof("flow-network mqtt connect try and connect to broker %s ", fmt.Sprintf("tcp://%s:%d", connection.Host, connection.Port))
	err = mqttClient.Connect()
	if err != nil {
		errMes := fmt.Sprintf("flow-network mqtt connect err: %s", err.Error())
		log.Error(errMes)
		inst.setError(errMes, false, false)
	} else {
		inst.mqttClient = mqttClient
		inst.mqttConnected = true
		inst.setError("", true, true)
	}
	inst.setSubTitle(fmt.Sprintf("%s:%d", connection.Host, connection.Port))

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
	if loopCount == 3 {
		go inst.subscribeToEachPoint()
		go inst.subscribeToMissingPoints()
		go inst.schedulesList()
		go inst.publish(loopCount)
	}
	retry := false
	if loopCount == 4 {
		go inst.fetchAllPointValues()
		retry = true
	}
	if loopCount > 5 && loopCount%retryCount == 0 {
		retry = true
	}

	if retry {
		if !inst.mqttConnected {
			log.Errorf("flow-network: reset mqtt connection as first time failed loop-count: %d, inst.mqttConnected: %v", loopCount, inst.mqttConnected)
			go inst.setConnection()
		}
		if inst.subscribeFailedPoints || !inst.mqttConnected {
			log.Errorf("flow-network: reset subscribe to each point as first time failed: %d, inst.subscribeFailedPoints: %v, inst.mqttConnected: %v", loopCount, inst.subscribeFailedPoints, inst.mqttConnected)
			go inst.subscribeToEachPoint()
		}
		if inst.subscribeFailedMissingPoints || !inst.mqttConnected {
			log.Errorf("flow-network: reset subscribe to missing points as first time failed: %d, inst.subscribeFailedMissingPoints: %v, inst.mqttConnected: %v", loopCount, inst.subscribeFailedMissingPoints, inst.mqttConnected)
			go inst.subscribeToMissingPoints()
		}
	}

	if inst.mqttConnected && !inst.error {
		inst.WritePinTrue(node.Connected)
	} else {
		inst.WritePinFalse(node.Connected)
	}
	if loopCount > 2 {
		go inst.publish(loopCount)
	}
	if retry { // get the points every 50 loops
		inst.connectionError()
		inst.fetchExistingPointValues() // refresh point COV
	}

	if loopCount%20 == 0 {
		inst.fetchSchedulesList()
	}

}

func (inst *Network) setSubTitle(msg string) {
	inst.SetSubTitle(msg)
}

func (inst *Network) setError(msg string, reset, setMQTTConnected bool) {
	if reset {
		inst.SetStatusError("error cleared")
		inst.SetErrorIcon(string(emoji.GreenCircle))
	} else {
		inst.SetStatusError(msg)
		inst.SetErrorIcon(string(emoji.RedCircle))
		inst.mqttConnected = setMQTTConnected
	}
}

func (inst *Network) connectionError() {
	if inst.error {
		inst.SetStatusError(string(inst.errorCode))
		inst.SetErrorIcon(string(emoji.RedCircle))
		log.Error(inst.errorCode)
		inst.mqttConnected = false
	} else {
		inst.setError("", true, false)
	}
}
