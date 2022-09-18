package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/modbuscli"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bstore"
)

// processMessage are the messages from the bacnet-server via the mqtt-broker

func (inst *Server) mqttRunner() {
	go func() {
		msg, ok := inst.bus().Recv()
		if ok {
			payload := bstore.NewPayload()
			err := payload.NewMessage(msg)
			if err != nil {
				return
			}
			t, id := payload.GetObject()
			pnt := inst.db().GetPointByObject(t, id)
			inst.db().WritePointValue(pnt.UUID, float.NonNil(payload.GetPresentValue()))
		}
	}()
}

// edge28 & rubix-io input values will come from rest
func (inst *Server) edgeRunner() {

}

// modbus will come from polling
// this is to only work for the IO-16
func (inst *Server) modbusRunner() {
	cli := modbuscli.New(&modbuscli.Modbus{})
	store := inst.db()
	pointsList := store.GetPointsByApplication(applications.Modbus)
	inst.modbusInputsRunner(cli, pointsList) // process the inputs
}

func (inst *Server) modbusInputsRunner(cli *modbuscli.Modbus, pointsList []*bstore.Point) {
	var err error
	var tempList []float64
	var voltList []float64
	var completedTemp bool
	var completedVolt bool
	store := inst.db()
	for _, point := range pointsList { // do modbus read
		if !point.IsWriteable {
			addr, _ := cli.BuildInput(point.IoType, point.ObjectID)
			slaveId := addr.DeviceAddr
			if !completedTemp && point.IoType == bstore.IoTypeTemp {
				tempList, err = cli.ReadTemps(slaveId)
				if err != nil {
					return
				}
				completedTemp = true
			}
			if !completedVolt && point.IoType == bstore.IoTypeVolts {
				tempList, err = cli.ReadVolts(slaveId)
				if err != nil {
					return
				}
				completedVolt = true
			}
			if err != nil {
				continue
			}
		}
	}
	for _, point := range pointsList {
		addr, _ := cli.BuildInput(point.IoType, point.ObjectID)
		objectId := addr.BacnetAddr
		var writeValue float64
		if point.ObjectType == bstore.AnalogInput {
			if point.ObjectID == bstore.ObjectID(objectId) {
				p := store.GetPointByObject(bstore.AnalogInput, point.ObjectID)
				io16Pin := addr.IoPin - 1
				if point.IoType == bstore.IoTypeTemp || point.IoType == bstore.IoTypeDigital { // update anypoint that is type temp
					if point.IoType == bstore.IoTypeDigital {
						writeValue = modbuscli.TempToDI(tempList[io16Pin]) // covert them temp value to a DI value
					} else {
						writeValue = tempList[io16Pin]
					}
				}
				if point.IoType == bstore.IoTypeVolts { // update anypoint that is type voltage
					writeValue = voltList[io16Pin]
				}
				store.WritePointValue(p.UUID, writeValue)
			}
		}
	}

}
