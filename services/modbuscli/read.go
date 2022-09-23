package modbuscli

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-modbus-go/modbus"
	log "github.com/sirupsen/logrus"
)

func tempRegs() (start int, finish int) {
	start = 0
	return start, start + 8
}

func voltRegs() (start int, finish int) {
	start = 250
	return start, start + 8
}

//DecodeUI use this for when reading as a temp, voltage or amps
func DecodeUI(in uint16) float64 {
	v := float64(in) / 100
	return v
}

// TempToDI read the io-16 as temp and convert to a DI
func TempToDI(v float64) float64 {
	if v > 1 {
		return 1
	}
	return 0
}

func (inst *Modbus) ReadTemps(slave int) (raw [8]float64, err error) {
	start, finish := tempRegs()
	registers, err := inst.readRegisters(slave, start, finish, false)
	if len(registers) < 8 {
		return raw, errors.New("read length must be 8")
	}
	return convert(registers), err
}

func (inst *Modbus) ReadVolts(slave int) (raw [8]float64, err error) {
	start, finish := voltRegs()
	registers, err := inst.readRegisters(slave, start, finish, false)
	if len(registers) < 8 {
		return raw, errors.New("read length must be 8")
	}
	return convert(registers), err
}

func convert(raw []byte) [8]float64 {
	decode := modbus.BytesToUint16s(modbus.BigEndian, raw)
	out := [8]float64{}
	for i, u := range decode {
		out[i] = DecodeUI(u)
	}
	return out
}

func (inst *Modbus) readRegisters(slave, start, finish int, holding bool) (raw []byte, err error) {
	err = inst.SetSlave(slave)
	if err != nil {
		return nil, err
	}
	if holding {
		registers, _, err := inst.client.ReadHoldingRegisters(uint16(start), uint16(finish))
		if err != nil {
			log.Error(err)
			return registers, err
		}
	} else {
		registers, _, err := inst.client.ReadInputRegisters(uint16(start), uint16(finish))
		if err != nil {
			log.Error(err)
			return registers, err
		}
	}
	return nil, err
}
