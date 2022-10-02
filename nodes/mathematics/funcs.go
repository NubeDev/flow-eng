package mathematics

import (
	"errors"
	"math"
)

func mathFunc(def string, x float64) (float64, error) {
	switch def {
	case acos:
		return math.Acos(x), nil
	case asin:
		return math.Asin(x), nil
	case atan:
		return math.Atan(x), nil
	case cbrt:
		return math.Cbrt(x), nil
	case cos:
		return math.Cos(x), nil
	case exp:
		return math.Exp(x), nil
	case log:
		return math.Log(x), nil
	case log10:
		return math.Log10(x), nil
	case sin:
		return math.Sin(x), nil
	case sqrt:
		return math.Sqrt(x), nil
	case tan:
		return math.Tan(x), nil
	}
	return 0, errors.New("math function not found")

}
