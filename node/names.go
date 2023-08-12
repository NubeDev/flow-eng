package node

type DataTypes string
type InputName string
type InputHelp string
type OutputHelp string
type OutputName string

const (
	TypeBool   DataTypes = "boolean"
	TypeString DataTypes = "string"
	TypeInt    DataTypes = "int"
	TypeFloat  DataTypes = "number"
	TypeNumber DataTypes = "number"
)

const (
	InputNamePrefix  string = "in"
	OutputNamePrefix string = "out"
)

const (
	Setpoint   InputName = "setpoint"
	Offset     InputName = "offset"
	CoolOffset InputName = "offset-cool"
	HeatOffset InputName = "offset-heat"
	DeadBand   InputName = "dead-band"

	ZoneTemp        InputName = "zone-temp"
	ClgOffset       InputName = "clg-offset"
	HtgOffset       InputName = "htg-offset"
	StgUpDelay      InputName = "stage-up-delay"
	ModeChangeDelay InputName = "mode-change-delay"
	EconoAllow      InputName = "econo-allow"
	OATemp          InputName = "oa-temp"
	EconoHigh       InputName = "econo-high"
	EconoLow        InputName = "econo-low"
	FanStatus       InputName = "fan-status"
	ClgLockout      InputName = "clg-lockout"
	HtgLockout      InputName = "htg-lockout"

	AlertDelayMins InputName = "alert-delay-mins"

	Comment   InputName = "comment"
	InNumber  InputName = "number"
	InString  InputName = "string"
	InBoolean InputName = "boolean"

	Enable InputName = "enable"

	MinOnTime  InputName = "min-on-time"
	MinOffTime InputName = "min-off-time"
	Interval   InputName = "interval"
	DutyCycle  InputName = "duty-cycle"
	Iterations InputName = "iterations"
	Period     InputName = "period"
	Amplitude  InputName = "amplitude"

	Duration InputName = "duration"

	TriggerOnCount InputName = "trigger on count"
	Count          InputName = "count"
	CountUp        InputName = "count-up"
	CountDown      InputName = "count-down"
	StepSize       InputName = "step-size"
	SetValue       InputName = "set-value"

	MinInput InputName = "min"
	MaxInput InputName = "max"

	InMin  InputName = "in-min"
	InMax  InputName = "in-max"
	OutMin InputName = "out-min"
	OutMax InputName = "out-max"

	Delete InputName = "delete"

	Size     InputName = "size"
	MaxSize  InputName = "max-size"
	MinSize  InputName = "min-size"
	Decimals InputName = "decimals"

	In1  InputName = "in1"
	In2  InputName = "in2"
	In3  InputName = "in3"
	In4  InputName = "in4"
	In10 InputName = "in10"
	In11 InputName = "in11"
	In12 InputName = "in12"
	In13 InputName = "in13"
	In14 InputName = "in14"
	In15 InputName = "in15"
	In16 InputName = "in16"

	InputA InputName = "a"
	InputB InputName = "b"
	InputC InputName = "c"
	InputD InputName = "d"

	X   InputName = "x"
	X0  InputName = "x0"
	X1  InputName = "x1"
	X2  InputName = "x2"
	X3  InputName = "x3"
	X4  InputName = "x4"
	X5  InputName = "x5"
	X6  InputName = "x6"
	X7  InputName = "x7"
	X8  InputName = "x8"
	X9  InputName = "x9"
	X10 InputName = "x10"

	In        InputName = "input"
	Match     InputName = "match"
	Value     InputName = "value"
	Threshold InputName = "threshold"

	Start InputName = "start"
	Stop  InputName = "stop"

	Connection   InputName = "topic"
	URL          InputName = "url"
	Topic        InputName = "topic"
	TriggerInput InputName = "trigger"
	Message      InputName = "message"
	Subject      InputName = "subject"

	Filter   InputName = "filter"
	Body     InputName = "body"
	Equation InputName = "equation"

	Ip          InputName = "ip"
	Time        InputName = "time"
	NetworkPort InputName = "port"

	Delay        InputName = "delay"
	DelaySeconds InputName = "delay (s)"
	Selection    InputName = "select"

	From InputName = "from"
	To   InputName = "to"

	UUID          InputName = "uuid"
	Name          InputName = "name"
	ObjectId      InputName = "object-id"
	ObjectType    InputName = "object-type"
	OverrideInput InputName = "override-value"

	RisingEdge  InputName = "rising-edge"
	FallingEdge InputName = "falling-edge"

	Switch  InputName = "switch"
	InTrue  InputName = "inTrue"
	InFalse InputName = "inFalse"

	Latch InputName = "latch"
	Set   InputName = "set"
	Reset InputName = "reset"

	ProcessValue InputName = "process-value"
	MinOut       InputName = "min-out"
	MaxOut       InputName = "max-out"
	InP          InputName = "in-p"
	InI          InputName = "in-i"
	InD          InputName = "in-d"
	PIDDirection InputName = "direction"
	Bias         InputName = "bias"
	Manual       InputName = "manual"
	IntervalSecs InputName = "interval-secs"

	AccumulationIn InputName = "accumulation-in"

	RotateLead InputName = "rotate-lead"
	Stage1     InputName = "stage-1"
	Stage2     InputName = "stage-2"

	Fallback InputName = "fallback"

	DryBulbTemp  InputName = "dry-bulb-temp"
	WetBulbTemp  InputName = "wet-bulb-temp"
	DewPointTemp InputName = "dew-point-temp"
	RelHumid     InputName = "rel-humidity-%"
)

