package nodes

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/bool"
	"github.com/NubeDev/flow-eng/nodes/compare"
	"github.com/NubeDev/flow-eng/nodes/constant"
	"github.com/NubeDev/flow-eng/nodes/conversion"
	debugging "github.com/NubeDev/flow-eng/nodes/debug"
	"github.com/NubeDev/flow-eng/nodes/functions"
	"github.com/NubeDev/flow-eng/nodes/hvac"
	"github.com/NubeDev/flow-eng/nodes/latch"
	"github.com/NubeDev/flow-eng/nodes/link"
	"github.com/NubeDev/flow-eng/nodes/math"
	broker "github.com/NubeDev/flow-eng/nodes/mqtt"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/nodes/protocols/driver"
	"github.com/NubeDev/flow-eng/nodes/protocols/flow"
	"github.com/NubeDev/flow-eng/nodes/statistics"
	switches "github.com/NubeDev/flow-eng/nodes/switch"
	"github.com/NubeDev/flow-eng/nodes/system"
	"github.com/NubeDev/flow-eng/nodes/timing"
)

const (
	disableMQTT = false // example now how to disable a node from the user being able to add it, will be moved to the config file
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
	and, _ := bool.NewAnd(nil)
	or, _ := bool.NewOr(nil)
	xor, _ := bool.NewXor(nil)
	not, _ := bool.NewNot(nil)
	toggle, _ := bool.NewToggle(nil)

	// compare
	comp, _ := compare.NewCompare(nil)
	between, _ := compare.NewBetween(nil)
	hysteresis, _ := compare.NewHysteresis(nil)

	// compare
	min, _ := statistics.NewMin(nil)
	max, _ := statistics.NewMax(nil)

	flowNetwork, _ := flow.NewNetwork(nil, nil)
	flowDevice, _ := flow.NewDevice(nil, nil)
	flowPoint, _ := flow.NewPoint(nil, nil)

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

	// latch
	numLatch, _ := latch.NewNumLatch(nil)
	stringLatch, _ := latch.NewStringLatch(nil)
	setResetLatch, _ := latch.NewSetResetLatch(nil)

	selectNode, _ := switches.NewSelectNum(nil)
	switchNode, _ := switches.NewSwitch(nil)

	linkInput, _ := link.NewInput(nil, nil)
	linkOutput, _ := link.NewOutput(nil, nil)

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

	// if disableMQTT {
	//	mqttSub = nil
	//	mqttPub = nil
	// }

	return node.BuildNodes(
		node.ConvertToSpec(constNum),
		node.ConvertToSpec(constStr),

		node.ConvertToSpec(add),
		node.ConvertToSpec(sub),
		node.ConvertToSpec(multiply),
		node.ConvertToSpec(divide),

		node.ConvertToSpec(and),
		node.ConvertToSpec(or),
		node.ConvertToSpec(xor),
		node.ConvertToSpec(not),
		node.ConvertToSpec(toggle),

		node.ConvertToSpec(comp),
		node.ConvertToSpec(between),
		node.ConvertToSpec(hysteresis),

		node.ConvertToSpec(min),
		node.ConvertToSpec(max),

		node.ConvertToSpec(flowLoopCount),

		node.ConvertToSpec(flowNetwork),
		node.ConvertToSpec(flowDevice),
		node.ConvertToSpec(flowPoint),

		node.ConvertToSpec(deadBand),

		node.ConvertToSpec(numLatch),
		node.ConvertToSpec(stringLatch),
		node.ConvertToSpec(setResetLatch),

		node.ConvertToSpec(delay),
		node.ConvertToSpec(inject),
		node.ConvertToSpec(delayOn),

		node.ConvertToSpec(funcNode),

		node.ConvertToSpec(switchNode),
		node.ConvertToSpec(selectNode),

		node.ConvertToSpec(stringToNum),
		node.ConvertToSpec(numToString),

		node.ConvertToSpec(linkInput),
		node.ConvertToSpec(linkOutput),

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

func Builder(body *node.Spec, db db.DB, opts ...interface{}) (node.Node, error) {
	body.AddDB(db)
	n, err := builderConst(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderMath(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderFlowNetworks(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderHVAC(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderLatch(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderSystem(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderBoolean(body)
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
	con := &link.Store{}
	switch body.GetName() {
	case logNode:
		return debugging.NewLog(body)
	case funcNode:
		return functions.NewFunc(body)
	case linkInput:
		return link.NewInput(body, con)
	case linkOutput:
		return link.NewOutput(body, con)
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

func builderFlowNetworks(body *node.Spec) (node.Node, error) {
	networksPool := driver.New(&driver.Networks{})
	switch body.GetName() {
	case flowNetwork:
		return flow.NewNetwork(body, networksPool)
	case flowDevice:
		return flow.NewDevice(body, networksPool)
	case flowPoint:
		return flow.NewPoint(body, networksPool)
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

func builderLatch(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case numLatch:
		return latch.NewNumLatch(body)
	case stringLatch:
		return latch.NewStringLatch(body)
	case setResetLatch:
		return latch.NewSetResetLatch(body)
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
	case switchNode:
		return switches.NewSwitch(body)
	case selectNum:
		return switches.NewSelectNum(body)
	}
	return nil, nil
}

func builderBoolean(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case and:
		return bool.NewAnd(body)
	case or:
		return bool.NewOr(body)
	case xor:
		return bool.NewXor(body)
	case not:
		return bool.NewNot(body)
	case toggle:
		return bool.NewToggle(body)
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
	store := points.New(names.RubixIO, nil, 0, 200, 200)
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
