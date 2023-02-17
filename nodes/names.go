package nodes

const (
	constNum  = "const-num"
	constBool = "const-boolean"
	constStr  = "const-string"

	abs          = "absolute-value"
	add          = "add"
	divide       = "divide"
	sub          = "subtract"
	multiply     = "multiply"
	modulo       = "modulo"
	floor        = "floor"
	ceiling      = "ceiling"
	mathAdvanced = "advanced"

	and    = "and"
	or     = "or"
	not    = "not"
	xor    = "xor"
	toggle = "toggle"

	gmail    = "gmail"
	pingNode = "ping"

	numLatch      = "numeric-latch"
	stringLatch   = "string-latch"
	setResetLatch = "set-reset-latch"

	boolWriteable   = "boolean-writeable"
	numWriteable    = "numeric-writeable"
	stringWriteable = "string-writeable"

	rampNode = "ramp"

	scaleNode = "scale"
	limitNode = "limit"
	fade      = "fade"
	rateLimit = "rate-limit"
	round     = "round"
	waveform  = "waveform"

	dataStore  = "data-store"
	jsonFilter = "json-filter"

	flowNetwork    = "flow-network"
	flowPoint      = "flow-point"
	flowSchedule   = "flow-schedule"
	flowPointWrite = "flow-point-write"

	flowLoopCount = "flow-loop-count"

	getHttpNode   = "http-get"
	writeHttpNode = "http-write"

	accumulationPeriod = "accumulation-period"
	deadBandNode       = "dead-band"
	leadLagSwitch      = "lead-lag-switch"
	pidNode            = "pid"
	pacControlNode     = "pac-control"
	sensorSelect       = "sensor-select"
	psychroDBRH        = "psychrometrics-db-rh"
	psychroDBDP        = "psychrometrics-db-dp"
	psychroDBWB        = "psychrometrics-db-wb"

	greaterThan = "greater-than"
	lessThan    = "less-than"
	equal       = "equal"
	equalString = "equal-string"
	between     = "between"
	hysteresis  = "hysteresis"
	funcNode    = "func"

	conversionString = "conversion-string"
	conversionNum    = "conversion-number"
	conversionBool   = "conversion-boolean"

	max       = "max"
	min       = "min"
	avg       = "average"
	minMaxAvg = "min-max-avg"
	rangeNode = "range"
	rank      = "rank"
	median    = "median"

	flatLine = "flatline"

	// filters
	preventNull        = "prevent-null"
	preventEqualFloat  = "prevent-equal-float"
	preventEqualString = "prevent-equal-string"
	onlyBetween        = "only-between"
	onlyGreater        = "only-greater"
	onlyLower          = "only-lower"
	preventDuplicates  = "prevent-duplicates"

	// count
	countNode       = "count"
	countNumNode    = "count-number"
	countStringNode = "count-string"

	// trigger
	covNode    = "change-of-value"
	random     = "random"
	injectNode = "inject"
	iterate    = "iterate"

	delay     = "delay"
	delayOn   = "delay-on"
	delayOff  = "delay-off"
	dutyCycle = "duty-cycle"
	minOnOff  = "min-on-off"
	oneShot   = "one-shot"
	stopwatch = "stopwatch"

	// switch
	switchNode      = "switch"
	selectNum       = "select-numeric"
	numOutputSelect = "output-select-numeric"

	linkInput      = "link-input-string"
	linkOutput     = "link-output-string"
	linkInputNum   = "link-input-number"
	linkOutputNum  = "link-output-number"
	linkInputBool  = "link-input-boolean"
	linkOutputBool = "link-output-boolean"

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

	subFlowFolder = "folder"
	inputFloat    = "input-float"
	inputBool     = "input-boolean"
	inputString   = "input-string"

	outputFloat  = "output-float"
	outputBool   = "output-boolean"
	outputString = "output-string"
)
