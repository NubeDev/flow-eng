package modbuscli

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-modbus-go/modbus"
	"testing"
)

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
