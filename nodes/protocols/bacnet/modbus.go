package bacnetio

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/modbuscli"
	"github.com/NubeIO/nubeio-rubix-lib-modbus-go/modbus"
	"github.com/enescakir/emoji"
	log "github.com/sirupsen/logrus"
	"time"
)

const errNoDevices = "no IO16 where selected so dont inti modbus"

func (inst *Server) initModbus(settings map[string]interface{}) (*modbuscli.Modbus, error) {
	schema, err := GetBacnetSchema(settings)
	if err != nil {
		log.Error(err)
		inst.SetStatusError("failed to set node settings")
		inst.SetErrorIcon(string(emoji.RedCircle))
		return nil, err
	}
	if schema.DeviceCount == noDevices {
		return nil, errors.New(errNoDevices)
	}
	inst.deviceCount = schema.DeviceCount
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
		return nil, err
	}
	return init, nil
}

// modbus will come from polling
// this is to only work for the IO-16
func (inst *Server) modbusRunner(settings map[string]interface{}) {
	log.Info("start modbus-runner")

	mb, err := inst.initModbus(settings)
	var dontPollModbus bool
	if err != nil {
		if err.Error() == errNoDevices {
			dontPollModbus = true
		} else {
			return
		}
	}
	var count int64
	for {
		log.Infof("modbus polling loop count: %d application-type: %s", count, inst.application)
		inst.pollingCount = count
		pointsListRead, _ := inst.getPointsReadOnly()
		if !dontPollModbus {
			inst.modbusInputsRunner(mb, pointsListRead) // process the inputs
			time.Sleep(modbusDelay * time.Millisecond)
			if count > 2 { // make sure all the inputs have been read before doing the writes, this is to try and stop toggling a relay output
				inst.modbusOutputsDispatch(mb)
			}
		}
		go inst.writeRunner() // publish all mqtt values
		time.Sleep(500 * time.Millisecond)
		inst.SetNotifyMessage(pollStats(inst.pollingCount))
		inst.SetNotifyIcon(string(emoji.GreenCircle)) // process the outs
		count++
	}
}

func modbusScaleOutput(value, offset float64) float64 {
	value = value + offset
	if value >= 10 {
		value = 10
	}
	if value <= 0 {
		value = 0
	}
	return value
}

func modbusBulkWrite(pointsList []*points.Point) [8]float64 {
	var out [8]float64
	for _, point := range pointsList {
		ioNumber, _ := points.ModbusBuildOutput(points.IoTypeVolts, point.ObjectID)
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
				} else {
					value = modbusScaleOutput(value, point.Offset) // point offset
				}
			}
			out[ioNumber.IoPin-1] = value
		}
	}
	return out
}

func (inst *Server) modbusOutputsDispatch(cli *modbuscli.Modbus) {
	pointsList := inst.GetModbusWriteablePoints()
	if pointsList == nil {
		log.Errorf("modbus modbusOutputsDispatch() points is empty")
		return

	}
	if len(pointsList.DeviceOne) > 0 {
		err := cli.Write(1, modbusBulkWrite(pointsList.DeviceOne))
		setReadError(1, err)
		if err != nil {
			log.Errorf("modbus write %s slave: %d", err.Error(), 1)
		}
		time.Sleep(modbusDelay * time.Millisecond)
	}
	if len(pointsList.DeviceTwo) > 0 {
		err := cli.Write(2, modbusBulkWrite(pointsList.DeviceTwo))
		setReadError(2, err)
		if err != nil {
			log.Errorf("modbus write %s slave: %d", err.Error(), 2)
		}
		time.Sleep(modbusDelay * time.Millisecond)
	}
	if len(pointsList.DeviceThree) > 0 {
		err := cli.Write(3, modbusBulkWrite(pointsList.DeviceThree))
		setReadError(3, err)
		if err != nil {
			log.Errorf("modbus write %s slave: %d", err.Error(), 3)
		}
		time.Sleep(modbusDelay * time.Millisecond)
	}
	if len(pointsList.DeviceFour) > 0 {
		err := cli.Write(4, modbusBulkWrite(pointsList.DeviceFour))
		setReadError(4, err)
		if err != nil {
			log.Errorf("modbus write %s slave: %d", err.Error(), 4)
		}
		time.Sleep(modbusDelay * time.Millisecond)
	}
}

var firstLoop = true

var readError1 bool
var readError2 bool
var readError3 bool
var readError4 bool

