package nodes

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/helpers/store"
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/bool"
	"github.com/NubeDev/flow-eng/nodes/compare"
	"github.com/NubeDev/flow-eng/nodes/constant"
	"github.com/NubeDev/flow-eng/nodes/conversion"
	"github.com/NubeDev/flow-eng/nodes/count"
	debugging "github.com/NubeDev/flow-eng/nodes/debug"
	"github.com/NubeDev/flow-eng/nodes/functions"
	"github.com/NubeDev/flow-eng/nodes/hvac"
	nodejson "github.com/NubeDev/flow-eng/nodes/json"
	"github.com/NubeDev/flow-eng/nodes/latch"
	"github.com/NubeDev/flow-eng/nodes/link"
	"github.com/NubeDev/flow-eng/nodes/math"
	"github.com/NubeDev/flow-eng/nodes/mathematics"
	broker "github.com/NubeDev/flow-eng/nodes/mqtt"
	"github.com/NubeDev/flow-eng/nodes/notify/gmail"
	"github.com/NubeDev/flow-eng/nodes/notify/ping"
	point "github.com/NubeDev/flow-eng/nodes/points"
	bacnetio "github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	"github.com/NubeDev/flow-eng/nodes/protocols/flow"
	"github.com/NubeDev/flow-eng/nodes/protocols/rest"
	"github.com/NubeDev/flow-eng/nodes/statistics"
	"github.com/NubeDev/flow-eng/nodes/streams"
	switches "github.com/NubeDev/flow-eng/nodes/switch"
	"github.com/NubeDev/flow-eng/nodes/system"
	"github.com/NubeDev/flow-eng/nodes/tigger"
	"github.com/NubeDev/flow-eng/nodes/timing"
	"github.com/NubeDev/flow-eng/nodes/transformations"
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
	delayMinOnOff, _ := bool.NewMinOnOff(nil, nil)
	// compare
	comp, _ := compare.NewCompareGreater(nil)
	compLess, _ := compare.NewCompareLess(nil)
	logicCompareEqual, _ := compare.NewCompareEqual(nil)
	between, _ := compare.NewBetween(nil)
	hysteresis, _ := compare.NewHysteresis(nil)

	// statistics
	min, _ := statistics.NewMin(nil)
	max, _ := statistics.NewMax(nil)
	avg, _ := statistics.NewAvg(nil)

	// streams
	flatLine, _ := streams.NewFlatline(nil)

	flowNetwork, _ := flow.NewNetwork(nil)
	flowPoint, _ := flow.NewFFPoint(nil)
	flowPointWrite, _ := flow.NewFFPointWrite(nil)

	flowLoopCount, _ := system.NewLoopCount(nil)

	conversionString, _ := conversion.NewString(nil)
	conversionNum, _ := conversion.NewNumber(nil)
	conversionBool, _ := conversion.NewBoolean(nil)

	funcNode, _ := functions.NewFunc(nil)

	jsonFilter, _ := nodejson.NewFilter(nil)
	dataStore, _ := nodejson.NewStore(nil)

	// hvac
	deadBand, _ := hvac.NewDeadBand(nil)
	pid, _ := hvac.NewPIDNode(nil)
	pacControl, _ := hvac.NewPACControl(nil)

	gmailNode, _ := gmail.NewGmail(nil)
	pingNode, _ := ping.NewPing(nil)

	// latch
	numLatch, _ := latch.NewNumLatch(nil)
	stringLatch, _ := latch.NewStringLatch(nil)
	setResetLatch, _ := latch.NewSetResetLatch(nil)

	selectNode, _ := switches.NewSelectNum(nil)
	switchNode, _ := switches.NewSwitch(nil)

	linkInput, _ := link.NewInput(nil, nil)
	linkInputNum, _ := link.NewInputNum(nil, nil)
	linkOutput, _ := link.NewOutput(nil, nil)
	linkOutputNum, _ := link.NewOutputNum(nil, nil)

	// math
	add, _ := math.NewAdd(nil)
	sub, _ := math.NewSub(nil)
	multiply, _ := math.NewMultiply(nil)
	divide, _ := math.NewDivide(nil)

	mathAdvanced, _ := mathematics.NewAdvanced(nil)

	// trigger
	countBool, _ := count.NewCountBool(nil)
	countNum, _ := count.NewCountNum(nil)
	countString, _ := count.NewCountString(nil)
	rampNode, _ := count.NewRamp(nil)

	// trigger
	inject, _ := trigger.NewInject(nil)
	randomFloat, _ := trigger.NewRandom(nil)

	// time
	delay, _ := timing.NewDelay(nil, nil)
	delayOn, _ := timing.NewDelayOn(nil, nil)
	delayOff, _ := timing.NewDelayOff(nil, nil)

	// number transformations
	scaleNode, _ := transformations.NewScale(nil)
	limitNode, _ := transformations.NewLimit(nil)

	// bacnet
	bacServer, _ := bacnetio.NewServer(nil, nil)
	bacPointAI, _ := bacnetio.NewAI(nil, nil)
	bacPointAO, _ := bacnetio.NewAO(nil, nil)
	bacPointAV, _ := bacnetio.NewAV(nil, nil)
	bacPointBV, _ := bacnetio.NewBV(nil, nil)

	mqttBroker, _ := broker.NewBroker(nil)
	mqttSub, _ := broker.NewMqttSub(nil)
	mqttPub, _ := broker.NewMqttPub(nil)

	logNode, _ := debugging.NewLog(nil)

	pointNum, _ := point.NewNumber(nil)
	pointBool, _ := point.NewBoolean(nil)

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
		node.ConvertToSpec(mathAdvanced),

		node.ConvertToSpec(and),
		node.ConvertToSpec(or),
		node.ConvertToSpec(xor),
		node.ConvertToSpec(not),
		node.ConvertToSpec(toggle),
		node.ConvertToSpec(delayMinOnOff),

		node.ConvertToSpec(comp),
		node.ConvertToSpec(compLess),
		node.ConvertToSpec(logicCompareEqual),
		node.ConvertToSpec(between),
		node.ConvertToSpec(hysteresis),

		node.ConvertToSpec(min),
		node.ConvertToSpec(max),
		node.ConvertToSpec(avg),

		node.ConvertToSpec(gmailNode),
		node.ConvertToSpec(pingNode),

		node.ConvertToSpec(flowLoopCount),

		node.ConvertToSpec(flowNetwork),
		node.ConvertToSpec(flowPoint),
		node.ConvertToSpec(flowPointWrite),

		node.ConvertToSpec(deadBand),
		node.ConvertToSpec(pid),
		node.ConvertToSpec(pacControl),

		node.ConvertToSpec(pointNum),
		node.ConvertToSpec(pointBool),

		node.ConvertToSpec(jsonFilter),
		node.ConvertToSpec(dataStore),

		node.ConvertToSpec(numLatch),
		node.ConvertToSpec(stringLatch),
		node.ConvertToSpec(setResetLatch),

		node.ConvertToSpec(delay),
		node.ConvertToSpec(delayOn),
		node.ConvertToSpec(delayOff),

		node.ConvertToSpec(randomFloat),
		node.ConvertToSpec(inject),

		node.ConvertToSpec(flatLine),

		node.ConvertToSpec(countBool),
		node.ConvertToSpec(countNum),
		node.ConvertToSpec(countString),

		node.ConvertToSpec(rampNode),

		node.ConvertToSpec(funcNode),

		node.ConvertToSpec(switchNode),
		node.ConvertToSpec(selectNode),

		node.ConvertToSpec(conversionString),
		node.ConvertToSpec(conversionNum),
		node.ConvertToSpec(conversionBool),

		node.ConvertToSpec(linkInput),
		node.ConvertToSpec(linkInputNum),
		node.ConvertToSpec(linkOutput),
		node.ConvertToSpec(linkOutputNum),

		node.ConvertToSpec(bacServer),
		node.ConvertToSpec(bacPointAI),
		node.ConvertToSpec(bacPointAO),
		node.ConvertToSpec(bacPointAV),
		node.ConvertToSpec(bacPointBV),

		node.ConvertToSpec(mqttBroker),
		node.ConvertToSpec(mqttSub),
		node.ConvertToSpec(mqttPub),

		node.ConvertToSpec(scaleNode),
		node.ConvertToSpec(limitNode),

		node.ConvertToSpec(logNode),

		node.ConvertToSpec(getNode),
		node.ConvertToSpec(writeNode),
	)
}

