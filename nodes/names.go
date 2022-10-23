package nodes

const (
	constNum  = "const-num"
	constBool = "const-bool"
	constStr  = "const-string"

	add      = "add"
	divide   = "divide"
	sub      = "subtract"
	multiply = "multiply"

	mathAdvanced = "advanced"

	and    = "and"
	or     = "or"
	not    = "not"
	xor    = "xor"
	toggle = "toggle"

	gmailNode = "gmail"
	pingNode  = "ping"

	numLatch      = "numeric-latch"
	stringLatch   = "string-latch"
	setResetLatch = "set-reset-latch"

	pointNumber  = "point-num"
	pointBoolean = "point-bool"

	rampNode = "ramp"

	scaleNode = "scale"
	limitNode = "limit"

	dataStore  = "data-store"
	jsonFilter = "json-filter"

	flowNetwork = "flow-network"
	flowPoint   = "flow-point"

	flowLoopCount = "flow-loop-count"

	getHttpNode   = "http-get"
	writeHttpNode = "http-write"

	deadBandNode = "dead-band"

	logicCompare = "compare"
	between      = "between"
	hysteresis   = "hysteresis"
	funcNode     = "func"

	conversionString = "conversion-string"
	conversionNum    = "conversion-number"
	conversionBool   = "conversion-bool"

	max = "min"
	min = "max"
	avg = "avg"

	flatLine = "flatLine"

	// count
	countNode = "count"

	// trigger
	randomFloat = "random-number"
	inject      = "inject"

	delay         = "delay"
	delayOn       = "delay-on"
	delayOff      = "delay-off"
	delayMinOnOff = "min on off"

	// switch
	switchNode = "switch"
	selectNum  = "select-numeric"

	linkInput  = "link-input"
	linkOutput = "link-output"

	bacnetServer = "bacnet-server"
	bacnetAI     = "analog-input"
	bacnetAO     = "analog-output"
	bacnetAV     = "analog-variable"
	bacnetBI     = "binary-input"
	bacnetBO     = "binary-output"
	bacnetBV     = "binary-variable"

	mqttBroker = "mqtt-broker"
	mqttSub    = "mqtt-subscribe"
	mqttPub    = "mqtt-publish"

	logNode = "log"
)
