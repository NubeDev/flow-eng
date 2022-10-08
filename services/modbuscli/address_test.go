package modbuscli

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
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
	s := points.New(names.Modbus, nil, 2, 200, 200)
	fmt.Println(init)

	s.AddPoint(&points.Point{
		ObjectType:  points.AnalogOutput,
		ObjectID:    1,
		IoType:      points.IoTypeVolts,
		IsIO:        true,
		IsWriteable: true,
	}, true)
	if err != nil {
		return
	}

	pointsList := s.GetPointsByApplication(names.Modbus)
	//var pointsToWrite []*points.Point
	for _, point := range pointsList {
		fmt.Println(point)

	}

	var a [8]float64
	a[0] = 11.1
	fmt.Println(a)

}
