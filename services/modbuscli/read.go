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

func pulseRegs() (start int, count int) {
	start = 400
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

func decodeDI(in bool) float64 {
	if in {
		return 1
	}
	return 0
}

func decodeDIs(in []bool) [8]float64 {
	var out [8]float64
	for i, b := range in {
		out[i] = decodeDI(b)
	}
	return out
}

func (inst *Modbus) ReadTemps(slave int) (raw [8]float64, err error) {
	start, count := tempRegs()
	registers, _, err := inst.readRegisters(slave, start, count, false, false)
	if err != nil {
		return raw, err
	}
	if len(registers) < 8 {
		return raw, errors.New("read length must be 8")
	}
	return convert(registers), err
}

func (inst *Modbus) ReadDIs(slave int) (raw [8]float64, err error) {
	start, count := tempRegs()
	registers, err := inst.readDiscreteInputs(slave, start, count)
	if err != nil {
		return raw, err
	}
	if len(registers) < 8 {
		return raw, errors.New("read length must be 8")
	}
	return decodeDIs(registers), err
}

func (inst *Modbus) ReadPulse(slave int) (raw [8]float64, err error) {
	start, count := pulseRegs()
	_, registers, err := inst.readRegisters(slave, start, count, false, true)
	if err != nil {
		return raw, err
	}
	if len(registers) < 8 {
		return raw, errors.New("read length must be 8")
	}
	return convertUint32(registers), err
}

func (inst *Modbus) ReadVolts(slave int) (raw [8]float64, err error) {
	start, count := voltRegs()
	registers, _, err := inst.readRegisters(slave, start, count, false, false)
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
	registers, _, err := inst.readRegisters(slave, start, count, false, false)
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

func convertUint32(raw []uint32) [8]float64 {
	out := [8]float64{}
	for i, u := range raw {
		out[i] = float64(u)
	}
	// log.Print(out)
	return out
}

func (inst *Modbus) readDiscreteInputs(slave, start, count int) (raw []bool, err error) {
	err = inst.SetSlave(slave)
	if err != nil {
		return nil, err
	}
	registers, _, err := inst.client.ReadDiscreteInputs(uint16(start), uint16(count))
	// log.Print(registers)
	return registers, err
}

func (inst *Modbus) readRegisters(slave, start, count int, holding, uint32 bool) (raw []byte, uint32Res []uint32, err error) {
	err = inst.SetSlave(slave)
	if err != nil {
		return nil, nil, err
	}
	if uint32 {
		inst.client.SetEncoding(modbus.BigEndian, modbus.LowWordFirst)
		registers, err := inst.client.ReadUint32s(uint16(start), uint16(count), modbus.InputRegister)
		// log.Print(registers)
		return nil, registers, err
	}

	if holding {
		registers, _, err := inst.client.ReadHoldingRegisters(uint16(start), uint16(count))
		// log.Print(registers)
		return registers, nil, err
	} else {
		// fmt.Println("READ INPUTS", slave, uint16(start), uint16(count))
		registers, _, err := inst.client.ReadInputRegisters(uint16(start), uint16(count))
		// log.Print(registers)
		return registers, nil, err
	}
}
