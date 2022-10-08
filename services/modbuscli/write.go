package modbuscli

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-modbus-go/modbus"
)

/*
Send all the pending writes
Sort by the bacnet address as in AO1,2,3 up to 8 (this would be modbus device address 1) AO9,10, 11 up to 16 (device address 2)
*/

func (inst *Modbus) Write(slave int, values [8]float64) (err error) {
	bulk, err := buildWrite(values[0], values[1], values[2], values[3], values[4], values[5], values[6], values[7])
	if err != nil {
		return err
	}
	_, err = inst.writeRegisters(slave, bulk)
	if err != nil {
		return err
	}
	return nil
}

func (inst *Modbus) writeRegisters(slave int, values []byte) (raw []byte, err error) {
	err = inst.SetSlave(slave)
	if err != nil {
		return nil, err
	}
	return inst.client.WriteMultipleRegisters(0, 8, values)
}

func encodeUO(in float64) uint16 {
	in = in * 100
	return uint16(in)
}

func buildWrite(uo1, uo2, uo3, uo4, uo5, uo6, uo7, uo8 float64) ([]byte, error) {
	if uo1 > 10 || uo2 > 10 || uo3 > 10 || uo4 > 10 || uo5 > 10 || uo6 > 10 || uo7 > 10 || uo8 > 10 {
		return nil, errors.New("modbus write values must be less then 10")
	}
	in := []uint16{
		encodeUO(uo1),
		encodeUO(uo2),
		encodeUO(uo3),
		encodeUO(uo4),
		encodeUO(uo5),
		encodeUO(uo6),
		encodeUO(uo7),
		encodeUO(uo8),
	}
	v := modbus.Uint16sToBytes(modbus.BigEndian, in)
	return v, nil
}
