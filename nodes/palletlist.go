package nodes

import (
	"errors"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/compare"
	debugging "github.com/NubeDev/flow-eng/nodes/debug"
	"github.com/NubeDev/flow-eng/nodes/functions"
	"github.com/NubeDev/flow-eng/nodes/logic"
	"github.com/NubeDev/flow-eng/nodes/math"
	broker "github.com/NubeDev/flow-eng/nodes/mqtt"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	"github.com/NubeDev/flow-eng/nodes/statistics"
	"github.com/NubeDev/flow-eng/nodes/timing"
)

const (
	disableMQTT = true // example now how to disable a node from the user being able to add it, will be moved to the config file
)

func All() []*node.Spec { // get all the nodes, will be used for the UI to list all the nodes
	// math
	constNum, _ := math.NewConst(nil)
	add, _ := math.NewAdd(nil)
	sub, _ := math.NewSub(nil)
	multiply, _ := math.NewMultiply(nil)
	divide, _ := math.NewDivide(nil)

	// bool
	and, _ := logic.NewAnd(nil)
	or, _ := logic.NewOr(nil)

	// compare
	comp, _ := compare.NewCompare(nil)
	between, _ := compare.NewBetween(nil)

	// compare
	min, _ := statistics.NewMin(nil)
	max, _ := statistics.NewMax(nil)

	// time
	delay, _ := timing.NewDelay(nil, nil)
	inject, _ := timing.NewInject(nil)

	funcNode, _ := functions.NewFunc(nil)

	// bacnet
	bacServer, _ := bacnet.NewServer(nil)
	bacPointAI, _ := bacnet.NewAI(nil)
	bacPointAO, _ := bacnet.NewAO(nil)
	bacPointBV, _ := bacnet.NewBV(nil)

	// pointbus
	mqttSub, _ := broker.NewMqttSub(nil)
	mqttPub, _ := broker.NewMqttPub(nil)

	logNode, _ := debugging.NewLog(nil)

	if disableMQTT {
		mqttSub = nil
		mqttPub = nil
	}

	return node.BuildNodes(
		node.ConvertToSpec(constNum),
		node.ConvertToSpec(add),
		node.ConvertToSpec(sub),
		node.ConvertToSpec(multiply),
		node.ConvertToSpec(divide),

		node.ConvertToSpec(and),
		node.ConvertToSpec(or),

		node.ConvertToSpec(comp),
		node.ConvertToSpec(between),

		node.ConvertToSpec(min),
		node.ConvertToSpec(max),

		node.ConvertToSpec(delay),
		node.ConvertToSpec(inject),

		node.ConvertToSpec(funcNode),

		node.ConvertToSpec(bacServer),
		node.ConvertToSpec(bacPointAI),
		node.ConvertToSpec(bacPointAO),
		node.ConvertToSpec(bacPointBV),

		node.ConvertToSpec(mqttSub),
		node.ConvertToSpec(mqttPub),

		node.ConvertToSpec(logNode),
	)
}

func Builder(body *node.Spec, opts ...interface{}) (node.Node, error) {
	n, err := builderMath(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderLogic(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderCompare(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderStatistics(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderTiming(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderMisc(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderProtocols(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderMQTT(body)
	if n != nil || err != nil {
		return n, err
	}
	return nil, errors.New(fmt.Sprintf("no nodes found with name:%s", body.GetName()))
}

func builderMisc(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case logNode:
		return debugging.NewLog(body)
	case funcNode:
		return functions.NewFunc(body)
	}
	return nil, nil
}

func builderMath(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case constNum:
		return math.NewConst(body)
	case add:
		return math.NewAdd(body)
	case sub:
		return math.NewSub(body)
	case multiply:
		return math.NewMultiply(body)
	case divide:
		return math.NewDivide(body)
	}
	return nil, nil
}

func builderLogic(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case and:
		return logic.NewAnd(body)
	case or:
		return logic.NewOr(body)
	}
	return nil, nil
}

func builderCompare(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case logicCompare:
		return compare.NewCompare(body)
	case between:
		return compare.NewBetween(body)
	}
	return nil, nil
}

func builderStatistics(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case min:
		return statistics.NewMin(body)
	case max:
		return statistics.NewMax(body)
	}
	return nil, nil
}

func builderTiming(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case delay:
		return timing.NewDelay(body, flowctrl.NewTimer())
	case inject:
		return timing.NewInject(body)
	}
	return nil, nil
}

func builderProtocols(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case bacnetServer:
		return bacnet.NewServer(body)
	case bacnetAI:
		return bacnet.NewAI(body)
	case bacnetAO:
		return bacnet.NewAO(body)
	case bacnetBV:
		return bacnet.NewBV(body)

	}

	return nil, nil
}

func builderMQTT(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case mqttSub:
		return broker.NewMqttSub(body)
	case mqttPub:
		return broker.NewMqttPub(body)
	}
	return nil, nil
}
