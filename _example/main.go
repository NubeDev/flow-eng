package main

import (
	"encoding/json"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
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
		node_, _ := node.New(n)
		graph.AddNode(node_)
	}

	// START >>> <<<====JUST FOR TEST ====>>>
	getA := graph.GetNode("a123")
	getB := graph.GetNode("b123")
	for _, output := range getA.GetOutputs() {
		for _, input := range getB.GetInputs() {
			output.Connect(input.InputPort)
		}
	}
	graph.ReplaceNode("a123", getA)
	graph.ReplaceNode("b123", getB)
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
		time.Sleep(5000 * time.Millisecond)
	}
}

func buildJson() {

	//var nodeA *node.Node
	var count = []string{"nodeA", "nodeB"}
	var nodesList []interface{}
	nodeA := &node.Node{
		Inputs: []*node.Input{&node.Input{
			InputPort: &node.InputPort{
				Name:     "",
				DataType: "",
				Connection: &node.Connection{
					NodeID:   "",
					NodePort: "",
				},
			},
		}},
		Outputs: []*node.Output{&node.Output{
			OutputPort: &node.OutputPort{
				Name:     "",
				DataType: "",
				Connections: []*node.Connection{&node.Connection{
					NodeID:   "",
					NodePort: "",
				},
				},
			},
			ValueFloat64: nil,
			ValueString:  nil,
		}},
	}
	for _, name := range count {
		fmt.Println(name)
		var node node.Node
		node.Info.Name = name
		node = *nodeA

		nodesList = append(nodesList, node)
	}

	//pprint.PrintJOSN(nodesList)
}
