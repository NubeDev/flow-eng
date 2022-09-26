package nodes

const (
	constNum = "const-num"
	constStr = "const-string"

	add      = "add"
	divide   = "divide"
	sub      = "subtract"
	multiply = "multiply"

	and = "and"
	or  = "or"
	not = "not"
	xor = "xor"

	flowLoopCount = "flow-loop-count"

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

	delay   = "delay"
	inject  = "inject"
	delayOn = "delayOn"

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

	mqttSub = "pointbus-subscribe"
	mqttPub = "pointbus-publish"

	logNode = "log"
)