func Builder(body *node.Spec, db db.DB, store *store.Store, opts ...interface{}) (node.Node, error) {
	body.AddDB(db)
	body.AddStore(store)
	n, err := builderConst(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderMath(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderFlowNetworks(body, opts)
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
	n, err = builderTransformations(body)
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
	n, err = builderPoints(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderTiming(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderTrigger(body)
	if n != nil || err != nil {
		return n, err
	}
	n, err = builderCount(body)
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
	n, err = builderProtocols(body, opts)
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
	case linkInputNum:
		return link.NewInputNum(body, con)
	case linkOutputNum:
		return link.NewOutputNum(body, con)
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

func builderTransformations(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case limitNode:
		return transformations.NewLimit(body)
	case scaleNode:
		return transformations.NewScale(body)
	}
	return nil, nil
}

func builderJson(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case jsonFilter:
		return nodejson.NewFilter(body)
	case dataStore:
		return nodejson.NewStore(body)
	}
	return nil, nil
}

func builderFlowNetworks(body *node.Spec, opts []interface{}) (node.Node, error) {
	switch body.GetName() {
	case flowNetwork:
		return flow.NewNetwork(body)
	case flowPoint:
		return flow.NewFFPoint(body)
	case flowPointWrite:
		return flow.NewFFPointWrite(body)
	}
	return nil, nil
}

func builderHVAC(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case deadBandNode:
		return hvac.NewDeadBand(body)
	case pidNode:
		return hvac.NewPIDNode(body)
	case pacControlNode:
		return hvac.NewPACControl(body)
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
	case mathAdvanced:
		return mathematics.NewAdvanced(body)
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
		return bool.NewMinOnOff(body, timer.NewTimer())
	}
	return nil, nil
}

func builderCompare(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case logicCompareGreater:
		return compare.NewCompareGreater(body)
	case logicCompareLess:
		return compare.NewCompareLess(body)
	case logicCompareEqual:
		return compare.NewCompareEqual(body)
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
	case avg:
		return statistics.NewAvg(body)
	}
	return nil, nil
}

func builderPoints(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case pointNumber:
		return point.NewNumber(body)
	case pointBoolean:
		return point.NewBoolean(body)
	}
	return nil, nil
}

