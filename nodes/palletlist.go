package nodes

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/nodes/notify"
	"github.com/NubeDev/flow-eng/nodes/trigger"

	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/helpers/store"
	"github.com/NubeDev/flow-eng/helpers/timer"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/boolean"
	"github.com/NubeDev/flow-eng/nodes/compare"
	"github.com/NubeDev/flow-eng/nodes/constant"
	"github.com/NubeDev/flow-eng/nodes/conversion"
	"github.com/NubeDev/flow-eng/nodes/count"
	debugging "github.com/NubeDev/flow-eng/nodes/debug"
	"github.com/NubeDev/flow-eng/nodes/filter"
	"github.com/NubeDev/flow-eng/nodes/functions"
	"github.com/NubeDev/flow-eng/nodes/hvac"
	nodejson "github.com/NubeDev/flow-eng/nodes/json"
	"github.com/NubeDev/flow-eng/nodes/latch"
	"github.com/NubeDev/flow-eng/nodes/link"
	"github.com/NubeDev/flow-eng/nodes/math"
	"github.com/NubeDev/flow-eng/nodes/mathematics"
	broker "github.com/NubeDev/flow-eng/nodes/mqtt"
	"github.com/NubeDev/flow-eng/nodes/numtransform"
	point "github.com/NubeDev/flow-eng/nodes/points"
	bacnetio "github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	"github.com/NubeDev/flow-eng/nodes/protocols/flow"
	"github.com/NubeDev/flow-eng/nodes/protocols/rest"
	"github.com/NubeDev/flow-eng/nodes/statistics"
	"github.com/NubeDev/flow-eng/nodes/streams"
	"github.com/NubeDev/flow-eng/nodes/subflow"
	switches "github.com/NubeDev/flow-eng/nodes/switch"
	"github.com/NubeDev/flow-eng/nodes/system"
	"github.com/NubeDev/flow-eng/nodes/timing"
)

const (
	disableNodes = true
)

