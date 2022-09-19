package main

import (
	flowctrl "github.com/NubeDev/flow-eng"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/nodes/math"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {

	const1, err := math.NewConst(nil) // new node
	if err != nil {
		log.Errorln(err)
		return
	}

	const2, err := math.NewConst(nil)
	if err != nil {
		log.Errorln(err)
		return
	}

	graph := flowctrl.New() // init the flow engine

	graph.AddNode(const1) // add the nodes to the runtime
	graph.AddNode(const2)

	// graph.AddNode(mqttSub)
	// graph.AddNode(mqttPub)

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
