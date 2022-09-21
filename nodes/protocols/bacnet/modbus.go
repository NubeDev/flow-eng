package bacnet

import (
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	points2 "github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	modbuscli2 "github.com/NubeDev/flow-eng/services/modbuscli"
	log "github.com/sirupsen/logrus"
	"time"
)

// modbus will come from polling
// this is to only work for the IO-16
func (inst *Server) modbusRunner() {
	log.Info("start modbus-runner")
	cli := &modbuscli2.Modbus{
		IsSerial: false,
		Address:  "192.168.15.202",
		Port:     502,
		Slave:    1,
	}
	init, err := cli.Init(cli)
	if err != nil {
		return
	}
	store := getStore()
	pointsList := store.GetPointsByApplication(applications.Modbus)
	inst.modbusInputsRunner(init, pointsList) // process the inputs
	time.Sleep(4 * time.Second)
	modbusLoop = false

}

func (inst *Server) modbusInputsRunner(cli *modbuscli2.Modbus, pointsList []*points2.Point) {
	var err error
	var tempList []float64
	var voltList []float64
	var completedTemp bool
	var completedVolt bool
	store := getStore()
	for _, point := range pointsList { // do modbus read
		if !point.IsWriteable {
			addr, _ := cli.BuildInput(point.IoType, point.ObjectID)
			slaveId := addr.DeviceAddr
			if !completedTemp && point.IoType == points2.IoTypeTemp {
				tempList, err = cli.ReadTemps(slaveId) // DO MODBUS READ FOR TEMPS
				if err != nil {
					//return
				}
				completedTemp = true
			}
			if !completedVolt && point.IoType == points2.IoTypeVolts {
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
		if point.ObjectType == points2.AnalogInput {
			if point.ObjectID == points2.ObjectID(objectId) {
				p := store.GetPointByObject(points2.AnalogInput, point.ObjectID)
				io16Pin := addr.IoPin
				if io16Pin <= 0 {
					log.Errorf("modbus-polling failed to get correct io-pin")
					continue
				}
				if point.IoType == points2.IoTypeTemp || point.IoType == points2.IoTypeDigital { // update anypoint that is type temp
					if point.IoType == points2.IoTypeDigital {
						writeValue = modbuscli2.TempToDI(tempList[io16Pin]) // covert them temp value to a DI value
					} else {
						writeValue = tempList[io16Pin]
					}
				}
				if point.IoType == points2.IoTypeVolts { // update anypoint that is type voltage
					writeValue = voltList[io16Pin]
				}
				store.WritePointValue(p.UUID, writeValue)
			}
		}
	}

}
