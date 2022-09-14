package node

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

type NodesList struct {
	Nodes interface{} `json:"nodes"`
}

func TestSpec_NodeConnection(t *testing.T) {

	var list []*Schema
	var (
		value = map[string]map[string][]Links{"flow": map[string][]Links{"links": []Links{
			Links{
				NodeId: "2",
				Socket: "flow",
			},
		}}}
	)

	s1 := &Schema{
		Id:   "1",
		Type: "time/delay",
		Metadata: &Metadata{
			PositionX: "271.5",
			PositionY: "-69",
		},
		Inputs: value,
	}

	s2 := &Schema{
		Id:   "2",
		Type: "time/delay",
		Metadata: &Metadata{
			PositionX: "271.5",
			PositionY: "-69",
		},
	}

	list = append(list, s1)
	list = append(list, s2)
	a := NodesList{
		Nodes: list,
	}

	pprint.PrintJOSN(a)

}

func TestSpec_NodeNonConnection(t *testing.T) {
	var list []*Schema
	var value = map[string]map[string]string{"duration": map[string]string{"value": "22"}}
	s1 := &Schema{
		Id:   "2",
		Type: "time/delay",
		Metadata: &Metadata{
			PositionX: "271.5",
			PositionY: "-69",
		},
		Inputs: value,
	}

	list = append(list, s1)
	a := NodesList{
		Nodes: list,
	}

	pprint.PrintJOSN(a)

}
