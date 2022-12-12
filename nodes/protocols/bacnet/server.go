package bacnetio

import (
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	log "github.com/sirupsen/logrus"
)

type Bacnet struct {
	Store       *points.Store
	MqttClient  *mqttclient.Client
	Application names.ApplicationName
	Ip          string `json:"ip"`
}

type Server struct {
	*node.Spec
	clients       *clients
	pingFailed    bool
	pingLock      bool
	reconnectedOk bool
	store         *points.Store
	application   names.ApplicationName
}

var runnersLock bool

type clients struct {
	mqttClient *mqttclient.Client
}

func bacnetOpts(opts *Bacnet) *Bacnet {
	if opts != nil {
		if opts.Store == nil {
			log.Error("bacnet store can not be empty")
		}
	}
	if opts == nil {
		return &Bacnet{}
	}
	return opts
}

var mqttQOS = mqttclient.AtMostOnce
var mqttRetain = false

func NewServer(body *node.Spec, opts *Bacnet) (node.Node, error) {
	opts = bacnetOpts(opts)
	var application = opts.Application
	var err error
	body = node.Defaults(body, serverNode, category)
	//inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputApplication := node.BuildOutput(node.Msg, node.TypeString, nil, body.Outputs)
	outputErr := node.BuildOutput(node.ErrMsg, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(outputApplication, outputErr)
	body.IsParent = true
	body = node.BuildNode(body, nil, outputs, body.Settings)
	clients := &clients{}
	server := &Server{body, clients, false, false, false, opts.Store, application}
	server.clients.mqttClient = opts.MqttClient
	body.SetSchema(BuildSchemaServer())
	if application == names.Modbus {
		log.Infof("bacnet-server: start application: %s device-ip: %s", application, opts.Ip)
	}
	return server, err
}

func (inst *Server) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		go inst.subscribe()
		go inst.mqttReconnect()
	}
	if inst.pingFailed || inst.reconnectedOk { // on failed resubscribe
	}
	if !inst.pingLock {
	}
	if !runnersLock {
		go inst.protocolRunner()
		runnersLock = true
	}
}
