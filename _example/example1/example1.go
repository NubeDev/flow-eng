package main

//
import (
	"encoding/json"
	"flag"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/_example/nodes"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/storage"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {

	storage.New("")

	filePath := flag.String("f", "../flow-eng/_example/example1/mqtt.json", "flow file")
	flag.Parse()
	fmt.Println("file:", *filePath)

	var nodesParsed []*node.BaseNode
	jsonFile, err := os.Open(*filePath)
	if err != nil {
		fmt.Println(err)
		return
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
		fmt.Println("BUILD connections:", n.GetName(), n.GetNodeName())
		err := graph.NodeConnector(n.GetID())

		if err != nil {
			fmt.Println("build connections ERROR", err)
			return
		}
	}

	for _, nn := range graph.GetNodes() {
		fmt.Println("REPLACE", nn.GetName(), nn.GetNodeName(), nn.GetID())
		graph.ReplaceNode(nn.GetID(), nn)

	}

	for _, n := range graph.GetNodes() {
		fmt.Println("GET NODES:", n.GetName(), n.GetNodeName())
	}

	runner := flowctrl.NewSerialRunner(graph)

	log.Println("Flow started")
	for {
		err := runner.Process()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
