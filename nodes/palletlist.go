package nodes

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/compare"
	"github.com/NubeDev/flow-eng/nodes/connection"
	"github.com/NubeDev/flow-eng/nodes/constant"
	"github.com/NubeDev/flow-eng/nodes/conversion"
	debugging "github.com/NubeDev/flow-eng/nodes/debug"
	"github.com/NubeDev/flow-eng/nodes/functions"
	"github.com/NubeDev/flow-eng/nodes/hvac"
	"github.com/NubeDev/flow-eng/nodes/logic"
	"github.com/NubeDev/flow-eng/nodes/math"
	broker "github.com/NubeDev/flow-eng/nodes/mqtt"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/nodes/statistics"
	switches "github.com/NubeDev/flow-eng/nodes/switch"
	"github.com/NubeDev/flow-eng/nodes/system"
	"github.com/NubeDev/flow-eng/nodes/timing"
)

const (
	disableMQTT = true // example now how to disable a node from the user being able to add it, will be moved to the config file
)

func All() []*node.Spec { // get all the nodes, will be used for the UI to list all the nodes
	constNum, _ := constant.NewConstNum(nil)
	constStr, _ := constant.NewString(nil)

	// math
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
	hysteresis, _ := compare.NewHysteresis(nil)

	// compare
	min, _ := statistics.NewMin(nil)
	max, _ := statistics.NewMax(nil)

	flowLoopCount, _ := system.NewLoopCount(nil)

	stringToNum, _ := conversion.NewStringToNum(nil)
	numToString, _ := conversion.NewNumToString(nil)

	// time
	delay, _ := timing.NewDelay(nil, nil)
	inject, _ := timing.NewInject(nil)
	delayOn, _ := timing.NewDelayOn(nil, nil)

	funcNode, _ := functions.NewFunc(nil)

	// hvac
	deadBand, _ := hvac.NewDeadBand(nil)

	selectNode, _ := switches.NewSelectNum(nil)

	connectionInput, _ := connection.NewInput(nil, nil)
	connectionOutput, _ := connection.NewOutput(nil, nil)

	// bacnet
	bacServer, _ := bacnet.NewServer(nil, nil)
	bacPointAI, _ := bacnet.NewAI(nil, nil)
	bacPointAO, _ := bacnet.NewAO(nil, nil)
	bacPointAV, _ := bacnet.NewAV(nil, nil)
	bacPointBV, _ := bacnet.NewBV(nil, nil)

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
		node.ConvertToSpec(constStr),

		node.ConvertToSpec(add),
		node.ConvertToSpec(sub),
		node.ConvertToSpec(multiply),
		node.ConvertToSpec(divide),

		node.ConvertToSpec(and),
		node.ConvertToSpec(or),

		node.ConvertToSpec(comp),
		node.ConvertToSpec(between),
		node.ConvertToSpec(hysteresis),

		node.ConvertToSpec(min),
		node.ConvertToSpec(max),

		node.ConvertToSpec(flowLoopCount),

		node.ConvertToSpec(deadBand),

		node.ConvertToSpec(delay),
		node.ConvertToSpec(inject),
		node.ConvertToSpec(delayOn),

		node.ConvertToSpec(funcNode),

		node.ConvertToSpec(selectNode),

		node.ConvertToSpec(stringToNum),
		node.ConvertToSpec(numToString),

		node.ConvertToSpec(connectionInput),
		node.ConvertToSpec(connectionOutput),

		node.ConvertToSpec(bacServer),
		node.ConvertToSpec(bacPointAI),
		node.ConvertToSpec(bacPointAO),
		node.ConvertToSpec(bacPointAV),
		node.ConvertToSpec(bacPointBV),

		node.ConvertToSpec(mqttSub),
		node.ConvertToSpec(mqttPub),

		node.ConvertToSpec(logNode),
	)
}

func Builder(body *node.Spec, opts ...interface{}) (node.Node, error) {
	n, err := builderConst(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderMath(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderHVAC(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderSystem(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderLogic(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderConversion(body)
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
	n, err = builderSwitch(body)
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
	con := &connection.Store{}
	switch body.GetName() {
	case logNode:
		return debugging.NewLog(body)
	case funcNode:
		return functions.NewFunc(body)
	case connectionInput:
		return connection.NewInput(body, con)
	case connectionOutput:
		return connection.NewOutput(body, con)
	}

	return nil, nil
}

func builderSystem(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case flowLoopCount:
		return system.NewLoopCount(body)
	}
	return nil, nil
}

func builderHVAC(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case constNum:
		return hvac.NewDeadBand(body)
	}
	return nil, nil
}

func builderConst(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case constNum:
		return constant.NewConstNum(body)
	case constStr:
		return constant.NewString(body)
	}
	return nil, nil
}

func builderMath(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
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

func builderConversion(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case numToString:
		return conversion.NewNumToString(body)
	case stringToNum:
		return conversion.NewStringToNum(body)
	}
	return nil, nil
}

func builderSwitch(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case selectNum:
		return switches.NewSelectNum(body)
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
	case hysteresis:
		return compare.NewHysteresis(body)
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
		return timing.NewDelay(body, timer.NewTimer())
	case inject:
		return timing.NewInject(body)
	case delayOn:
		return timing.NewDelayOn(body, timer.NewTimer())
	}
	return nil, nil
}

func builderProtocols(body *node.Spec) (node.Node, error) {
	store := points.New(applications.RubixIO, nil, 0, 200, 200)
	switch body.GetName() {
	case bacnetServer:
		return bacnet.NewServer(body, store)
	case bacnetAI:
		return bacnet.NewAI(body, store)
	case bacnetAO:
		return bacnet.NewAO(body, store)
	case bacnetAV:
		return bacnet.NewAV(body, store)
	case bacnetBV:
		return bacnet.NewBV(body, store)

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
