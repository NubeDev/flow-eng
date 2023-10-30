package bacnetio

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/mqttclient"
	"github.com/dustin/go-humanize"
	"github.com/enescakir/emoji"
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
	clients                *clients
	pingFailed             bool
	pingLock               bool
	reconnectedOk          bool
	store                  *points.Store
	application            names.ApplicationName
	loopCount              uint64
	firstMessageFromBacnet bool
	deviceCount            string
	deviceCountNumber      int
	deviceSelectionCount   string
	pollingCount           int64
	finishModbusLoop       bool
	devStats1              string
	devStats2              string
	devStats3              string
	devStats4              string
	devAddr1               int
	devAddr2               int
	devAddr3               int
	devAddr4               int
	isImperial             bool
	bacnetServerDeviceId   int
	bacnetServerPort       int
	bacnetServerIP         string
}

var runnersLock bool

const startProtocolRunnerCount = 10

type clients struct {
	mqttClient *mqttclient.Client
}

func bacnetOpts(opts ...any) *Bacnet {
	var bn *Bacnet
	if len(opts) == 1 {
		bn = opts[0].(*Bacnet)
	} else {
		bn = &Bacnet{}
	}
	return bn
}

var mqttQOS = mqttclient.AtMostOnce
var mqttRetain = false

const (
	devStats1 = "dev-1-stats"
	devStats2 = "dev-2-stats"
	devStats3 = "dev-3-stats"
	devStats4 = "dev-4-stats"
)

