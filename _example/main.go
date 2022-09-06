package main

import (
	"github.com/NubeDev/flow-eng/_example/nodes"
	pprint "github.com/NubeDev/flow-eng/print"
)

type typeInput struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Connection *connection `json:"connection"`
}

type connection struct {
	NodeID string `json:"nodeID"`
	Port   string `json:"port"`
}

type typeOutput struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Connections []*connection `json:"connections"`
}

type commonNode struct {
	NodeID  string        `json:"nodeID"` // abc
	Name    string        `json:"name"`   // my node
	Node    string        `json:"node"`   // PASS
	Inputs  []*typeInput  `json:"inputs"`
	Outputs []*typeOutput `json:"outputs"`
}

func main() {

	nodeA := nodes.New(&nodes.Node{
		NodeID: "a123",
		Name:   "a123",
		Node:   "PASS",
		Inputs: []*nodes.TypeInput{&nodes.TypeInput{
			PortCommon: &nodes.PortCommon{
				Name: "in1",
				Type: "int8",
				Connection: &nodes.Connection{
					NodeID: "",
					Port:   "",
				},
			},
		}},
		Outputs: []*nodes.TypeOutput{&nodes.TypeOutput{
			PortCommonOut: &nodes.PortCommonOut{
				Name: "out1",
				Type: "int8",
				Connections: []*nodes.Connection{&nodes.Connection{
					NodeID: "b123",
					Port:   "in1",
				}},
			},
		}},
	})

	nodeB := nodes.New(&nodes.Node{
		NodeID: "b123",
		Name:   "b123",
		Node:   "PASS",
		Inputs: []*nodes.TypeInput{&nodes.TypeInput{
			PortCommon: &nodes.PortCommon{
				Name: "in1",
				Type: "int8",
				Connection: &nodes.Connection{
					NodeID: "b123",
					Port:   "out1",
				},
			},
		}},
		Outputs: []*nodes.TypeOutput{&nodes.TypeOutput{
			PortCommonOut: &nodes.PortCommonOut{
				Name: "out1",
				Type: "int8",
				Connections: []*nodes.Connection{&nodes.Connection{
					NodeID: "",
					Port:   "",
				}},
			},
		}},
	})

	var nodesList []*nodes.Node
	nodesList = append(nodesList, nodeA)
	nodesList = append(nodesList, nodeB)
	pprint.PrintJOSN(nodesList)

	//nodeB := nodes.NewPass(nil)

	//nameA_ := &genericNode{
	//	NodeName: "nodeA",
	//	NodeType: "PASS",
	//	Node:     nodeA,
	//}
	//
	//nameB_ := &genericNode{
	//	NodeName: "nodeB",
	//	NodeType: "PASS",
	//	Node:     nodeB,
	//}
	//
	//var nodesParsed []*genericNode
	//
	//nodesParsed = append(nodesParsed, nameA_)
	//nodesParsed = append(nodesParsed, nameB_)

}
