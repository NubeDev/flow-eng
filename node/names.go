package node

type DataTypes string
type InputName string
type OutputName string

const (
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

	Comment  InputName = "comment"
	InNumber InputName = "number"
	InString InputName = "string"

	Enable InputName = "enable"

	Interval InputName = "interval"

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

	Input_ InputName = "input"
	InputA InputName = "a"
	InputB InputName = "b"
	InputC InputName = "c"
	InputD InputName = "d"

	Connection   InputName = "topic"
	Topic        InputName = "topic"
	TriggerInput InputName = "trigger"
	Message      InputName = "message"
	Subject      InputName = "subject"

	Filter   InputName = "filter"
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
	Result OutputName = "result"

	ErrMsg OutputName = "error"
	Msg    OutputName = "message"

	Flatline OutputName = "flatline"

	Trigger OutputName = "trigger"
	Toggle  OutputName = "toggle"
	Out     OutputName = "out"

	OutNot OutputName = "out not"

	Out1 OutputName = "out1"
	Out2 OutputName = "out2"
	Out3 OutputName = "out3"
	Out4 OutputName = "out4"

	Above OutputName = "above"
	Below OutputName = "below"

	GraterThan OutputName = "grater"
	LessThan   OutputName = "less"
	Equal      OutputName = "equal"
)
