package edge28lib

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"testing"
)

var c = New("192.168.15.141")

func TestNew(t *testing.T) {
	server, err := c.client.PingServer()
	pprint.Print(server)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestGetDIs(t *testing.T) {
	p1 := &points.Point{
		ObjectType: points.BinaryInput,
		ObjectID:   1,
		IoType:     points.IoTypeDigital,
	}
	p2 := &points.Point{
		ObjectType: points.BinaryInput,
		ObjectID:   2,
		IoType:     points.IoTypeDigital,
	}
	server, err := c.GetDIs(p1, p2)
	fmt.Println(err)
	pprint.PrintJOSN(server)
}

func TestGetUIs(t *testing.T) {
	p1 := &points.Point{
		ObjectType: points.AnalogInput,
		ObjectID:   1,
		IoType:     points.IoTypeTemp,
	}
	p2 := &points.Point{
		ObjectType: points.AnalogInput,
		ObjectID:   2,
		IoType:     points.IoTypeDigital,
	}
	server, err := c.GetUIs(p1, p2)
	fmt.Println(err)
	pprint.PrintJOSN(server)
}