func NewServer(body *node.Spec, opts ...any) (node.Node, error) {
	bn := bacnetOpts(opts...)
	var application = bn.Application
	var err error
	body = node.Defaults(body, serverNode, Category)
	deviceError1 := node.BuildOutput(devStats1, node.TypeString, nil, body.Outputs)
	deviceError2 := node.BuildOutput(devStats2, node.TypeString, nil, body.Outputs)
	deviceError3 := node.BuildOutput(devStats3, node.TypeString, nil, body.Outputs)
	deviceError4 := node.BuildOutput(devStats4, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(deviceError1, deviceError2, deviceError3, deviceError4)
	body.IsParent = true
	body = node.BuildNode(body, nil, outputs, body.Settings)
	clients := &clients{}
	server := &Server{
		body,
		clients,
		false,
		false,
		false,
		bn.Store,
		application,
		0,
		false,
		"",
		0,
		"",
		0,
		false,
		"",
		"",
		"",
		"",
		1,
		2,
		3,
		4,
		false,
		0,
		0,
		"",
	}
	server.clients.mqttClient = bn.MqttClient
	body.SetSchema(BuildSchemaServer())
	if application == names.Modbus {
		log.Debugf("bacnet-server: start application: %s device-ip: %s", application, bn.Ip)
	}
	return server, err
}

func (inst *Server) Process() {

	loopCount, firstLoop := inst.Loop()
	if firstLoop {
		inst.setDeviceCount()
		err := inst.setDeviceAddress()
		if err != nil {
			log.Error(err)
		}
		go inst.subscribe()
		go inst.mqttReconnect()
		inst.deviceSelectionCount = inst.getDeviceCount()

		config, err := inst.GetRubixOSConfig()
		if err == nil {
			if config != nil {
				inst.isImperial = config.UnitsImperial
			}
		}
		port, deviceID, ip, err := inst.getServerConfig()
		if err != nil {
			log.Errorf("bacnet-units get config: %s", err.Error())
		} else {
			inst.bacnetServerDeviceId = deviceID
			inst.bacnetServerPort = port
			inst.bacnetServerIP = ip
		}

	}
	if loopCount == 3 { // publish all the point names
		p, ok := inst.getPoints()
		if ok {
			for _, point := range p {
				go inst.mqttPublishNames(point)
			}
		}
	}
	if !runnersLock {
		if loopCount == startProtocolRunnerCount {
			go inst.protocolRunner()
			runnersLock = true
		}
	}
	inst.loopCount = loopCount
	if loopCount%200 == 0 {
		go inst.mqttReconnect()
		p, ok := inst.getPoints()
		if ok {
			for _, point := range p {
				go inst.mqttPublishNames(point)
			}
		}
	}
	message := fmt.Sprintf("poll count: %s", pollStats(inst.pollingCount))
	inst.SetSubTitle(message)
	inst.WritePin(devStats1, deviceMsg(inst.devStats1, inst.devAddr1))
	inst.WritePin(devStats2, deviceMsg(inst.devStats2, inst.devAddr2))
	inst.WritePin(devStats3, deviceMsg(inst.devStats3, inst.devAddr3))
	inst.WritePin(devStats4, deviceMsg(inst.devStats4, inst.devAddr4))

}

func deviceMsg(message string, addr int) string {
	return fmt.Sprintf("addr: %d %s", addr, message)
}

func pollStats(pollingCount int64) string {
	return humanize.Comma(pollingCount)
}

func setUUID(parentID string, objType points.ObjectType, id points.ObjectID) string {
	return fmt.Sprintf("%s:%s:%d", parentID, objType, id)
}

func (inst *Server) getPV(objType points.ObjectType, id points.ObjectID) (*float64, error) {
	pnt, ok := inst.getPoint(objType, id)
	if ok {
		return pnt.PresentValue, nil

	}
	return float.New(0), nil
}

func (inst *Server) setWaringMessage(message string, clear bool) {
	if clear {
		inst.SetStatusMessage("")
	}
	inst.SetWaringMessage(message)
	inst.SetWaringIcon(string(emoji.OrangeCircle))
}

func (inst *Server) setStatusMessage(message string, clear bool) {
	if clear {
		inst.SetStatusMessage("")
	}
	inst.SetStatusMessage(message)
	inst.SetNotifyIcon(string(emoji.GreenCircle))
}

func (inst *Server) setStatusError(message string, clear bool) {
	if clear {
		inst.SetStatusError("")
	} else {
		inst.SetStatusError(message)
		inst.SetErrorIcon(string(emoji.RedCircle))
	}

}

func (inst *Server) writePV(objType points.ObjectType, id points.ObjectID, value float64) error {
	pnt, ok := inst.getPoint(objType, id)
	if ok {
		newPV := conversions.ValueTransformOnRead(value, pnt.ScaleEnable, pnt.Factor, pnt.ScaleInMin, pnt.ScaleInMax, pnt.ScaleOutMin, pnt.ScaleOutMax, pnt.Offset)
		if pnt.PresentValue == nil || newPV != *pnt.PresentValue {
			pnt.PendingMQTTPublish = true
			pnt.ModbusUpdated = true
			pnt.PresentValue = float.New(newPV)
			err := inst.updatePoint(objType, id, pnt)
			if err != nil {
				return err
			}
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
		err := inst.updatePoint(objType, id, p)
		valuePri := points.GetHighest(array)
		if valuePri != nil {
			log.Debugf("bacnet write mqtt value to objType: %s -> objId: %d value: %f pri: %d", objType, id, valuePri.Value, valuePri.Number)
		}
		if err != nil {
			return err
		}
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
				if ok && point != nil {
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
		if point.DeviceNumber == 1 {
			if point.IsWriteable {
				out.DeviceOne = append(out.DeviceOne, point)
				sort.Slice(out.DeviceOne[:], func(i, j int) bool { // sort by the modbus register
					return out.DeviceOne[i].ModbusRegister < out.DeviceOne[j].ModbusRegister
				})
			}
		}
		if point.DeviceNumber == 2 {
			if point.IsWriteable {
				out.DeviceTwo = append(out.DeviceTwo, point)
				sort.Slice(out.DeviceTwo[:], func(i, j int) bool { // sort by the modbus register
					return out.DeviceTwo[i].ModbusRegister < out.DeviceTwo[j].ModbusRegister
				})
			}
		}
		if point.DeviceNumber == 3 {
			if point.IsWriteable {
				out.DeviceThree = append(out.DeviceThree, point)
				sort.Slice(out.DeviceThree[:], func(i, j int) bool { // sort by the modbus register
					return out.DeviceThree[i].ModbusRegister < out.DeviceThree[j].ModbusRegister
				})
			}
		}
		if point.DeviceNumber == 4 {
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

func (inst *Server) setDeviceAddress() error {
	schema, err := inst.getSchema()
	if err != nil || schema == nil {
		inst.devAddr1 = 1
		inst.devAddr2 = 2
		inst.devAddr3 = 3
		inst.devAddr4 = 4
		return errors.New("failed to get bacnet-server settings")
	}
	if schema.DeviceAddress1 <= 0 {
		inst.devAddr1 = 1
	} else {
		inst.devAddr1 = schema.DeviceAddress1
	}
	if schema.DeviceAddress2 <= 0 {
		inst.devAddr2 = 2
	} else {
		inst.devAddr2 = schema.DeviceAddress2
	}
	if schema.DeviceAddress3 <= 0 {
		inst.devAddr3 = 3
	} else {
		inst.devAddr3 = schema.DeviceAddress3
	}
	if schema.DeviceAddress4 <= 0 {
		inst.devAddr4 = 4
	} else {
		inst.devAddr4 = schema.DeviceAddress4
	}
	return nil
}

func (inst *Server) getSchema() (*BacnetSchema, error) {
	return GetBacnetSchema(inst.GetSettings())
}

func (inst *Server) setDevStats1(msg string) {
	if inst.deviceCountNumber < 1 {
		inst.devStats1 = "not added"
	} else {
		inst.devStats1 = msg
	}
}

func (inst *Server) setDevStats2(msg string) {
	if inst.deviceCountNumber < 2 {
		inst.devStats2 = "not added"
	} else {
		inst.devStats2 = msg
	}
}

func (inst *Server) setDevStats3(msg string) {
	if inst.deviceCountNumber < 3 {
		inst.devStats3 = "not added"
	} else {
		inst.devStats3 = msg
	}
}

func (inst *Server) setDevStats4(msg string) {
	if inst.deviceCountNumber < 4 {
		inst.devStats4 = "not added"
	} else {
		inst.devStats4 = msg
	}
}

func (inst *Server) getDeviceCount() string {
	schema, err := inst.getSchema()
	if err != nil {
		return ""
	}
	if schema != nil {
		return schema.DeviceCount
	}
	return ""
}

func (inst *Server) setDeviceCount() {
	var count int
	schema, err := inst.getSchema()
	if err != nil {
		return
	}
	if schema != nil {
		inst.deviceCount = schema.DeviceCount
	}
	deviceCount := inst.deviceCount
	if strings.Contains(deviceCount, "1x") {
		count = 1
	}
	if strings.Contains(deviceCount, "2x") {
		count = 2
	}
	if strings.Contains(deviceCount, "3x") {
		count = 3
	}
	if strings.Contains(deviceCount, "4x") {
		count = 4
	}
	inst.deviceCountNumber = count
}
