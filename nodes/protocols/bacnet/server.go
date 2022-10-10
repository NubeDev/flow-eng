package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	edge28lib "github.com/NubeDev/flow-eng/services/edge28"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	rubixIO "github.com/NubeDev/flow-eng/services/rubixio"
	log "github.com/sirupsen/logrus"
)

type Bacnet struct {
	Store       *points.Store
	MqttClient  *mqttclient.Client
	Application names.ApplicationName
}

type Server struct {
	*node.Spec
	clients       *clients
	pingFailed    bool
	pingLock      bool
	runnersLock   bool
	reconnectedOk bool
	store         *points.Store
	application   names.ApplicationName
}

type clients struct {
	mqttClient *mqttclient.Client
	rio        *rubixIO.RubixIO
	edge28     *edge28lib.Edge28
}

func bacnetOpts(opts *Bacnet) *Bacnet {
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
	ip := "0.0.0.0"
	body = node.Defaults(body, serverNode, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputApplication := node.BuildOutput(node.Msg, node.TypeString, nil, body.Outputs)
	outputErr := node.BuildOutput(node.ErrMsg, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(outputApplication, outputErr)
	parameters := &node.Parameters{
		Application: &node.Application{
			Application: names.BACnet,
			IsChild:     false,
		},
		MaxNodeCount: 1,
	}
	body.Parameters = node.BuildParameters(parameters) // if node is already added then show the user
	body.IsParent = true
	body = node.BuildNode(body, inputs, outputs, nil)
	clients := &clients{}
	server := &Server{body, clients, false, false, false, false, opts.Store, application}
	server.clients.mqttClient = opts.MqttClient
	if application == names.RubixIO || application == names.RubixIOAndModbus {
		rubixIOUICount, rubixIOUOCount := points.CalcPointCount(1, application)
		rio := &rubixIO.RubixIO{}
		rio = rubixIO.New(&rubixIO.RubixIO{
			IP:          ip,
			StartAddrUI: rubixIOUICount,
			StartAddrUO: rubixIOUOCount,
			StartAddrDO: 2,
		})
		server.clients.rio = rio
		log.Infof("bacnet-server: start application: %s device-ip: %s", application, ip)
	}
	if application == names.Edge {
		edge28 := edge28lib.New(ip)
		server.clients.edge28 = edge28
		log.Infof("bacnet-server: start application: %s device-ip: %s", application, ip)
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
	if !inst.runnersLock {
		go inst.protocolRunner(inst.application)
		inst.runnersLock = true
	}
}

func (inst *Server) Cleanup() {}
