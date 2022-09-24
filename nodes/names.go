package nodes

const (
	constNum = "const-num"
	add      = "add"
	divide   = "divide"
	sub      = "subtract"
	multiply = "multiply"

	and     = "and"
	or      = "or"
	not     = "not"
	greater = "greater"
	less    = "less"

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

	selectNum        = "select-numeric"
	connectionInput  = "input"
	connectionOutput = "output"

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
