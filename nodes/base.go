package nodes

import (
	"errors"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/node"
	math2 "github.com/NubeDev/flow-eng/nodes/math"
	broker2 "github.com/NubeDev/flow-eng/nodes/mqtt"
	timing2 "github.com/NubeDev/flow-eng/nodes/timing"
)

func All() []node.Node { // get all the nodes, will be used for the UI to list all the nodes
	// math
	constNum, _ := math2.NewConst(nil)
	add, _ := math2.NewAdd(nil)
	sub, _ := math2.NewSub(nil)
	// time
	delay, _ := timing2.NewDelay(nil, nil)
	inject, _ := timing2.NewInject(nil)
	// mqtt
	mqttSub, _ := broker2.NewMqttSub(nil)
	mqttPub, _ := broker2.NewMqttPub(nil)
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
		return math2.NewConst(body)
	case add:
		return math2.NewAdd(body)
	case sub:
		return math2.NewSub(body)
	}
	return nil, nil
}

func builderTiming(body *node.BaseNode) (node.Node, error) {
	switch body.GetName() {
	case delay:
		return timing2.NewDelay(body, flowctrl.NewTimer())
	case inject:
		return timing2.NewInject(body)
	}
	return nil, nil
}

func builderMQTT(body *node.BaseNode) (node.Node, error) {
	switch body.GetName() {
	case mqttSub:
		return broker2.NewMqttSub(body)
	case mqttPub:
		return broker2.NewMqttPub(body)
	}
	return nil, nil
}
