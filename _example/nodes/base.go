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

func All() []node.Node {
	add, _ := math.NewAdd(nil)
	sub, _ := math.NewSub(nil)
	return node.BuildNodes(add, sub)
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
	case add:
		return math.NewAdd(body)
	case sub:
		return math.NewSub(body)
	case delay:
		return timing.NewDelay(body, flowctrl.NewTimer())
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
