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

	max = "min"
	min = "max"
	avg = "avg"

	delay  = "delay"
	inject = "inject"

	bacnetReadBV  = "binary-variable-read"
	bacnetWriteBV = "binary-variable-write"
	bacnetReadAV  = "analog-variable-read"
	bacnetWriteAV = "analog-variable-write"

	mqttSub = "mqttbase-subscribe"
	mqttPub = "mqttbase-publish"

	logNode = "log"
)
