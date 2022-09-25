package filter

import (
	"errors"
	"math"
)

// NODEs will single in/out
// only-true
// only-false
// prevent-null
// prevent-duplicates

// NODEs will double in
// only-between
// only-lower
// prevent-equal

const (
	onlyTrue    = "onlyTrue"
	onlyFalse   = "onlyFalse"
	preventNull = "preventNull"
)

func mathFunc(def string, x float64) (float64, error) {
	switch def {
	case onlyTrue:
		return onlyTrueFunc(x), nil
	case onlyFalse:
		return math.Tan(x), nil
	case preventNull:
		return math.Tan(x), nil
	}
	return 0, errors.New("math function not found")
}

func onlyTrueFunc(v float64) float64 {
	if v >= 1 {
		return 1
	}
	return 0
}

func onlyFalseFunc(v float64) float64 {
	if v <= 1 {
		return 1
	}
	return 0
}

func preventNullFunc(v interface{}) bool {
	if v != nil {
		return true
	}
	return false
}

// preventDuplicates if new value == last value return true
func preventDuplicatesFunc(newValue, lastValue interface{}) bool {
	if newValue == lastValue {
		return true
	}
	return false
}
