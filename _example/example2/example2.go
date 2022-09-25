package main

import (
	flowctrl "github.com/NubeDev/flow-eng"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/constant"
	"github.com/NubeDev/flow-eng/nodes/math"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {

	const1, err := constant.NewConst(nil) // new node
	if err != nil {
		log.Errorln(err)
		return
	}

	err = const1.OverrideInputValue(node.In, 11.0)
	if err != nil {
		log.Errorln(err)
		return
	}
	const2, err := constant.NewConst(nil)
	if err != nil {
		log.Errorln(err)
		return
	}
	err = const2.OverrideInputValue(node.In, 33.3)
	if err != nil {
		log.Errorln(err)
		return
	}
	add, err := math.NewAdd(nil) // new math (add) node
	if err != nil {
		log.Errorln(err)
		return
	}

	graph := flowctrl.New() // init the flow engine

	graph.AddNode(const1) // add the nodes to the runtime
	graph.AddNode(const2)
	graph.AddNode(add)

	// graph.AddNode(mqttSub)
	// graph.AddNode(mqttPub)

	err = graph.ManualNodeConnector(const1, node.Out, add, node.InputA) // connect const-1 and 2 to the add node
	if err != nil {
		log.Errorln(err)
		return
	}

	err = graph.ManualNodeConnector(const2, node.Out, add, node.InputB)
	if err != nil {
		log.Errorln(err)
		return
	}

	graph.ReBuildFlow(true)

	runner := flowctrl.NewSerialRunner(graph) // make the runner for lopping
	pprint.PrintJOSN(graph.GetNodesSpec())

	for {
		err := runner.Process()
		// random := float.RandFloat(0, 1)
		// err = const2.OverrideInputValue(node.In1, random)
		if err != nil {
			log.Errorln(err)
			return
		}
		time.Sleep(1000 * time.Millisecond) // change loop time
	}
}
