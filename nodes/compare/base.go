package compare

const (
	category = "compare"
)

const (
	GreaterThan = "greater-than"
	LessThan    = "less-than"
	Equal       = "equal"
	betweenNode = "between"
	hysteresis  = "hysteresis"
)

func B2F(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func zeroToOne(b float64) float64 {
	if b > 0 {
		return 0
	}
	return 1
}
