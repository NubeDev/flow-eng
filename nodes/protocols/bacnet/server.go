package bacnetio

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	log "github.com/sirupsen/logrus"
	"sort"
	"strings"
)

type Bacnet struct {
	Store       *points.Store
	MqttClient  *mqttclient.Client
	Application names.ApplicationName
	Ip          string `json:"ip"`
}

type Server struct {
	*node.Spec
	clients                *clients
	pingFailed             bool
	pingLock               bool
	reconnectedOk          bool
	store                  *points.Store
	application            names.ApplicationName
	loopCount              uint64
	firstMessageFromBacnet bool
	deviceCount            string
	pollingCount           float64
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
	outputMsg := node.BuildOutput(node.Msg, node.TypeString, nil, body.Outputs)
	pollingCount := node.BuildOutput(node.PollingCount, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(outputMsg, pollingCount)
	body.IsParent = true
	body = node.BuildNode(body, nil, outputs, body.Settings)
	clients := &clients{}
	server := &Server{body, clients, false, false, false, opts.Store, application, 0, false, "", 0}
	server.clients.mqttClient = opts.MqttClient
	body.SetSchema(BuildSchemaServer())
	if application == names.Modbus {
		log.Infof("bacnet-server: start application: %s device-ip: %s", application, opts.Ip)
	}
	return server, err
}

func (inst *Server) Process() {
	loopCount, firstLoop := inst.Loop()
	if firstLoop {
		go inst.setMessage()
		go inst.subscribe()
		go inst.mqttReconnect()
		p, ok := inst.getPoints()
		if ok {
			for _, point := range p {
				go inst.mqttPublishNames(point)
			}
		}

	}
	if inst.pingFailed || inst.reconnectedOk { // on failed resubscribe
	}
	if !inst.pingLock {
	}
	if !runnersLock {
		go inst.protocolRunner()
		if loopCount == 1 {
			runnersLock = true
		}
	}
	inst.loopCount = loopCount
	if loopCount%100 == 0 {
		p, ok := inst.getPoints()
		if ok {
			for _, point := range p {
				go inst.mqttPublishNames(point)
			}
		}
	}
	inst.WritePinFloat(node.PollingCount, inst.pollingCount)
}

func setUUID(parentID string, objType points.ObjectType, id points.ObjectID) string {
	return fmt.Sprintf("%s:%s:%d", parentID, objType, id)
}

func (inst *Server) getPV(objType points.ObjectType, id points.ObjectID) (float64, error) {
	pnt, ok := inst.getPoint(objType, id)
	if ok {
		return pnt.PresentValue, nil

	}
	return 0, nil
}

func (inst *Server) writePV(objType points.ObjectType, id points.ObjectID, value float64) error {
	pnt, ok := inst.getPoint(objType, id)
	if ok {
		pnt.PresentValue = value
		err := inst.updatePoint(objType, id, pnt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *Server) updatePoint(objType points.ObjectType, id points.ObjectID, point *points.Point) error {
	s := inst.GetStore()
	if s == nil {
		return nil
	}
	s.Set(setUUID(inst.GetID(), objType, id), point, 0)
	return nil
}

// updateFromBACnet this is a value that has come from the bacnet-server over MQTT
func (inst *Server) updateFromBACnet(objType points.ObjectType, id points.ObjectID, array *points.PriArray) error {
	p, _ := inst.getPoint(objType, id)
	if p != nil {
		p.WriteValueFromBACnet = array
		p.PendingWriteValueFromBACnet = true
		inst.updatePoint(objType, id, p)
	}
	return nil
}

func (inst *Server) getPointsReadOnly() ([]*points.Point, bool) {
	p, _ := inst.getPoints()
	var pointsList []*points.Point
	for _, point := range p {
		if !point.IsWriteable {
			pointsList = append(pointsList, point)
		}
	}
	return pointsList, false
}

func (inst *Server) getPointsWriteOnly() ([]*points.Point, bool) {
	p, _ := inst.getPoints()
	var pointsList []*points.Point
	for _, point := range p {
		if point.IsWriteable {
			pointsList = append(pointsList, point)
		}
	}
	return pointsList, false
}

func (inst *Server) getPoints() ([]*points.Point, bool) {
	s := inst.GetStore()
	if s == nil {
		return nil, false
	}
	var pointsList []*points.Point
	for id, item := range s.All() {
		parts := strings.Split(id, ":")
		if len(parts) > 0 {
			if parts[0] == inst.GetID() {
				point, ok := item.Object.(*points.Point)
				if ok {
					pointsList = append(pointsList, point)
				}
			}
		}
	}
	return pointsList, true
}

func (inst *Server) getPoint(objType points.ObjectType, id points.ObjectID) (*points.Point, bool) {
	s := inst.GetStore()
	if s == nil {
		return nil, false
	}
	d, ok := s.Get(setUUID(inst.GetID(), objType, id))
	if ok {
		return d.(*points.Point), true
	}
	return nil, false
}

func (inst *Server) GetModbusWriteablePoints() *points.ModbusPoints {
	out := &points.ModbusPoints{
		DeviceOne:   []*points.Point{},
		DeviceTwo:   []*points.Point{},
		DeviceThree: []*points.Point{},
		DeviceFour:  []*points.Point{},
	}
	p, _ := inst.getPoints()
	for _, point := range p {
		if point.ModbusDevAddr == 1 {
			if point.IsWriteable {
				out.DeviceOne = append(out.DeviceOne, point)
				sort.Slice(out.DeviceOne[:], func(i, j int) bool { // sort by the modbus register
					return out.DeviceOne[i].ModbusRegister < out.DeviceOne[j].ModbusRegister
				})
			}
		}
		if point.ModbusDevAddr == 2 {
			if point.IsWriteable {
				out.DeviceTwo = append(out.DeviceTwo, point)
				sort.Slice(out.DeviceTwo[:], func(i, j int) bool { // sort by the modbus register
					return out.DeviceTwo[i].ModbusRegister < out.DeviceTwo[j].ModbusRegister
				})
			}
		}
		if point.ModbusDevAddr == 3 {
			if point.IsWriteable {
				out.DeviceThree = append(out.DeviceThree, point)
				sort.Slice(out.DeviceThree[:], func(i, j int) bool { // sort by the modbus register
					return out.DeviceThree[i].ModbusRegister < out.DeviceThree[j].ModbusRegister
				})
			}
		}
		if point.ModbusDevAddr == 4 {
			if point.IsWriteable {
				out.DeviceFour = append(out.DeviceFour, point)
				sort.Slice(out.DeviceFour[:], func(i, j int) bool { // sort by the modbus register
					return out.DeviceFour[i].ModbusRegister < out.DeviceFour[j].ModbusRegister
				})
			}
		}
	}
	return out
}

func (inst *Server) setMessage() {
	schema, err := GetBacnetSchema(inst.GetSettings())
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("start modbus polling on port: %s", schema.Serial)
	inst.WritePin(node.Msg, fmt.Sprintf("port:%s & %s:IO16s", schema.Serial, schema.DeviceCount))
}
