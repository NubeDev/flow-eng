package edge28lib

import (
	"fmt"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"testing"
)

func TestEdge28_WriteDO(t *testing.T) {
	arr := points.NewPriArrayAt15(0)
	p1 := &points.Point{
		ObjectType: points.BinaryOutput,
		ObjectID:   1,
		IoType:     points.IoTypeDigital,
		WriteValue: arr,
	}
	do, err := c.WriteDO(p1)
	fmt.Println(do, err)
	if err != nil {
		return
	}
}

func TestEdge28_WriteUOasDO(t *testing.T) {
	arr := points.NewPriArrayAt15(9.9)
	p1 := &points.Point{
		ObjectType: points.AnalogOutput,
		ObjectID:   2,
		IoType:     points.IoTypeVolts,
		WriteValue: arr,
	}
	do, err := c.WriteUO(p1)
	fmt.Println(do, err)
	if err != nil {
		return
	}
}
