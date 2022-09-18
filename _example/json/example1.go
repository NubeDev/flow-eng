package main

import (
	"encoding/json"
	"flag"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/helpers/mqttbase"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes"
	"github.com/NubeDev/flow-eng/storage"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"time"
)

func main() {

	storage.New("")

	filePath := flag.String("f", "../flow-eng/_example/json/bacnet.json", "flow file")
	flag.Parse()
	fmt.Println("file:", *filePath)

	var nodesParsed []*node.Spec
	jsonFile, err := os.Open(*filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	sort.Slice(nodesParsed, func(i, j int) bool {
		return false
	})
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &nodesParsed)

	graph := flowctrl.New()

	m, err := mqttbase.NewMqtt()
	if err != nil {
		return
	}

	m.Connect()

	if m.Connected() {
		m.Publish("start bacnet", "test")
	}
	fmt.Println(4444, m)

	for _, n := range nodesParsed {
		node_, err := nodes.Builder(n, m)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("ADD:", node_.GetName(), node_.GetNodeName(), "ERR", err)
		graph.AddNode(node_)
	}

	graph.ReBuildFlow(true)

	for _, n := range graph.GetNodes() {
		fmt.Println("GET NODES:", n.GetName(), n.GetNodeName())
	}

	runner := flowctrl.NewSerialRunner(graph)
	// pprint.PrintJOSN(graph.GetNodes())

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
