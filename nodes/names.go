package nodes

const (
	constNum = "const-num"
	constStr = "const-string"

	add      = "add"
	divide   = "divide"
	sub      = "subtract"
	multiply = "multiply"

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
	jsonFilter    = "json-filter"

	flowNetwork = "network"
	flowDevice  = "device"
	flowPoint   = "point"

	flowLoopCount = "flow-loop-count"

	getNode = "get"

	deadBandNode = "dead-band"

	logicCompare = "compare"
	between      = "between"
	hysteresis   = "hysteresis"
	funcNode     = "func"

	stringToNum = "string-to-number"
	numToString = "number-to-string"

	max = "min"
	min = "max"
	avg = "avg"

	flatline = "flatline"

	delay   = "delay"
	inject  = "inject"
	delayOn = "delayOn"

	// switch
	switchNode = "switch"
	selectNum  = "select-numeric"

	linkInput  = "link-input"
	linkOutput = "link-output"

	bacnetServer = "server"
	bacnetAI     = "analog-input"
	bacnetAO     = "analog-output"
	bacnetAV     = "analog-variable"
	bacnetBI     = "binary-input"
	bacnetBO     = "binary-output"
	bacnetBV     = "binary-variable"

	mqttSub = "mqtt-subscribe"
	mqttPub = "mqtt-publish"

	logNode = "log"
)