const (
	Ok OutputName = "ok"
	// Result OutputName = "result"

	ErrMsg       OutputName = "error"
	Msg          OutputName = "message"
	PollingStats OutputName = "polling-status"
	PollingCount OutputName = "polling-count"

	FlatLine OutputName = "flatline"

	Complete OutputName = "complete"

	Trigger OutputName = "trigger"
	Toggle  OutputName = "toggle"

	String  OutputName = "string"
	Boolean OutputName = "boolean"
	Float   OutputName = "float"

	OutNot OutputName = "out not"

	OutPayload      OutputName = "payload"
	PeriodStart     OutputName = "period-start"
	PeriodStop      OutputName = "period-stop"
	NextStart       OutputName = "next-start"
	NextStop        OutputName = "next-stop"
	PeriodStartUnix OutputName = "period-start-unix"
	PeriodStopUnix  OutputName = "period-stop-unix"
	NextStartUnix   OutputName = "next-start-unix"
	NextStopUnix    OutputName = "next-stop-unix"

	TimeString OutputName = "time-string"
	Hour       OutputName = "hour"
	Min        OutputName = "min"
	Sec        OutputName = "sec"
	Ms         OutputName = "ms"
	LongString OutputName = "long-string"
	TzOffset   OutputName = "tz-offset"
	UnixSecs   OutputName = "unix-secs"
	DateString OutputName = "date-string"
	DayString  OutputName = "day-string"
	DayOfWeek  OutputName = "day-of-week"
	Date       OutputName = "date"
	Month      OutputName = "month"
	Year       OutputName = "year"

	OutTopic OutputName = "topic"
	Out      OutputName = "output"
	Response OutputName = "response"
	OnStart  OutputName = "on-start"
	Elapsed  OutputName = "elapsed"

	MinOutput OutputName = "min"
	MaxOutput OutputName = "max"
	AvgOutput OutputName = "avg"

	LastUpdated     OutputName = "last-updated"
	CurrentPriority OutputName = "current-priority"

	Connected OutputName = "connected"

	CountOut OutputName = "count"
	OutEqTo  OutputName = "out=to"

	MinOnActive  OutputName = "min-on-active"
	MinOffActive OutputName = "min-off-active"

	Out1 OutputName = "out1"
	Out2 OutputName = "out2"
	Out3 OutputName = "out3"
	Out4 OutputName = "out4"

	Above OutputName = "above"
	Below OutputName = "below"

	GreaterThan      OutputName = "greater"
	GreaterThanEqual OutputName = "greater or equal"
	LessThanEqual    OutputName = "less or equal"
	LessThan         OutputName = "less"
	Equal            OutputName = "equal"
	NotEqual         OutputName = "not-equal"

	ClgMode        OutputName = "clg-mode"
	HtgMode        OutputName = "htg-mode"
	CompStage      OutputName = "comp-stage"
	EconoMode      OutputName = "econo-mode"
	OADamper       OutputName = "oa-damper"
	ReversingValve OutputName = "reversing-valve"
	Compressor1    OutputName = "compressor-1"
	Compressor2    OutputName = "compressor-2"

	LeadUnit     OutputName = "lead-unit"
	LeadUnitBool OutputName = "lead-unit-boolean"
	EnableA      OutputName = "enable-A"
	EnableB      OutputName = "enable-B"

	ValidOutput OutputName = "valid-output"

	PeriodConsumption OutputName = "period-consumption"
	LastPeriodEndVal  OutputName = "last-period-end-value"
	PeriodDuration    OutputName = "period-duration"
	NextTrigger       OutputName = "next-trigger"

	WetBulbTempO     OutputName = "wet-bulb-temp"
	DewPointTempO    OutputName = "dew-point-temp"
	RelHumPercO      OutputName = "rel-humidity-%"
	HumRatioO        OutputName = "humidity-ratio"
	VaporPres        OutputName = "vapor-pressure"
	MoistAirEnthalpy OutputName = "moist-air-enthalpy"
	MoistAirVolume   OutputName = "moist-air-volume"
	DegreeSaturation OutputName = "degree-of-saturation"
)