func setReadError(slaveID int, err error) {
	if slaveID == 1 {
		if err != nil {
			readError1 = true
		} else {
			readError1 = false
		}
	}
	if slaveID == 2 {
		if err != nil {
			readError2 = true
		} else {
			readError2 = false
		}
	}
	if slaveID == 3 {
		if err != nil {
			readError3 = true
		} else {
			readError3 = false
		}
	}
	if slaveID == 4 {
		if err != nil {
			readError4 = true
		} else {
			readError4 = false
		}
	}
}

func (inst *Server) modbusInputsRunner(cli *modbuscli.Modbus, pointsList []*points.Point) {
	var err error
	var tempList [8]float64
	var voltList [8]float64
	var currentList [8]float64
	var diList [8]float64
	var completedTemp bool
	var completedVolt bool
	var completedCurrent bool
	var completedDis bool
	var returnedValue float64
	var slaveId int
	for _, point := range pointsList { // do modbus read
		if !point.IsWriteable {
			addr, _ := points.ModbusBuildInput(point.IoType, point.ObjectID)
			slaveId = addr.DeviceAddr
			io16Pin := addr.IoPin - 1
			if slaveId <= 0 {
				log.Errorf("modbus slave addrress cant not be less to 1")
				continue
			}

			if firstLoop { // setup all the pulse inputs
				inst.modbusPointSetup(cli, point, slaveId, addr.IoPin)
			}

			if !completedTemp && (point.IoType == points.IoTypeTemp) {
				tempList, err = cli.ReadTemps(slaveId) // DO MODBUS READ FOR TEMPS OR DIs
				if err != nil {
					log.Errorf("modbus read temp %s slave: %d", err.Error(), slaveId)
				} else {
					returnedValue = tempList[io16Pin]
					err := inst.writePV(point.ObjectType, point.ObjectID, returnedValue)
					if err != nil {
						log.Errorf("modbus modbusInputsRunner() writePv %s slave: %d", err.Error(), slaveId)
					}
				}
				setReadError(slaveId, err)
				time.Sleep(modbusDelay * time.Millisecond)
			}
			if !completedVolt && point.IoType == points.IoTypeVolts {
				voltList, err = cli.ReadVolts(slaveId) // DO MODBUS READ FOR VOLTS
				if err != nil {
					log.Errorf("modbus read voltages %s slave: %d", err.Error(), slaveId)
				} else {
					returnedValue = voltList[io16Pin]
					err := inst.writePV(point.ObjectType, point.ObjectID, returnedValue)
					if err != nil {
						log.Errorf("modbus modbusInputsRunner() writePv %s slave: %d", err.Error(), slaveId)
					}
				}
				setReadError(slaveId, err)
				time.Sleep(modbusDelay * time.Millisecond)
			}
			if !completedCurrent && point.IoType == points.IoTypeCurrent {
				currentList, err = cli.ReadCurrent(slaveId) // DO MODBUS READ FOR VOLTS
				if err != nil {
					log.Errorf("modbus read current %s slave: %d", err.Error(), slaveId)
				} else {
					returnedValue = currentList[io16Pin]
					err := inst.writePV(point.ObjectType, point.ObjectID, returnedValue)
					if err != nil {
						log.Errorf("modbus modbusInputsRunner() writePv %s slave: %d", err.Error(), slaveId)
					}
				}
				setReadError(slaveId, err)
				time.Sleep(modbusDelay * time.Millisecond)
			}
			// update the store
			if !completedDis && point.IoType == points.IoTypeDigital { // update anypoint that is type temp
				diList, err = cli.ReadDIs(slaveId) // DO MODBUS READ FOR TEMPS OR DIs
				if err != nil {
					log.Errorf("modbus read DIs %s slave: %d", err.Error(), slaveId)
				} else {
					returnedValue = diList[io16Pin]
					err := inst.writePV(point.ObjectType, point.ObjectID, returnedValue)
					if err != nil {
						log.Errorf("modbus modbusInputsRunner() writePv %s slave: %d", err.Error(), slaveId)
					}
				}
				setReadError(slaveId, err)
				time.Sleep(modbusDelay * time.Millisecond)
			}

			if !completedDis && (point.IoType == points.IoTypePulseOnRise || point.IoType == points.IoTypePulseOnFall) { // update anypoint that is type temp
				diList, err = cli.ReadPulse(slaveId) // DO MODBUS READ FOR PULSE
				if err != nil {
					log.Errorf("modbus read pulse %s slave: %d", err.Error(), slaveId)
				} else {
					returnedValue = diList[io16Pin]
					err := inst.writePV(point.ObjectType, point.ObjectID, returnedValue)
					if err != nil {
						log.Errorf("modbus modbusInputsRunner() writePv %s slave: %d", err.Error(), slaveId)
					}
				}
				setReadError(slaveId, err)
				time.Sleep(modbusDelay * time.Millisecond)
			}
		}
	}
	if readError1 {
		inst.setDevStats1("offline")
	} else {
		inst.setDevStats1("ok")
	}
	if readError2 {
		inst.setDevStats2("offline")
	} else {
		inst.setDevStats2("ok")
	}
	if readError3 {
		inst.setDevStats3("offline")
	} else {
		inst.setDevStats3("ok")
	}
	if readError4 {
		inst.setDevStats4("offline")
	} else {
		inst.setDevStats4("ok")
	}
	firstLoop = false
}

