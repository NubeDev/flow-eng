package main

import (
	"encoding/json"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/_example/nodes"
	"github.com/NubeDev/flow-eng/node"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {
	buildJson()

	var nodesParsed []*node.Node
	jsonFile, err := os.Open("../flow-eng/_example/test.json")
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &nodesParsed)

	graph := flowctrl.New()
	for _, n := range nodesParsed {
		node_, _ := nodes.New(n)
		graph.AddNode(node_)
	}

	// START >>> <<<====JUST FOR TEST ====>>>
	getA := graph.GetNode("pass_a6160d20")
	getB := graph.GetNode("pass_23e61419")
	for _, output := range getA.GetOutputs() {
		for _, input := range getB.GetInputs() {
			output.OutputPort.Connect(input.InputPort)
		}
	}
	graph.ReplaceNode("pass_a6160d20", getA)
	graph.ReplaceNode("pass_23e61419", getB)
	// END >>> <<<====JUST FOR TEST ====>>>

	for _, ordered := range graph.Get().Graphs {
		for _, runner := range ordered.Runners {
			// fmt.Println("RUNNER-1", runner.Name(), len(runner.Outputs()), "LEN", "UUID", runner.UUID())
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
		fmt.Println(err)
		if err != nil {
			panic(err)
		}
		time.Sleep(5 * time.Second)
	}
}

func buildJson() {
	nodeA, _ := nodes.New(&node.Node{
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

	nodeB, _ := nodes.New(&node.Node{
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

	var nodesList []interface{}

	nodesList = append(nodesList, nodeA)
	nodesList = append(nodesList, nodeB)
}
