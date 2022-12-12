package bacnetio

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/modbuscli"
	"github.com/NubeIO/nubeio-rubix-lib-modbus-go/modbus"
	log "github.com/sirupsen/logrus"
	"time"
)

// modbus will come from polling
// this is to only work for the IO-16
func (inst *Server) modbusRunner(settings map[string]interface{}) {
	log.Info("start modbus-runner")
	schema, err := GetBacnetSchema(settings)
	if err != nil {
		log.Error(err)
		return
	}
	port := "/dev/ttyUSB0"
	if schema.Serial != "" {
		port = schema.Serial
	}
	cli := &modbuscli.Modbus{
		IsSerial: true,
		Serial: &modbus.Serial{
			SerialPort: port,
		},
	}
	log.Infof("start modbus polling on port: %s", port)
	init, err := cli.Init(cli)
	if err != nil {
		log.Error(err)
		return
	}
	var count int
	for {
		log.Infof("modbus polling loop count: %d application-type: %s", count, inst.application)
		fmt.Println(names.Modbus)
		pointsList := inst.store.GetPointsByApplication(names.Modbus)
		inst.modbusInputsRunner(init, pointsList) // process the inputs
		inst.modbusOutputsDispatch(init)          // process the outs
		time.Sleep(1 * time.Second)
		count++
	}
}

func modbusBulkWrite(pointsList []*points.Point) [8]float64 {
	var out [8]float64
	for i, point := range pointsList {
		v := points.GetHighest(point.WriteValue)
		if v != nil {
			out[i] = v.Value
		}
	}
	return out
}

func (inst *Server) modbusOutputsDispatch(cli *modbuscli.Modbus) {

	pointsList := inst.store.GetModbusWriteablePoints(true)
	if pointsList == nil {
		//return
	}
	if len(pointsList.DeviceOne) > 0 {
		err := cli.Write(1, modbusBulkWrite(pointsList.DeviceOne))
		if err != nil {
			log.Error(err)
		}
		for _, point := range pointsList.DeviceOne {
			inst.store.CompletePendingWriteCount(point)
		}
	}
	if len(pointsList.DeviceTwo) > 0 {
		err := cli.Write(2, modbusBulkWrite(pointsList.DeviceTwo))
		if err != nil {
			log.Error(err)
		}
		for _, point := range pointsList.DeviceOne {
			inst.store.CompletePendingWriteCount(point)
		}
	}

}

func (inst *Server) modbusInputsRunner(cli *modbuscli.Modbus, pointsList []*points.Point) {
	var err error
	var tempList [8]float64
	var voltList [8]float64
	var completedTemp bool
	var completedVolt bool
	fmt.Println(222, len(pointsList), pointsList[0].IoType)
	pprint.PrintJOSN(pointsList[0])
	for _, point := range pointsList { // do modbus read
		if !point.IsWriteable {
			addr, _ := points.ModbusBuildInput(point.IoType, point.ObjectID)
			pprint.PrintJOSN(addr)
			slaveId := addr.DeviceAddr
			if !completedTemp && (point.IoType == points.IoTypeTemp || point.IoType == points.IoTypeDigital) {
				tempList, err = cli.ReadTemps(slaveId) // DO MODBUS READ FOR TEMPS
				if err != nil {
					log.Errorf("modbus read temp %s", err.Error())
				}
				completedTemp = true
			}
			if !completedVolt && point.IoType == points.IoTypeVolts {
				voltList, err = cli.ReadVolts(slaveId) // DO MODBUS READ FOR VOLTS
				if err != nil {
					log.Errorf("modbus read voltages %s", err.Error())
				}
				completedVolt = true
			}
		}
	}
	for _, point := range pointsList {
		addr, _ := points.ModbusBuildInput(point.IoType, point.ObjectID)
		objectId := addr.BacnetAddr
		var returnedValue float64
		if point.ObjectType == points.AnalogInput {
			if point.ObjectID == points.ObjectID(objectId) {
				p := inst.store.GetPointByObject(points.AnalogInput, point.ObjectID)
				io16Pin := addr.IoPin - 1
				if io16Pin < 0 || io16Pin > len(tempList) {
					log.Errorf("modbus-polling failed to get correct io-pin")
					continue
				}
				if point.IoType == points.IoTypeTemp || point.IoType == points.IoTypeDigital { // update anypoint that is type temp
					if point.IoType == points.IoTypeDigital {
						returnedValue = modbuscli.TempToDI(tempList[io16Pin]) // covert them temp value to a DI value
					} else {
						returnedValue = tempList[io16Pin]
					}
				}
				if point.IoType == points.IoTypeVolts { // update anypoint that is type voltage
					returnedValue = voltList[io16Pin]
				}
				inst.store.WriteValueFromRead(p, returnedValue)
			}
		}
	}

}