func All() []*node.Spec { // get all the nodes, will be used for the UI to list all the nodes
	constBool, _ := constant.NewBoolean(nil)
	constNum, _ := constant.NewNumber(nil)
	constStr, _ := constant.NewString(nil)

	// boolean
	and, _ := boolean.NewAnd(nil)
	or, _ := boolean.NewOr(nil)
	xor, _ := boolean.NewXor(nil)
	not, _ := boolean.NewNot(nil)
	toggle, _ := boolean.NewToggle(nil)

	// compare
	greaterthan, _ := compare.NewGreaterThan(nil)
	lessthan, _ := compare.NewLessThan(nil)
	equal, _ := compare.NewEqual(nil)
	equalString, _ := compare.NewEqualString(nil)
	between, _ := compare.NewBetween(nil)
	hysteresis, _ := compare.NewHysteresis(nil)

	// statistics
	min, _ := statistics.NewMin(nil)
	max, _ := statistics.NewMax(nil)
	avg, _ := statistics.NewAvg(nil)
	minMaxAvg, _ := statistics.NewMinMaxAvg(nil)
	rangeNode, _ := statistics.NewRange(nil)
	rank, _ := statistics.NewRank(nil)
	median, _ := statistics.NewMedian(nil)

	// streams
	flatLine, _ := streams.NewFlatline(nil)

	flowNetwork, _ := flow.NewNetwork(nil)
	flowPoint, _ := flow.NewFFPoint(nil)
	flowPointWrite, _ := flow.NewFFPointWrite(nil)
	flowSchedule, _ := flow.NewFFSchedule(nil)

	flowLoopCount, _ := system.NewLoopCount(nil)
	subFlowFolder, _ := subflow.NewSubFlowFolder(nil)
	subInputFloat, _ := subflow.NewSubFlowInputFloat(nil)
	subInputString, _ := subflow.NewSubFlowInputString(nil)
	subInputBool, _ := subflow.NewSubFlowInputBool(nil)

	subOutputFloat, _ := subflow.NewSubFlowOutputFloat(nil)
	subOutputBool, _ := subflow.NewSubFlowOutputBool(nil)
	subOutputString, _ := subflow.NewSubFlowOutputString(nil)

	conversionString, _ := conversion.NewString(nil)
	conversionNum, _ := conversion.NewNumber(nil)
	conversionBool, _ := conversion.NewBoolean(nil)

	funcNode, _ := functions.NewFunc(nil)

	jsonFilter, _ := nodejson.NewFilter(nil)
	dataStore, _ := nodejson.NewStore(nil)

	// hvac
	accumPeriod, _ := hvac.NewAccumulationPeriod(nil)
	deadBand, _ := hvac.NewDeadBand(nil)
	leadLagSwitch, _ := hvac.NewLeadLagSwitch(nil)
	pid, _ := hvac.NewPIDNode(nil)
	pacControl, _ := hvac.NewPACControl(nil)
	sensorSelect, _ := hvac.NewSensorSelect(nil)
	psychroDBRH, _ := hvac.NewPsychroDBRH(nil)
	psychroDBDP, _ := hvac.NewPsychroDBDP(nil)
	psychroDBWB, _ := hvac.NewPsychroDBWB(nil)

	gmail, _ := notify.NewGmail(nil)
	pingNode, _ := notify.NewPing(nil)

	// latch
	numLatch, _ := latch.NewNumLatch(nil)
	stringLatch, _ := latch.NewStringLatch(nil)
	setResetLatch, _ := latch.NewSetResetLatch(nil)

	numOutputSelect, _ := switches.NewNumOutputSelect(nil)
	selectNode, _ := switches.NewSelectNum(nil)
	switchNode, _ := switches.NewSwitch(nil)

	linkInput, _ := link.NewStringLinkInput(nil, nil)
	linkInputNum, _ := link.NewNumLinkInput(nil, nil)
	linkInputBool, _ := link.NewBoolLinkInput(nil, nil)
	linkOutput, _ := link.NewStringLinkOutput(nil, nil)
	linkOutputNum, _ := link.NewNumLinkOutput(nil, nil)
	linkOutputBool, _ := link.NewBoolLinkOutput(nil, nil)

	// math
	abs, _ := math.NewAbsolute(nil)
	add, _ := math.NewAdd(nil)
	sub, _ := math.NewSub(nil)
	multiply, _ := math.NewMultiply(nil)
	divide, _ := math.NewDivide(nil)
	modulo, _ := math.NewModulo(nil)
	ceiling, _ := math.NewCeiling(nil)
	floor, _ := math.NewFloor(nil)

	mathAdvanced, _ := mathematics.NewAdvanced(nil)

	// trigger
	countNode, _ := count.NewCount(nil)
	countNum, _ := count.NewCountNum(nil)
	countString, _ := count.NewCountString(nil)

	// triggers
	cov, _ := trigger.NewCOVNode(nil)
	random, _ := trigger.NewRandom(nil)
	iterate, _ := trigger.NewIterate(nil)

	// time
	delay, _ := timing.NewDelay(nil)
	delayOn, _ := timing.NewDelayOn(nil, nil)
	delayOff, _ := timing.NewDelayOff(nil, nil)
	dutyCycle, _ := timing.NewDutyCycle(nil)
	minOnOff, _ := timing.NewMinOnOff(nil)
	oneShot, _ := timing.NewOneShot(nil)
	stopwatch, _ := timing.NewStopwatch(nil)

	// numtransform
	fade, _ := numtransform.NewFade(nil)
	limitNode, _ := numtransform.NewLimit(nil)
	rateLimit, _ := numtransform.NewRateLimit(nil)
	scaleNode, _ := numtransform.NewScale(nil)
	waveform, _ := numtransform.NewWaveform(nil)

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

	boolWriteable, _ := point.NewBooleanWriteable(nil)
	numWriteable, _ := point.NewNumericWriteable(nil)
	stringWriteable, _ := point.NewStringWriteable(nil)

	// filter
	preventNull, _ := filter.NewPreventNull(nil)
	preventEqualFloat, _ := filter.NewPreventEqualFloat(nil)
	preventEqualString, _ := filter.NewPreventEqualString(nil)
	onlyBetween, _ := filter.NewOnlyBetween(nil)
	onlyGreater, _ := filter.NewOnlyGreater(nil)
	onlyLower, _ := filter.NewOnlyLower(nil)
	preventDuplicates, _ := filter.NewPreventDuplicates(nil)

	// rest
	getNode, _ := rest.NewGet(nil)
	writeNode, _ := rest.NewHttpWrite(nil)

	if disableNodes {
		dataStore = nil
		jsonFilter = nil
		pingNode = nil
		getNode = nil
		writeNode = nil
		funcNode = nil
		mqttBroker = nil
		mqttSub = nil
		mqttPub = nil
	}

	return node.BuildNodes(

		node.ConvertToSpec(bacServer),
		node.ConvertToSpec(bacPointAI),
		node.ConvertToSpec(bacPointAO),
		node.ConvertToSpec(bacPointAV),
		node.ConvertToSpec(bacPointBV),

		node.ConvertToSpec(flowNetwork),
		node.ConvertToSpec(flowPoint),
		node.ConvertToSpec(flowPointWrite),
		node.ConvertToSpec(flowSchedule),

		node.ConvertToSpec(getNode),
		node.ConvertToSpec(writeNode),

		node.ConvertToSpec(mqttBroker),
		node.ConvertToSpec(mqttSub),
		node.ConvertToSpec(mqttPub),

		node.ConvertToSpec(and),
		node.ConvertToSpec(or),
		node.ConvertToSpec(xor),
		node.ConvertToSpec(not),
		node.ConvertToSpec(toggle),

		node.ConvertToSpec(greaterthan),
		node.ConvertToSpec(lessthan),
		node.ConvertToSpec(equal),
		node.ConvertToSpec(equalString),
		node.ConvertToSpec(between),
		node.ConvertToSpec(hysteresis),

		node.ConvertToSpec(constBool),
		node.ConvertToSpec(constNum),
		node.ConvertToSpec(constStr),

		node.ConvertToSpec(conversionString),
		node.ConvertToSpec(conversionNum),
		node.ConvertToSpec(conversionBool),

		node.ConvertToSpec(countNode),
		node.ConvertToSpec(countNum),
		node.ConvertToSpec(countString),

		node.ConvertToSpec(logNode),

		node.ConvertToSpec(preventNull),
		node.ConvertToSpec(preventEqualFloat),
		node.ConvertToSpec(preventEqualString),
		node.ConvertToSpec(onlyBetween),
		node.ConvertToSpec(onlyGreater),
		node.ConvertToSpec(onlyLower),
		node.ConvertToSpec(preventDuplicates),

		node.ConvertToSpec(funcNode),

		node.ConvertToSpec(accumPeriod),
		node.ConvertToSpec(deadBand),
		node.ConvertToSpec(leadLagSwitch),
		node.ConvertToSpec(pid),
		node.ConvertToSpec(pacControl),
		node.ConvertToSpec(sensorSelect),
		node.ConvertToSpec(psychroDBRH),
		node.ConvertToSpec(psychroDBDP),
		node.ConvertToSpec(psychroDBWB),

		node.ConvertToSpec(jsonFilter),
		node.ConvertToSpec(dataStore),

		node.ConvertToSpec(numLatch),
		node.ConvertToSpec(stringLatch),
		node.ConvertToSpec(setResetLatch),

		node.ConvertToSpec(linkInputNum),
		node.ConvertToSpec(linkInputBool),
		node.ConvertToSpec(linkOutputNum),
		node.ConvertToSpec(linkOutputBool),
		node.ConvertToSpec(linkInput),
		node.ConvertToSpec(linkOutput),

		node.ConvertToSpec(abs),
		node.ConvertToSpec(add),
		node.ConvertToSpec(sub),
		node.ConvertToSpec(multiply),
		node.ConvertToSpec(divide),
		node.ConvertToSpec(modulo),
		node.ConvertToSpec(ceiling),
		node.ConvertToSpec(floor),
		node.ConvertToSpec(mathAdvanced),

		node.ConvertToSpec(fade),
		node.ConvertToSpec(limitNode),
		node.ConvertToSpec(rateLimit),
		node.ConvertToSpec(scaleNode),
		node.ConvertToSpec(waveform),

		node.ConvertToSpec(gmail),
		node.ConvertToSpec(pingNode),

		node.ConvertToSpec(boolWriteable),
		node.ConvertToSpec(numWriteable),
		node.ConvertToSpec(stringWriteable),

		node.ConvertToSpec(min),
		node.ConvertToSpec(max),
		node.ConvertToSpec(avg),
		node.ConvertToSpec(minMaxAvg),
		node.ConvertToSpec(rangeNode),
		node.ConvertToSpec(rank),
		node.ConvertToSpec(median),

		node.ConvertToSpec(flatLine),

		node.ConvertToSpec(flowLoopCount),
		node.ConvertToSpec(subFlowFolder),
		node.ConvertToSpec(subInputFloat),
		node.ConvertToSpec(subInputString),
		node.ConvertToSpec(subInputBool),
		node.ConvertToSpec(subOutputFloat),
		node.ConvertToSpec(subOutputBool),
		node.ConvertToSpec(subOutputString),

		node.ConvertToSpec(numOutputSelect),
		node.ConvertToSpec(switchNode),
		node.ConvertToSpec(selectNode),

		node.ConvertToSpec(delay),
		node.ConvertToSpec(delayOn),
		node.ConvertToSpec(delayOff),
		node.ConvertToSpec(dutyCycle),
		node.ConvertToSpec(minOnOff),
		node.ConvertToSpec(oneShot),
		node.ConvertToSpec(stopwatch),

		node.ConvertToSpec(cov),
		node.ConvertToSpec(random),
		node.ConvertToSpec(iterate),
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
	n, err = builderFilter(body)
	if n != nil || err != nil {
		return n, err
	}
	return nil, errors.New(fmt.Sprintf("no nodes found with name:%s", body.GetName()))
}

func builderMisc(body *node.Spec) (node.Node, error) {
	con := &link.Store{}
	switch body.GetName() {
	case subFlowFolder:
		return subflow.NewSubFlowFolder(body)
	case inputFloat:
		return subflow.NewSubFlowInputFloat(body)
	case inputString:
		return subflow.NewSubFlowInputString(body)
	case inputBool:
		return subflow.NewSubFlowInputBool(body)
	case outputFloat:
		return subflow.NewSubFlowOutputFloat(body)
	case outputBool:
		return subflow.NewSubFlowOutputBool(body)
	case outputString:
		return subflow.NewSubFlowOutputString(body)
	case logNode:
		return debugging.NewLog(body)
	case funcNode:
		return functions.NewFunc(body)
	case linkInput:
		return link.NewStringLinkInput(body, con)
	case linkOutput:
		return link.NewStringLinkOutput(body, con)
	case linkInputNum:
		return link.NewNumLinkInput(body, con)
	case linkOutputNum:
		return link.NewNumLinkOutput(body, con)
	case linkInputBool:
		return link.NewBoolLinkInput(body, con)
	case linkOutputBool:
		return link.NewBoolLinkOutput(body, con)
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

func builderFilter(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	// case onlyTrue:
	// 	return filter.NewOnlyTrue(body)
	// case onlyFalse:
	// 	return filter.NewOnlyFalse(body)
	case preventNull:
		return filter.NewPreventNull(body)
	case preventEqualFloat:
		return filter.NewPreventEqualFloat(body)
	case preventEqualString:
		return filter.NewPreventEqualString(body)
	case onlyBetween:
		return filter.NewOnlyBetween(body)
	case onlyGreater:
		return filter.NewOnlyGreater(body)
	case onlyLower:
		return filter.NewOnlyLower(body)
	case preventDuplicates:
		return filter.NewPreventDuplicates(body)
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
	case gmail:
		return notify.NewGmail(body)
	case pingNode:
		return notify.NewPing(body)
	}
	return nil, nil
}

func builderTransformations(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case limitNode:
		return numtransform.NewLimit(body)
	case scaleNode:
		return numtransform.NewScale(body)
	case fade:
		return numtransform.NewFade(body)
	case rateLimit:
		return numtransform.NewRateLimit(body)
	case waveform:
		return numtransform.NewWaveform(body)
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
	case flowSchedule:
		return flow.NewFFSchedule(body)
	}
	return nil, nil
}

func builderHVAC(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case accumulationPeriod:
		return hvac.NewAccumulationPeriod(body)
	case deadBandNode:
		return hvac.NewDeadBand(body)
	case leadLagSwitch:
		return hvac.NewLeadLagSwitch(body)
	case pidNode:
		return hvac.NewPIDNode(body)
	case pacControlNode:
		return hvac.NewPACControl(body)
	case sensorSelect:
		return hvac.NewSensorSelect(body)
	case psychroDBRH:
		return hvac.NewPsychroDBRH(body)
	case psychroDBDP:
		return hvac.NewPsychroDBDP(body)
	case psychroDBWB:
		return hvac.NewPsychroDBWB(body)
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
	case abs:
		return math.NewAbsolute(body)
	case add:
		return math.NewAdd(body)
	case sub:
		return math.NewSub(body)
	case multiply:
		return math.NewMultiply(body)
	case divide:
		return math.NewDivide(body)
	case modulo:
		return math.NewModulo(body)
	case ceiling:
		return math.NewCeiling(body)
	case floor:
		return math.NewFloor(body)
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
	case numOutputSelect:
		return switches.NewNumOutputSelect(body)
	}
	return nil, nil
}

