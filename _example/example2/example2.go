package main

import (
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/_example/nodes/math"
	broker "github.com/NubeDev/flow-eng/_example/nodes/mqtt"
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {

	const1, err := math.NewConst(nil) // new node
	if err != nil {
		log.Errorln(err)
		return
	}

	err = const1.OverrideInputValue(node.In1, 11)
	if err != nil {
		log.Errorln(err)
		return
	}
	const2, err := math.NewConst(nil)
	if err != nil {
		log.Errorln(err)
		return
	}
	err = const2.OverrideInputValue(node.In1, 11)
	if err != nil {
		log.Errorln(err)
		return
	}
	add, err := math.NewAdd(nil) // new math (add) node
	if err != nil {
		log.Errorln(err)
		return
	}

	mqttSub, err := broker.NewMqttSub(nil) // new math (add) node
	if err != nil {
		log.Errorln(err)
		return
	}

	mqttPub, err := broker.NewMqttPub(nil) // new math (add) node
	if err != nil {
		log.Errorln(err)
		return
	}

	graph := flowctrl.New() // init the flow engine

	graph.AddNode(const1) // add the nodes to the runtime
	graph.AddNode(const2)
	graph.AddNode(add)
	graph.AddNode(mqttSub)
	graph.AddNode(mqttPub)

	err = graph.ManualNodeConnector(const1, add, node.Out1, node.In1) // connect const-1 and 2 to the add node
	if err != nil {
		log.Errorln(err)
		return
	}

	err = graph.ManualNodeConnector(const2, add, node.Out1, node.In2)
	if err != nil {
		log.Errorln(err)
		return
	}
	err = graph.ManualNodeConnector(mqttSub, add, node.Out1, node.In3)
	if err != nil {
		log.Errorln(err)
		return
	}

	err = graph.ManualNodeConnector(add, mqttPub, node.Out1, node.In1)
	if err != nil {
		log.Errorln(err)
		return
	}

	graph.ReplaceNode(const1) // add the nodes to the runtime
	graph.ReplaceNode(const2)
	graph.ReplaceNode(add)
	graph.ReplaceNode(mqttSub)
	graph.ReplaceNode(mqttPub)

	runner := flowctrl.NewSerialRunner(graph) // make the runner for lopping

	for {
		err := runner.Process()
		//random := float.RandFloat(0, 1)
		//err = const2.OverrideInputValue(node.In1, random)
		if err != nil {
			log.Errorln(err)
			return
		}
		time.Sleep(1000 * time.Millisecond) // change loop time
	}
}
