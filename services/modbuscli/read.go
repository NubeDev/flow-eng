package modbuscli

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-modbus-go/modbus"
)

func tempRegs() (start int, count int) {
	start = 0
	return start, start + 8
}

func voltRegs() (start int, count int) {
	start = 200
	return start, 8
}

func currentRegs() (start int, count int) {
	start = 300
	return start, 8
}

// DecodeUI use this for when reading as a temp, voltage or amps
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
	start, count := tempRegs()
	registers, err := inst.readRegisters(slave, start, count, false)
	if err != nil {
		return raw, err
	}
	if len(registers) < 8 {
		return raw, errors.New("read length must be 8")
	}
	return convert(registers), err
}

func (inst *Modbus) ReadVolts(slave int) (raw [8]float64, err error) {
	start, count := voltRegs()
	registers, err := inst.readRegisters(slave, start, count, false)
	if err != nil {
		return raw, err
	}
	if len(registers) < 8 {
		return raw, errors.New("read length must be 8")
	}
	return convert(registers), err
}

func (inst *Modbus) ReadCurrent(slave int) (raw [8]float64, err error) {
	start, count := currentRegs()
	registers, err := inst.readRegisters(slave, start, count, false)
	if err != nil {
		return raw, err
	}
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
	// log.Print(out)
	return out
}

func (inst *Modbus) readRegisters(slave, start, count int, holding bool) (raw []byte, err error) {
	err = inst.SetSlave(slave)
	if err != nil {
		return nil, err
	}
	if holding {
		registers, _, err := inst.client.ReadHoldingRegisters(uint16(start), uint16(count))
		// log.Print(registers)
		return registers, err
	} else {
		// fmt.Println("READ INPUTS", slave, uint16(start), uint16(count))
		registers, _, err := inst.client.ReadInputRegisters(uint16(start), uint16(count))
		// log.Print(registers)
		return registers, err
	}
}
