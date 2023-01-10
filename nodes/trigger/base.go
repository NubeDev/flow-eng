package trigger

type TimeUnits string

const (
	Category    = "trigger"
	COV         = "change-of-value"
	Iterator    = "iterator"
	Inject      = "inject"
	RandomFloat = "random-number"
)

const (
	Milliseconds TimeUnits = "milliseconds"
	Seconds      TimeUnits = "seconds"
	Minutes      TimeUnits = "minutes"
	Hours        TimeUnits = "hours"
)
