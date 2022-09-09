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

	var nodesParsed []*node.BaseNode
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
		node_, err := nodes.Builder(n)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("ADD:", node_.GetName(), node_.GetNodeName(), "ERR", err)
		graph.AddNode(node_)
	}

	for _, n := range graph.GetNodes() {
		fmt.Println("*******build connections:", n.GetName(), n.GetNodeName())
		err := graph.NodeConnector(n.GetID())
		fmt.Println("build connections", err)
		if err != nil {
			return
		}
	}

	for _, n := range graph.GetNodes() {
		fmt.Println("REPLACE", n.GetName(), n.GetNodeName())
		graph.ReplaceNode(n.GetID(), n)
	}
	//pprint.PrintJOSN(graph.GetNodes())
	runner := flowctrl.NewSerialRunner(graph)

	log.Println("Flow started")
	for {
		err := runner.Process()
		fmt.Println(err)
		if err != nil {
			panic(err)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func buildJson() {
	// var nodeA *node.BaseNode
	var count = []string{"nodeA", "nodeB"}
	var nodesList []interface{}
	nodeA := &node.BaseNode{
		Inputs: []*node.Input{&node.Input{
			InputPort: &node.InputPort{
				Name:     "in1",
				DataType: "",
				Connection: &node.InputConnection{
					NodeID:   "aa",
					NodePort: "aaa",
				},
			},
		}},
		Outputs: []*node.Output{&node.Output{
			OutputPort: &node.OutputPort{
				Name:     "out1",
				DataType: "",
				Connections: []*node.OutputConnection{&node.OutputConnection{
					NodeID:   "bbb",
					NodePort: "34dfg",
				},
				},
			},
			Value: nil,
		}},
	}

	// pprint.PrintJOSN(a)
	for _, name := range count {

		var node node.BaseNode
		node.Info.Name = name
		node = *nodeA

		nodesList = append(nodesList, node)
	}
}
