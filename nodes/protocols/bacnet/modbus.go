package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/modbuscli"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/points"
	log "github.com/sirupsen/logrus"
	"time"
)

// modbus will come from polling
// this is to only work for the IO-16
func (inst *Server) modbusRunner() {
	log.Info("start modbus-runner")
	cli := &modbuscli.Modbus{
		IsSerial: false,
		Address:  "192.168.15.202",
		Port:     502,
		Slave:    1,
	}
	init, err := cli.Init(cli)
	if err != nil {
		return
	}
	store := inst.db()
	pointsList := store.GetPointsByApplication(applications.Modbus)
	inst.modbusInputsRunner(init, pointsList) // process the inputs
	time.Sleep(4 * time.Second)
	modbusLoop = false

}

func (inst *Server) modbusInputsRunner(cli *modbuscli.Modbus, pointsList []*points.Point) {
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
			if !completedTemp && point.IoType == points.IoTypeTemp {
				tempList, err = cli.ReadTemps(slaveId) // DO MODBUS READ FOR TEMPS
				if err != nil {
					//return
				}
				completedTemp = true
			}
			if !completedVolt && point.IoType == points.IoTypeVolts {
				tempList, err = cli.ReadVolts(slaveId) // DO MODBUS READ FOR VOLTS
				if err != nil {
					//return
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
		if point.ObjectType == points.AnalogInput {
			if point.ObjectID == points.ObjectID(objectId) {
				p := store.GetPointByObject(points.AnalogInput, point.ObjectID)
				io16Pin := addr.IoPin
				if io16Pin <= 0 {
					log.Errorf("modbus-polling failed to get correct io-pin")
					continue
				}
				if point.IoType == points.IoTypeTemp || point.IoType == points.IoTypeDigital { // update anypoint that is type temp
					if point.IoType == points.IoTypeDigital {
						writeValue = modbuscli.TempToDI(tempList[io16Pin]) // covert them temp value to a DI value
					} else {
						writeValue = tempList[io16Pin]
					}
				}
				if point.IoType == points.IoTypeVolts { // update anypoint that is type voltage
					writeValue = voltList[io16Pin]
				}
				store.WritePointValue(p.UUID, writeValue)
			}
		}
	}

}
