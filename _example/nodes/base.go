package nodes

import (
	"errors"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/_example/nodes/math"
	broker "github.com/NubeDev/flow-eng/_example/nodes/mqtt"
	"github.com/NubeDev/flow-eng/_example/nodes/timing"
	"github.com/NubeDev/flow-eng/node"
)

func All() []node.Node { // get all the nodes, will be used for the UI to list all the nodes
	// math
	constNum, _ := math.NewConst(nil)
	add, _ := math.NewAdd(nil)
	sub, _ := math.NewSub(nil)
	// time
	delay, _ := timing.NewDelay(nil, nil)
	inject, _ := timing.NewInject(nil)
	// mqtt
	mqttSub, _ := broker.NewMqttSub(nil)
	mqttPub, _ := broker.NewMqttPub(nil)
	return node.BuildNodes(
		constNum,
		add,
		sub,
		delay,
		inject,
		mqttSub,
		mqttPub,
	)
}

func Builder(body *node.BaseNode) (node.Node, error) {
	n, err := builderMath(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderTiming(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderMQTT(body)
	if n != nil || err != nil {
		return n, err
	}
	return nil, errors.New(fmt.Sprintf("no nodes found with name:%s", body.GetName()))
}

func builderMath(body *node.BaseNode) (node.Node, error) {
	switch body.GetName() {
	case constNum:
		return math.NewConst(body)
	case add:
		return math.NewAdd(body)
	case sub:
		return math.NewSub(body)
	}
	return nil, nil
}

func builderTiming(body *node.BaseNode) (node.Node, error) {
	switch body.GetName() {
	case delay:
		return timing.NewDelay(body, flowctrl.NewTimer())
	case inject:
		return timing.NewInject(body)
	}
	return nil, nil
}

func builderMQTT(body *node.BaseNode) (node.Node, error) {
	switch body.GetName() {
	case mqttSub:
		return broker.NewMqttSub(body)
	case mqttPub:
		return broker.NewMqttPub(body)
	}
	return nil, nil
}
