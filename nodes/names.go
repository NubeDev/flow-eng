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
	jsonFilter    = "json-filter"

	flowNetwork = "network"
	flowDevice  = "device"
	flowPoint   = "point"

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

	delay         = "delay"
	inject        = "inject"
	delayOn       = "delayOn"
	delayMinOnOff = "min on off"

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