func builderCount(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case countBoolNode:
		return count.NewCountBool(body)
	case countNumNode:
		return count.NewCountNum(body)
	case countStringNode:
		return count.NewCountString(body)
	case rampNode:
		return count.NewRamp(body)
	}
	return nil, nil
}

func builderTrigger(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case inject:
		return trigger.NewInject(body)
	case randomFloat:
		return trigger.NewRandom(body)
	}
	return nil, nil
}

func builderTiming(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case delay:
		return timing.NewDelay(body, timer.NewTimer())
	case delayOn:
		return timing.NewDelayOn(body, timer.NewTimer())
	case delayOff:
		return timing.NewDelayOff(body, timer.NewTimer())
	}
	return nil, nil
}

func builderProtocols(body *node.Spec, opts []interface{}) (node.Node, error) {
	bacOpts := &bacnetio.Bacnet{}
	if len(opts) > 0 {
		_, ok := opts[0].(*bacnetio.Bacnet)
		if ok {
			bacOpts = opts[0].(*bacnetio.Bacnet)
		}
	}
	switch body.GetName() {
	case bacnetServer:
		return bacnetio.NewServer(body, bacOpts)
	case bacnetAI:
		return bacnetio.NewAI(body, bacOpts)
	case bacnetAO:
		return bacnetio.NewAO(body, bacOpts)
	case bacnetAV:
		return bacnetio.NewAV(body, bacOpts)
	case bacnetBV:
		return bacnetio.NewBV(body, bacOpts)

	}
	return nil, nil
}

func builderMQTT(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case mqttBroker:
		return broker.NewBroker(body)
	case mqttSub:
		return broker.NewMqttSub(body)
	case mqttPub:
		return broker.NewMqttPub(body)
	}
	return nil, nil
}
