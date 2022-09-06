package back2

import (
	"fmt"
	"github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/_example/nodes"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type genericNode struct {
	NodeName string      `json:"nodeName"`
	NodeType string      `json:"nodeType"`
	Node     interface{} `json:"node"`
}

//type nodesList struct {
//	Node []genericNode `json:"nodes"`
//}

func main() {

	nodeA := nodes.NewPass("nodeA")

	nodeB := nodes.NewPass("nodeB")

	nameA_ := &genericNode{
		NodeName: "nodeA",
		NodeType: "PASS",
		Node:     nodeA,
	}

	nameB_ := &genericNode{
		NodeName: "nodeB",
		NodeType: "PASS",
		Node:     nodeB,
	}
	//////
	var nodesParsed []*genericNode
	////
	nodesParsed = append(nodesParsed, nameA_)
	nodesParsed = append(nodesParsed, nameB_)

	pprint.PrintJOSN(nodesParsed)
	//var nodesParsed []*genericNode
	jsonFile, err := os.Open("../flowctrl/_example/test.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	//
	//// we initialize our Users array
	////nodesParsed []*genericNode
	//
	//// we unmarshal our byteArray which contains our
	//// jsonFile's content into 'users' which we defined above
	//json.Unmarshal(byteValue, &nodesParsed)
	//pprint.PrintJOSN(nodesParsed)
	//value := gjson.ParseBytes(byteValue)

	//aa := value.Value().(*nodes.Pass)
	//pprint.Print(aa)

	graph := flowctrl.New()
	for _, n := range nodesParsed { // init nodes
		switch n.NodeType {
		case "PASS":
			//pprint.Print(n.Node)
			//aa := n.Node.(*nodes.Pass)
			//pprint.Print(aa)
			//var pass *nodes.Pass
			//err := mapstructure.Decode(n.Node, &pass)
			//
			//pprint.Print(pass)
			//
			//if pass.Info().Name == "nodeA" {
			//	fmt.Printf("Pointer addr GetNode: %p\n", pass)
			//}
			//if err != nil {
			//	return
			//}

			//nodeA := nodes.NewPass(n.NodeName)
			//fmt.Println("ADD NEW NODE", n.NodeName)
			//fmt.Printf("Pointer addr: %p\n", nodeA)
			graph.AddNode(nodes.NewPass(n.NodeName))

		}
	}

	for _, n := range nodesParsed { // init nodes
		switch n.NodeType {
		case "PASS":
			//fmt.Println("UPDATE NODE")
			//var pass *nodes.Pass
			//
			//err := mapstructure.Decode(n.Node, &pass)
			//fmt.Println(err)
			//if err != nil {
			//	return
			//}
			////value := gjson.Get(json, "name.last")
			//pprint.PrintJOSN(pass)
			nodeAa := graph.GetNode(n.NodeName).(*nodes.Pass)
			if nodeAa.Info().Name == "nodeA" {
				fmt.Println(n.NodeName, "SHOW A")
				fmt.Printf("Pointer addr GetNode:  %p\n", nodeAa)
				//initWriter := adapter.NewInt8(nodeA.Input)
				//initWriter.Set(123)
			}

			if nodeAa.Info().Name == "nodeA" {
				nodeBb := graph.GetNode("nodeB").(*nodes.Pass)
				fmt.Println(11111, nodeBb.Info().Name)
				nodeAa.Output.Connect(nodeBb.Input)
				if nodeAa.Info().Name == "nodeA" {
					graph.ReplaceNode("nodeA", nodeAa)
				}
				if nodeAa.Info().Name == "nodeB" {
					graph.ReplaceNode("nodeB", nodeBb)
				}

			}

		}
	}

	for _, ordered := range graph.Get().Graphs {
		for _, runner := range ordered.Runners {
			//fmt.Println("RUNNER-1", runner.Name(), len(runner.Outputs()), "LEN")

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
	//log.Println("Final read:", endReader.Get())
}
