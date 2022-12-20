package bacnetio

import (
	"fmt"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/modbuscli"
	"github.com/NubeIO/nubeio-rubix-lib-modbus-go/modbus"
	"github.com/enescakir/emoji"
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
		inst.SetStatusError("failed to set node settings")
		inst.SetErrorIcon(string(emoji.RedCircle))
		return
	}
	port := "/dev/ttyAMA0"
	if schema.Serial != "" {
		port = schema.Serial
	}
	cli := &modbuscli.Modbus{
		IsSerial: true,
		Serial: &modbus.Serial{
			SerialPort: port,
		},
	}
	init, err := cli.Init(cli)
	if err != nil {
		inst.SetStatusError(fmt.Sprintf("failed to set serial port: %s", port))
		inst.SetErrorIcon(string(emoji.RedCircle))
		log.Error(err)
		return
	}
	var count float64
	for {
		log.Infof("modbus polling loop count: %f application-type: %s", count, inst.application)
		inst.pollingCount = count
		pointsListRead, _ := inst.getPointsReadOnly()
		inst.modbusInputsRunner(init, pointsListRead) // process the inputs
		inst.modbusOutputsDispatch(init)              // process the outs
		time.Sleep(500 * time.Millisecond)
		count++
	}
}

func modbusBulkWrite(pointsList []*points.Point) [8]float64 {
	var out [8]float64
	for i, point := range pointsList {
		v := points.GetHighest(point.WriteValue)
		if v != nil {
			var value = v.Value
			if point.ObjectType == points.AnalogOutput {
				if value >= 10 {
					value = 10
				}
				if point.IoType == points.IoTypeDigital {
					if value > 0 {
						value = 10
					}
				}
			}
			out[i] = value
		}
	}
	return out
}

func (inst *Server) modbusOutputsDispatch(cli *modbuscli.Modbus) {
	pointsList := inst.GetModbusWriteablePoints()
	if pointsList == nil {

	}
	if len(pointsList.DeviceOne) > 0 {
		err := cli.Write(1, modbusBulkWrite(pointsList.DeviceOne))
		if err != nil {
			log.Error(err)
		}
	}
	if len(pointsList.DeviceTwo) > 0 {
		err := cli.Write(2, modbusBulkWrite(pointsList.DeviceTwo))
		if err != nil {
			log.Error(err)
		}
	}
	if len(pointsList.DeviceThree) > 0 {
		err := cli.Write(3, modbusBulkWrite(pointsList.DeviceThree))
		if err != nil {
			log.Error(err)
		}
	}
	if len(pointsList.DeviceFour) > 0 {
		err := cli.Write(4, modbusBulkWrite(pointsList.DeviceFour))
		if err != nil {
			log.Error(err)
		}
	}

}

func (inst *Server) modbusInputsRunner(cli *modbuscli.Modbus, pointsList []*points.Point) {
	var err error
	var tempList [8]float64
	var voltList [8]float64
	var completedTemp bool
	var completedVolt bool
	var returnedValue float64
	for _, point := range pointsList { // do modbus read
		if !point.IsWriteable {
			addr, _ := points.ModbusBuildInput(point.IoType, point.ObjectID)
			slaveId := addr.DeviceAddr
			if slaveId <= 0 {
				log.Errorf("modbus slave addrress cant not be less to 1")
				continue
			}
			if !completedTemp && (point.IoType == points.IoTypeTemp || point.IoType == points.IoTypeDigital) {
				tempList, err = cli.ReadTemps(slaveId) // DO MODBUS READ FOR TEMPS
				if err != nil {
					log.Errorf("modbus read temp %s", err.Error())
				}
			}
			if !completedVolt && point.IoType == points.IoTypeVolts {
				voltList, err = cli.ReadVolts(slaveId) // DO MODBUS READ FOR VOLTS
				if err != nil {
					log.Errorf("modbus read voltages %s", err.Error())
				}
			}
			// update the store
			io16Pin := addr.IoPin - 1
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
			inst.writePV(point.ObjectType, point.ObjectID, returnedValue)
		}
	}

}
