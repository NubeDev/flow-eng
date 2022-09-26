package ffclient

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	"testing"
)

func TestFlowClient_PointWriteByName(t *testing.T) {
	c := New(&Connection{})

	name, err := c.PointWriteByName("net", "dev", "pnt", &Priority{P1: float.New(33)})
	fmt.Println(name, err)
	name, err = c.PointWrite("pnt_c82138298ea14f40", &Priority{P1: float.New(44)})
	fmt.Println(name, err)

}
