package main

import (
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/_example/nodes"
	"log"
	"time"

	"github.com/NubeDev/flow-eng/node"
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
		NodeID:   "a123",
		Name:     "a123",
		NodeType: "PASS",
		InputList: []*node.TypeInput{&node.TypeInput{
			PortCommon: &node.PortCommon{
				Name: "in1",
				Type: "int8",
				Connection: &node.Connection{
					NodeID: "",
					Port:   "",
				},
			},
		}},
		OutputList: []*node.TypeOutput{&node.TypeOutput{
			PortCommonOut: &node.PortCommonOut{
				Name: "out1",
				Type: "int8",
				Connections: []*node.Connection{&node.Connection{
					NodeID: "b123",
					Port:   "in1",
				}},
			},
		}},
	})

	nodeB := nodes.New(&nodes.Node{
		NodeID:   "b123",
		Name:     "b123",
		NodeType: "PASS",
		InputList: []*node.TypeInput{&node.TypeInput{
			PortCommon: &node.PortCommon{
				Name: "in1",
				Type: "int8",
				Connection: &node.Connection{
					NodeID: "b123",
					Port:   "out1",
				},
			},
		}},
		OutputList: []*node.TypeOutput{&node.TypeOutput{
			PortCommonOut: &node.PortCommonOut{
				Name: "out1",
				Type: "int8",
				Connections: []*node.Connection{&node.Connection{
					NodeID: "",
					Port:   "",
				}},
			},
		}},
	})

	var nodesList []*nodes.Node
	nodesList = append(nodesList, nodeA)
	nodesList = append(nodesList, nodeB)

	graph := flowctrl.New()
	graph.AddNode(nodeA)
	graph.AddNode(nodeB)

	getA := graph.GetNode(nodeA.Name).(*nodes.Node)
	getB := graph.GetNode(nodeB.Name).(*nodes.Node)

	for _, output := range getA.OutputList {
		for _, input := range getB.InputList {
			fmt.Println(11111, "SET")
			output.OutputPort.Connect(input.InputPort)
		}

	}

	graph.ReplaceNode(nodeA.Name, getA)
	graph.ReplaceNode(nodeB.Name, getB)

	for _, ordered := range graph.Get().Graphs {
		for _, runner := range ordered.Runners {
			//fmt.Println("RUNNER-1", runner.Name(), len(runner.Outputs()), "LEN", "UUID", runner.UUID())

			for _, port := range runner.Outputs() {
				for _, connector := range port.Connectors() {

					fmt.Println("RUNNER-TO-------------", connector.FromUUID(), connector.ToUUID())
				}
			}
		}
	}

	runner := flowctrl.NewSerialRunner(graph)

	log.Println("Flow started")
	for {

		err := runner.Process()
		if err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Second)
		//// wait for delayed node to propagate data
		//if endReader.Get() != 0 {
		//	break
		//}
	}

}
