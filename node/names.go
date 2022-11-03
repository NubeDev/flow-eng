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
	SetPoint   InputName = "set-point"
	Offset     InputName = "offset"
	CoolOffset InputName = "offset-cool"
	HeatOffset InputName = "offset-heat"
	DeadBand   InputName = "dead-band"

	AlertDelayMins InputName = "alert-delay-mins"

	Comment   InputName = "comment"
	InNumber  InputName = "number"
	InString  InputName = "string"
	InBoolean InputName = "boolean"

	Enable InputName = "enable"

	OnInterval  InputName = "on-interval"
	OffInterval InputName = "off-interval"
	Interval    InputName = "interval"

	Duration InputName = "duration"

	TriggerOnCount InputName = "trigger on count"
	Count          InputName = "count"
	CountUp        InputName = "count-up"
	CountDown      InputName = "count-down"

	Min InputName = "min"
	Max InputName = "max"

	InMin  InputName = "in-min"
	InMax  InputName = "in-max"
	OutMin InputName = "out-min"
	OutMax InputName = "out-max"

	Delete InputName = "delete"

	Size    InputName = "size"
	MaxSize InputName = "max-size"
	MinSize InputName = "min-size"

	In   InputName = "in"
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
)

const (
	Ok OutputName = "ok"
	//Result OutputName = "result"

	ErrMsg OutputName = "error"
	Msg    OutputName = "message"

	FlatLine OutputName = "flatline"

	Trigger OutputName = "trigger"
	Toggle  OutputName = "toggle"

	String  OutputName = "string"
	Boolean OutputName = "boolean"
	Float   OutputName = "float"

	OutNot OutputName = "out not"

	OutTopic OutputName = "topic"
	Out      OutputName = "out"

	Connected OutputName = "connected"

	CountOut OutputName = "count"

	MinOn  OutputName = "min-on"
	MinOff OutputName = "min-off"

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
)
