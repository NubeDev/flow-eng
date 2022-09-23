package modbuscli

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-modbus-go/modbus"
)

func encodeUO(in float64) uint16 {
	in = in * 100
	return uint16(in)
}

func BulkWrite(uo1, uo2, uo3, uo4, uo5, uo6, uo7, uo8 float64) ([]byte, error) {
	if uo1 > 10 || uo2 > 10 || uo3 > 10 || uo4 > 10 || uo5 > 10 || uo6 > 10 || uo7 > 10 || uo8 > 10 {
		return nil, errors.New("values must be less then 10")
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