func builderBoolean(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case and:
		return boolean.NewAnd(body)
	case or:
		return boolean.NewOr(body)
	case xor:
		return boolean.NewXor(body)
	case not:
		return boolean.NewNot(body)
	case toggle:
		return boolean.NewToggle(body)
	}
	return nil, nil
}

func builderCompare(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case greaterThan:
		return compare.NewGreaterThan(body)
	case lessThan:
		return compare.NewLessThan(body)
	case equal:
		return compare.NewEqual(body)
	case equalString:
		return compare.NewEqualString(body)
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
	case minMaxAvg:
		return statistics.NewMinMaxAvg(body)
	case rangeNode:
		return statistics.NewRange(body)
	case rank:
		return statistics.NewRank(body)
	case median:
		return statistics.NewMedian(body)
	}
	return nil, nil
}

func builderPoints(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case boolWriteable:
		return point.NewBooleanWriteable(body)
	case numWriteable:
		return point.NewNumericWriteable(body)
	case stringWriteable:
		return point.NewStringWriteable(body)
	}
	return nil, nil
}

func builderCount(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case countNode:
		return count.NewCount(body)
	case countNumNode:
		return count.NewCountNum(body)
	case countStringNode:
		return count.NewCountString(body)
	}
	return nil, nil
}

func builderTrigger(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case covNode:
		return trigger.NewCOVNode(body)
	case random:
		return trigger.NewRandom(body)
	case iterate:
		return trigger.NewIterate(body)
	}
	return nil, nil
}

func builderTiming(body *node.Spec) (node.Node, error) {
	switch body.GetName() {
	case delay:
		return timing.NewDelay(body)
	case delayOn:
		return timing.NewDelayOn(body, timer.NewTimer())
	case delayOff:
		return timing.NewDelayOff(body, timer.NewTimer())
	case dutyCycle:
		return timing.NewDutyCycle(body)
	case minOnOff:
		return timing.NewMinOnOff(body)
	case oneShot:
		return timing.NewOneShot(body)
	case stopwatch:
		return timing.NewStopwatch(body)
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