func (inst *Server) modbusPointSetup(cli *modbuscli.Modbus, point *points.Point, slaveId, ioPin int) {
	if ioPin == 1 { // UI1
		var value uint16
		if point.IoType == points.IoTypePulseOnRise { // 8: Pulse RISE
			value = 8
			err := cli.WriteRegisterInt16(slaveId, 5200, value)
			if err != nil {
				log.Errorf("modbus WriteRegister for pulse err: %s slave: %d io16Pin: %d value: %d", err.Error(), slaveId, ioPin, value)
			}
		} else if point.IoType == points.IoTypePulseOnFall { // 9: Pulse FALL
			value = 9
			err := cli.WriteRegisterInt16(slaveId, 5200, value)
			if err != nil {
				log.Errorf("modbus WriteRegister for pulse err: %s slave: %d io16Pin: %d value: %d", err.Error(), slaveId, ioPin, value)
			}
		} else {
			err := cli.WriteRegisterInt16(slaveId, 5200, 0)
			if err != nil {
				log.Errorf("modbus WriteRegister for pulse err: %s slave: %d io16Pin: %d value: %d", err.Error(), slaveId, ioPin, value)
			}
		}
	}
	if ioPin == 2 { // UI2
		var value uint16
		if point.IoType == points.IoTypePulseOnRise { // 8: Pulse RISE
			value = 8
			err := cli.WriteRegisterInt16(slaveId, 5201, value)
			if err != nil {
				log.Errorf("modbus WriteRegister for pulse err: %s slave: %d io16Pin: %d value: %d", err.Error(), slaveId, ioPin, value)
			}
		} else if point.IoType == points.IoTypePulseOnFall { // 9: Pulse FALL
			value = 9
			err := cli.WriteRegisterInt16(slaveId, 5201, value)
			if err != nil {
				log.Errorf("modbus WriteRegister for pulse err: %s slave: %d io16Pin: %d value: %d", err.Error(), slaveId, ioPin, value)
			}
		} else {
			err := cli.WriteRegisterInt16(slaveId, 5201, 0)
			if err != nil {
				log.Errorf("modbus WriteRegister for pulse err: %s slave: %d io16Pin: %d value: %d", err.Error(), slaveId, ioPin, value)
			}
		}
	}
	if ioPin == 3 { // UI2
		var value uint16
		if point.IoType == points.IoTypePulseOnRise { // 8: Pulse RISE
			value = 8
			err := cli.WriteRegisterInt16(slaveId, 5202, value)
			if err != nil {
				log.Errorf("modbus WriteRegister for pulse err: %s slave: %d io16Pin: %d value: %d", err.Error(), slaveId, ioPin, value)
			}
		} else if point.IoType == points.IoTypePulseOnFall { // 9: Pulse FALL
			value = 9
			err := cli.WriteRegisterInt16(slaveId, 5202, value)
			if err != nil {
				log.Errorf("modbus WriteRegister for pulse err: %s slave: %d io16Pin: %d value: %d", err.Error(), slaveId, ioPin, value)
			}
		} else {
			err := cli.WriteRegisterInt16(slaveId, 5202, 0)
			if err != nil {
				log.Errorf("modbus WriteRegister for pulse err: %s slave: %d io16Pin: %d value: %d", err.Error(), slaveId, ioPin, value)
			}
		}
	}
}
