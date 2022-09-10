package main

import (
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/_example/nodes/math"
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {

	const1, err := math.NewConst(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = const1.OverrideInputValue(node.In1, 11)
	if err != nil {
		fmt.Println(err)
		return
	}
	const2, err := math.NewConst(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = const2.OverrideInputValue(node.In1, 22)
	if err != nil {
		fmt.Println(err)
		return
	}
	add, err := math.NewAdd(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	graph := flowctrl.New()
	err = graph.ManualNodeConnector(const1, add, node.Out1, node.In1)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = graph.ManualNodeConnector(const2, add, node.Out1, node.In2)
	if err != nil {
		fmt.Println(err)
		return
	}

	graph.AddNode(const1)
	graph.AddNode(const2)
	graph.AddNode(add)

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
