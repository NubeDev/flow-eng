package rubixIO

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/services/clients/rubixcli"
	"testing"
)

func TestNew(t *testing.T) {

	c := New()
	ping, err := c.rest.Ping()
	fmt.Println(err)
	if err != nil {
		return
	}
	pprint.Print(ping)

}

func TestWrite(t *testing.T) {

	c := New()
	out1 := &rubixcli.Output{
		IoNumber: "UO1",
		Value:    0,
	}
	out2 := &rubixcli.Output{
		IoNumber: "DO2",
		Value:    1,
	}
	write, err := c.rest.BulkWrite(c.rest.BulkWriteBuilder(out1, out2))
	fmt.Println(err)
	if err != nil {
		return
	}
	pprint.Print(write)

}
