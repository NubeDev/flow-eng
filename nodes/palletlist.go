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
	nodejson "github.com/NubeDev/flow-eng/nodes/json"
	"github.com/NubeDev/flow-eng/nodes/latch"
	"github.com/NubeDev/flow-eng/nodes/link"
	"github.com/NubeDev/flow-eng/nodes/math"
	broker "github.com/NubeDev/flow-eng/nodes/mqtt"
	"github.com/NubeDev/flow-eng/nodes/notify/gmail"
	"github.com/NubeDev/flow-eng/nodes/notify/ping"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/nodes/protocols/driver"
	"github.com/NubeDev/flow-eng/nodes/protocols/flow"
	"github.com/NubeDev/flow-eng/nodes/protocols/rest"
	"github.com/NubeDev/flow-eng/nodes/statistics"
	"github.com/NubeDev/flow-eng/nodes/streams"
	switches "github.com/NubeDev/flow-eng/nodes/switch"
	"github.com/NubeDev/flow-eng/nodes/system"
	"github.com/NubeDev/flow-eng/nodes/timing"
)

const (
	disableMQTT = false // example now how to disable a node from the user being able to add it, will be moved to the config file
)

func All() []*node.Spec { // get all the nodes, will be used for the UI to list all the nodes
	constNum, _ := constant.NewNumber(nil)
	constBool, _ := constant.NewBoolean(nil)
	constStr, _ := constant.NewString(nil)

	// bool
	and, _ := bool.NewAnd(nil)
	or, _ := bool.NewOr(nil)
	xor, _ := bool.NewXor(nil)
	not, _ := bool.NewNot(nil)
	toggle, _ := bool.NewToggle(nil)
	delayMinOnOff, _ := bool.NewMinOn(nil, nil)
	// compare
	comp, _ := compare.NewCompare(nil)
	between, _ := compare.NewBetween(nil)
	hysteresis, _ := compare.NewHysteresis(nil)

	// statistics
	min, _ := statistics.NewMin(nil)
	max, _ := statistics.NewMax(nil)

	// streams
	flatline, _ := streams.NewFlatline(nil)

	flowNetwork, _ := flow.NewNetwork(nil, nil)
	flowDevice, _ := flow.NewDevice(nil, nil)
	flowPoint, _ := flow.NewPoint(nil, nil)

	flowLoopCount, _ := system.NewLoopCount(nil)

	conversionString, _ := conversion.NewString(nil)
	conversionNum, _ := conversion.NewNumber(nil)
	conversionBool, _ := conversion.NewBoolean(nil)

	funcNode, _ := functions.NewFunc(nil)

	jsonFilter, _ := nodejson.NewFilter(nil)

	// hvac
	deadBand, _ := hvac.NewDeadBand(nil)

	gmailNode, _ := gmail.NewGmail(nil)
	pingNode, _ := ping.NewPing(nil)

	// latch
	numLatch, _ := latch.NewNumLatch(nil)
	stringLatch, _ := latch.NewStringLatch(nil)
	setResetLatch, _ := latch.NewSetResetLatch(nil)

	selectNode, _ := switches.NewSelectNum(nil)
	switchNode, _ := switches.NewSwitch(nil)

	linkInput, _ := link.NewInput(nil, nil)
	linkOutput, _ := link.NewOutput(nil, nil)

	// math
	add, _ := math.NewAdd(nil)
	sub, _ := math.NewSub(nil)
	multiply, _ := math.NewMultiply(nil)
	divide, _ := math.NewDivide(nil)

	// time
	delay, _ := timing.NewDelay(nil, nil)

	inject, _ := timing.NewInject(nil)
	delayOn, _ := timing.NewDelayOn(nil, nil)

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

	// rest
	getNode, _ := rest.NewGet(nil)
	writeNode, _ := rest.NewHttpWrite(nil)

	// if disableMQTT {
	//	mqttSub = nil
	//	mqttPub = nil
	// }

	return node.BuildNodes(
		node.ConvertToSpec(constNum),
		node.ConvertToSpec(constBool),
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
		node.ConvertToSpec(delayMinOnOff),

		node.ConvertToSpec(comp),
		node.ConvertToSpec(between),
		node.ConvertToSpec(hysteresis),

		node.ConvertToSpec(min),
		node.ConvertToSpec(max),

		node.ConvertToSpec(gmailNode),
		node.ConvertToSpec(pingNode),

		node.ConvertToSpec(flowLoopCount),

		node.ConvertToSpec(flowNetwork),
		node.ConvertToSpec(flowDevice),
		node.ConvertToSpec(flowPoint),

		node.ConvertToSpec(deadBand),

		node.ConvertToSpec(jsonFilter),

		node.ConvertToSpec(numLatch),
		node.ConvertToSpec(stringLatch),
		node.ConvertToSpec(setResetLatch),

		node.ConvertToSpec(delay),
		node.ConvertToSpec(inject),
		node.ConvertToSpec(delayOn),

		node.ConvertToSpec(flatline),

		node.ConvertToSpec(funcNode),

		node.ConvertToSpec(switchNode),
		node.ConvertToSpec(selectNode),

		node.ConvertToSpec(conversionString),
		node.ConvertToSpec(conversionNum),
		node.ConvertToSpec(conversionBool),

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

		node.ConvertToSpec(getNode),
		node.ConvertToSpec(writeNode),
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
	n, err = builderJson(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderRest(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderNotify(body)
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
	n, err = builderStreams(body)
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

func builderRest(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case getHttpNode:
		return rest.NewGet(body)
	case writeHttpNode:
		return rest.NewHttpWrite(body)
	}
	return nil, nil
}

func builderNotify(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case gmailNode:
		return gmail.NewGmail(body)
	case pingNode:
		return ping.NewPing(body)
	}
	return nil, nil
}

func builderJson(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case jsonFilter:
		return nodejson.NewFilter(body)
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

func builderStreams(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case flatLine:
		return streams.NewFlatline(body)
	}
	return nil, nil
}

func builderConst(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case constNum:
		return constant.NewNumber(body)
	case constStr:
		return constant.NewString(body)
	case constBool:
		return constant.NewBoolean(body)
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
	case conversionString:
		return conversion.NewString(body)
	case conversionNum:
		return conversion.NewNumber(body)
	case conversionBool:
		return conversion.NewBoolean(body)
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
	case delayMinOnOff:
		return bool.NewMinOn(body, timer.NewTimer())
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
	store := points.New(names.Edge, nil, 0, 200, 200)
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
