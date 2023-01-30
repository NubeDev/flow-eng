package modbuscli

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-modbus-go/modbus"
	"testing"
)

func Test_Write(t *testing.T) {

	cli := &Modbus{
		IsSerial: true,
		Serial: &modbus.Serial{
			SerialPort: "/dev/ttyUSB0",
		},
	}
	init, err := cli.Init(cli)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = init.WriteRegister(2, 0, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

}

func Test_ReadTemp(t *testing.T) {

	cli := &Modbus{
		IsSerial: true,
		Serial: &modbus.Serial{
			SerialPort: "/dev/ttyUSB0",
		},
	}
	init, err := cli.Init(cli)
	if err != nil {
		fmt.Println(err)
		return
	}
	registers, err := init.ReadTemps(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(registers)

}

func Test_ReadVoltage(t *testing.T) {

	cli := &Modbus{
		IsSerial: true,
		Serial: &modbus.Serial{
			SerialPort: "/dev/ttyUSB0",
		},
	}
	init, err := cli.Init(cli)
	if err != nil {
		fmt.Println(err)
		return
	}
	registers, err := init.ReadVolts(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(registers)

}

func Test_Read(t *testing.T) {

	cli := &Modbus{
		IsSerial: true,
		Serial: &modbus.Serial{
			SerialPort: "/dev/ttyUSB0",
		},
	}
	init, err := cli.Init(cli)
	if err != nil {
		fmt.Println(err)
		return
	}
	registers, err := init.readRegisters(1, 200, 2, false)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(registers)

}
