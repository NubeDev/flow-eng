package main

import (
	"encoding/json"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/_example/nodes"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {

	var nodesParsed []*nodes.Node
	jsonFile, err := os.Open("../flow-eng/_example/test.json")
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &nodesParsed)
	//pprint.PrintJOSN(nodesParsed)

	var nodeA *nodes.Node
	var nodeB *nodes.Node

	for _, n := range nodesParsed {
		if n.GetName() == "nodeA" {
			nodeA = n
		}
		if n.GetName() == "nodeB" {
			nodeB = n
		}
	}

	nodeA, _ = nodes.New(nodeA)
	nodeB, _ = nodes.New(nodeB)

	graph := flowctrl.New()
	graph.AddNode(nodeA)
	graph.AddNode(nodeB)

	getA := graph.GetNode(nodeA.GetID())
	getB := graph.GetNode(nodeB.GetID())

	for _, output := range getA.GetOutputs() {
		for _, input := range getB.GetInputs() {
			output.OutputPort.Connect(input.InputPort)
		}
	}

	graph.ReplaceNode(nodeA.GetID(), getA)
	graph.ReplaceNode(nodeB.GetID(), getB)

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

func buildJson() {

	nodeA, _ := nodes.New(&nodes.Node{
		InputList: []*node.TypeInput{&node.TypeInput{
			PortCommon: &node.PortCommon{
				Name: "in1",
				Type: "int8",
				Connection: &node.Connection{
					NodeID:   "",
					NodePort: "",
				},
			},
		}},
		OutputList: []*node.TypeOutput{&node.TypeOutput{
			PortCommonOut: &node.PortCommonOut{
				Name: "out1",
				Type: "int8",
				Connections: []*node.Connection{&node.Connection{
					NodeID:   "",
					NodePort: "in1",
				}},
			},
		}},
	})

	nodeB, _ := nodes.New(&nodes.Node{
		InputList: []*node.TypeInput{&node.TypeInput{
			PortCommon: &node.PortCommon{
				Name: "in1",
				Type: "int8",
				Connection: &node.Connection{
					NodeID:   "",
					NodePort: "out1",
				},
			},
		}},
		OutputList: []*node.TypeOutput{&node.TypeOutput{
			PortCommonOut: &node.PortCommonOut{
				Name: "out1",
				Type: "int8",
				Connections: []*node.Connection{&node.Connection{
					NodeID:   "",
					NodePort: "",
				}},
			},
		}},
	})

	var nodesList []*nodes.Node

	nodesList = append(nodesList, nodeA)
	nodesList = append(nodesList, nodeB)
	pprint.PrintJOSN(nodesList)
}
