package bacnet

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"testing"
)

func TestNewStore(t *testing.T) {
	NewStore(applications.Edge, nil)

	var st []*Point

	ai1 := &Point{
		Application: applications.Edge,
		ObjectType:  AnalogInput,
		ObjectID:    1,
	}
	ai2 := &Point{
		Application: applications.Edge,
		ObjectType:  AnalogOutput,
		ObjectID:    1,
	}
	st = append(st, ai1)
	fmt.Println(CheckExistingPoint(st, ai2))

	pprint.PrintJOSN(st)

}
