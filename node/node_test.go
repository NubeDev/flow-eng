package node

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

type NodesList struct {
	Nodes interface{} `json:"nodes"`
}

func TestBaseNode_GetInfo(t *testing.T) {

	var list []*Schema

	var links = map[string]Links{"a": Links{
		Value: 222,
	}}
	s1 := &Schema{
		Id:   "111",
		Type: "logic/numberConstant",
		Metadata: &Metadata{
			PositionX: "123",
			PositionY: "1235",
		},
		Inputs: &Inputs{
			Links: links,
		},
	}

	list = append(list, s1)
	a := NodesList{
		Nodes: list,
	}

	pprint.PrintJOSN(a)

}
