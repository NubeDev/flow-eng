package points

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"testing"
)

func Test_calcModbusRubix(t *testing.T) {
	uiStart, uoStart := calcModbusRubix(2, true, true, false)
	fmt.Println(uiStart, uoStart)
}

func TestNewStore(t *testing.T) {
	bs := New(applications.Edge, nil, 1, 1, 1)

	var st []*Point
	pprint.PrintJOSN(bs)

	ai1 := &Point{
		Application: applications.Edge,
		ObjectType:  AnalogInput,
		ObjectID:    1,
	}
	ai2 := &Point{
		Application: applications.Edge,
		ObjectType:  AnalogInput,
		ObjectID:    2,
	}
	err, _ := bs.AddPoint(ai1)
	if err != nil {
		fmt.Println(err)
		return
	}
	err, _ = bs.AddPoint(ai2)
	if err != nil {
		fmt.Println(err)
		return
	}

	pprint.PrintJOSN(st)
}
