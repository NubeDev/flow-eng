package modbuscli

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	log "github.com/sirupsen/logrus"
)

func tempRegs() (start int, finish int) {
	start = 1
	return start, start + 8
}

func voltRegs() (start int, finish int) {
	start = 250
	return start, start + 8
}

// TempToDI read the io-16 as temp and convert to a DI
func TempToDI(v float64) float64 {
	if v > 100 { // TODO check calc
		return 1
	}
	return 0
}

func (inst *Modbus) ReadTemps(slave int) (raw []float64, err error) {
	start, finish := tempRegs()
	registers, err := inst.readRegisters(slave, start, finish, false)
	return convert(registers), err
}

func (inst *Modbus) ReadVolts(slave int) (raw []float64, err error) {
	start, finish := voltRegs()
	registers, err := inst.readRegisters(slave, start, finish, false)
	return convert(registers), err
}

func convert(raw []byte) []float64 {
	a := float.RandFloat(1, 2)
	return []float64{122.2, a, 3, 4, 5, 6, 7, 8}
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
