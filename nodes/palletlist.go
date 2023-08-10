package nodes

import (
	"github.com/NubeDev/flow-eng/nodes/notify"
	switches "github.com/NubeDev/flow-eng/nodes/switch"
	"github.com/NubeDev/flow-eng/nodes/trigger"
	"github.com/NubeDev/flow-eng/pallet"

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

	"github.com/NubeDev/flow-eng/nodes/system"
	"github.com/NubeDev/flow-eng/nodes/timing"
)

const (
	disableNodes = true
)

//     if disableNodes {
//         dataStore = nil
//         jsonFilter = nil
//         pingNode = nil
//         getNode = nil
//         writeNode = nil
//         funcNode = nil
//         mqttBroker = nil
//         mqttSub = nil
//         mqttPub = nil
//         logNode = nil
//     }

func RegisterAllNodes() {
	// builderMisc
	pallet.RegisterNodeBuilder(subflow.Category, subFlowFolder, subflow.NewSubFlowFolder)
	pallet.RegisterNodeBuilder(subflow.Category, inputFloat, subflow.NewSubFlowInputFloat)
	pallet.RegisterNodeBuilder(subflow.Category, inputString, subflow.NewSubFlowInputString)
	pallet.RegisterNodeBuilder(subflow.Category, inputBool, subflow.NewSubFlowInputBool)
	pallet.RegisterNodeBuilder(subflow.Category, outputFloat, subflow.NewSubFlowOutputFloat)
	pallet.RegisterNodeBuilder(subflow.Category, outputBool, subflow.NewSubFlowOutputBool)
	pallet.RegisterNodeBuilder(subflow.Category, outputString, subflow.NewSubFlowOutputString)
	pallet.RegisterNodeBuilder(debugging.Category, logNode, debugging.NewLog)
	pallet.RegisterNodeBuilder(functions.Category, funcNode, functions.NewFunc)
	pallet.RegisterNodeBuilder(link.Category, linkInput, link.NewStringLinkInput)
	pallet.RegisterNodeBuilder(link.Category, linkOutput, link.NewStringLinkOutput)
	pallet.RegisterNodeBuilder(link.Category, linkInputNum, link.NewNumLinkInput)
	pallet.RegisterNodeBuilder(link.Category, linkOutputNum, link.NewNumLinkOutput)
	pallet.RegisterNodeBuilder(link.Category, linkInputBool, link.NewBoolLinkInput)
	pallet.RegisterNodeBuilder(link.Category, linkOutputBool, link.NewBoolLinkOutput)

	// builderSystem
	pallet.RegisterNodeBuilder(system.Category, flowLoopCount, system.NewLoopCount)

	// builderFilter
	pallet.RegisterNodeBuilder(filter.Category, preventNull, filter.NewPreventNull)
	pallet.RegisterNodeBuilder(filter.Category, preventEqualFloat, filter.NewPreventEqualFloat)
	pallet.RegisterNodeBuilder(filter.Category, preventEqualString, filter.NewPreventEqualString)
	pallet.RegisterNodeBuilder(filter.Category, onlyBetween, filter.NewOnlyBetween)
	pallet.RegisterNodeBuilder(filter.Category, onlyGreater, filter.NewOnlyGreater)
	pallet.RegisterNodeBuilder(filter.Category, onlyLower, filter.NewOnlyLower)
	pallet.RegisterNodeBuilder(filter.Category, preventDuplicates, filter.NewPreventDuplicates)

	// builderRest
	pallet.RegisterNodeBuilder(rest.Category, getHttpNode, rest.NewGet)
	pallet.RegisterNodeBuilder(rest.Category, writeHttpNode, rest.NewHttpWrite)

	// builderNotify
	pallet.RegisterNodeBuilder(notify.Category, gmail, notify.NewGmail)
	pallet.RegisterNodeBuilder(notify.Category, pingNode, notify.NewPing)

	// builderTransformations
	pallet.RegisterNodeBuilder(numtransform.Category, limitNode, numtransform.NewLimit)
	pallet.RegisterNodeBuilder(numtransform.Category, scaleNode, numtransform.NewScale)
	pallet.RegisterNodeBuilder(numtransform.Category, fade, numtransform.NewFade)
	pallet.RegisterNodeBuilder(numtransform.Category, rateLimit, numtransform.NewRateLimit)
	pallet.RegisterNodeBuilder(numtransform.Category, round, numtransform.NewRound)
	pallet.RegisterNodeBuilder(numtransform.Category, waveform, numtransform.NewWaveform)
	pallet.RegisterNodeBuilder(numtransform.Category, polynomial, numtransform.NewPolynomial)

	// builderJson
	pallet.RegisterNodeBuilder(nodejson.Category, jsonFilter, nodejson.NewFilter)
	pallet.RegisterNodeBuilder(nodejson.Category, dataStore, nodejson.NewStore)

	// builderFlowNetworks
	pallet.RegisterNodeBuilder(flow.Category, flowNetwork, flow.NewNetwork)
	pallet.RegisterNodeBuilder(flow.Category, flowPoint, flow.NewFFPoint)
	pallet.RegisterNodeBuilder(flow.Category, flowPointWrite, flow.NewFFPointWrite)
	pallet.RegisterNodeBuilder(flow.Category, flowSchedule, flow.NewFFSchedule)

	// builderHVAC
	pallet.RegisterNodeBuilder(hvac.Category, accumulationPeriod, hvac.NewAccumulationPeriod)
	pallet.RegisterNodeBuilder(hvac.Category, deadBandNode, hvac.NewDeadBand)
	pallet.RegisterNodeBuilder(hvac.Category, leadLagSwitch, hvac.NewLeadLagSwitch)
	pallet.RegisterNodeBuilder(hvac.Category, pidNode, hvac.NewPIDNode)
	pallet.RegisterNodeBuilder(hvac.Category, pacControlNode, hvac.NewPACControl)
	pallet.RegisterNodeBuilder(hvac.Category, sensorSelect, hvac.NewSensorSelect)
	pallet.RegisterNodeBuilder(hvac.Category, psychroDBRH, hvac.NewPsychroDBRH)
	pallet.RegisterNodeBuilder(hvac.Category, psychroDBDP, hvac.NewPsychroDBDP)
	pallet.RegisterNodeBuilder(hvac.Category, psychroDBWB, hvac.NewPsychroDBWB)

	// builderLatch
	pallet.RegisterNodeBuilder(latch.Category, numLatch, latch.NewNumLatch)
	pallet.RegisterNodeBuilder(latch.Category, stringLatch, latch.NewStringLatch)
	pallet.RegisterNodeBuilder(latch.Category, setResetLatch, latch.NewSetResetLatch)

	// builderStreams
	pallet.RegisterNodeBuilder(streams.Category, flatLine, streams.NewFlatline)
	pallet.RegisterNodeBuilder(streams.Category, rollingAverage, streams.NewRollingAverage)

	// builderConst
	pallet.RegisterNodeBuilder(constant.Category, constNum, constant.NewNumber)
	pallet.RegisterNodeBuilder(constant.Category, constStr, constant.NewString)
	pallet.RegisterNodeBuilder(constant.Category, constBool, constant.NewBoolean)

	// builderMath
	pallet.RegisterNodeBuilder(math.Category, abs, math.NewAbsolute)
	pallet.RegisterNodeBuilder(math.Category, add, math.NewAdd)
	pallet.RegisterNodeBuilder(math.Category, sub, math.NewSub)
	pallet.RegisterNodeBuilder(math.Category, multiply, math.NewMultiply)
	pallet.RegisterNodeBuilder(math.Category, divide, math.NewDivide)
	pallet.RegisterNodeBuilder(math.Category, modulo, math.NewModulo)
	pallet.RegisterNodeBuilder(math.Category, ceiling, math.NewCeiling)
	pallet.RegisterNodeBuilder(math.Category, floor, math.NewFloor)
	pallet.RegisterNodeBuilder(mathematics.Category, mathAdvanced, mathematics.NewAdvanced)

	// builderConversion
	pallet.RegisterNodeBuilder(conversion.Category, conversionString, conversion.NewString)
	pallet.RegisterNodeBuilder(conversion.Category, conversionNum, conversion.NewNumber)
	pallet.RegisterNodeBuilder(conversion.Category, conversionBool, conversion.NewBoolean)

	// builderSwitch
	pallet.RegisterNodeBuilder(switches.Category, switchNode, switches.NewSwitch)
	pallet.RegisterNodeBuilder(switches.Category, selectNum, switches.NewSelectNum)
	pallet.RegisterNodeBuilder(switches.Category, numOutputSelect, switches.NewNumOutputSelect)

	// builderBoolean
	pallet.RegisterNodeBuilder(boolean.Category, and, boolean.NewAnd)
	pallet.RegisterNodeBuilder(boolean.Category, or, boolean.NewOr)
	pallet.RegisterNodeBuilder(boolean.Category, xor, boolean.NewXor)
	pallet.RegisterNodeBuilder(boolean.Category, not, boolean.NewNot)
	pallet.RegisterNodeBuilder(boolean.Category, toggle, boolean.NewToggle)

	// builderCompare
	pallet.RegisterNodeBuilder(compare.Category, greaterThan, compare.NewGreaterThan)
	pallet.RegisterNodeBuilder(compare.Category, lessThan, compare.NewLessThan)
	pallet.RegisterNodeBuilder(compare.Category, equal, compare.NewEqual)
	pallet.RegisterNodeBuilder(compare.Category, equalString, compare.NewEqualString)
	pallet.RegisterNodeBuilder(compare.Category, between, compare.NewBetween)
	pallet.RegisterNodeBuilder(compare.Category, hysteresis, compare.NewHysteresis)

	// builderStatistics
	pallet.RegisterNodeBuilder(statistics.Category, min, statistics.NewMin)
	pallet.RegisterNodeBuilder(statistics.Category, max, statistics.NewMax)
	pallet.RegisterNodeBuilder(statistics.Category, avg, statistics.NewAvg)
	pallet.RegisterNodeBuilder(statistics.Category, minMaxAvg, statistics.NewMinMaxAvg)
	pallet.RegisterNodeBuilder(statistics.Category, rangeNode, statistics.NewRange)
	pallet.RegisterNodeBuilder(statistics.Category, rank, statistics.NewRank)
	pallet.RegisterNodeBuilder(statistics.Category, median, statistics.NewMedian)

	// builderPoints
	pallet.RegisterNodeBuilder(point.Category, boolWriteable, point.NewBooleanWriteable)
	pallet.RegisterNodeBuilder(point.Category, numWriteable, point.NewNumericWriteable)
	pallet.RegisterNodeBuilder(point.Category, stringWriteable, point.NewStringWriteable)

	// builderCount
	pallet.RegisterNodeBuilder(count.Category, countNode, count.NewCount)
	pallet.RegisterNodeBuilder(count.Category, countNumNode, count.NewCountNum)
	pallet.RegisterNodeBuilder(count.Category, countStringNode, count.NewCountString)

	// builderTrigger
	pallet.RegisterNodeBuilder(trigger.Category, covNode, trigger.NewCOVNode)
	pallet.RegisterNodeBuilder(trigger.Category, random, trigger.NewRandom)
	pallet.RegisterNodeBuilder(trigger.Category, iterate, trigger.NewIterate)

	// builderTiming
	pallet.RegisterNodeBuilder(timing.Category, clock, timing.NewClock)
	pallet.RegisterNodeBuilder(timing.Category, date, timing.NewDate)
	pallet.RegisterNodeBuilder(timing.Category, delay, timing.NewDelay)
	pallet.RegisterNodeBuilder(timing.Category, delayOn, timing.NewDelayOn)
	pallet.RegisterNodeBuilder(timing.Category, delayOff, timing.NewDelayOff)
	pallet.RegisterNodeBuilder(timing.Category, dutyCycle, timing.NewDutyCycle)
	pallet.RegisterNodeBuilder(timing.Category, minOnOff, timing.NewMinOnOff)
	pallet.RegisterNodeBuilder(timing.Category, oneShot, timing.NewOneShot)
	pallet.RegisterNodeBuilder(timing.Category, stopwatch, timing.NewStopwatch)

	// builderProtocols
	pallet.RegisterNodeBuilder(bacnetio.Category, bacnetServer, bacnetio.NewServer)
	pallet.RegisterNodeBuilder(bacnetio.Category, bacnetAI, bacnetio.NewAI)
	pallet.RegisterNodeBuilder(bacnetio.Category, bacnetAO, bacnetio.NewAO)
	pallet.RegisterNodeBuilder(bacnetio.Category, bacnetAV, bacnetio.NewAV)
	pallet.RegisterNodeBuilder(bacnetio.Category, bacnetBV, bacnetio.NewBV)

	// builderMQTT
	pallet.RegisterNodeBuilder(broker.Category, mqttBroker, broker.NewBroker)
	pallet.RegisterNodeBuilder(broker.Category, mqttSub, broker.NewMqttSub)
	pallet.RegisterNodeBuilder(broker.Category, mqttPub, broker.NewMqttPub)
}
