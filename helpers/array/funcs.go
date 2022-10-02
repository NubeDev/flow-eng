package array

func Between(x, min, max float64) (between, below, above bool) {
	if (x >= min) && (x <= max) {
		between = true
	}
	if x < min {
		below = true
	}
	if x > max {
		above = true
	}
	return between, below, above
}

func Compare(array []float64) (greater, less, equal bool) {
	if len(array) < 1 {
		return false, false, false
	}
	val1 := array[0]
	val2 := array[1]
	if val1 > val2 {
		greater = true
	}
	if val1 < val2 {
		less = true
	}
	if val1 == val2 {
		equal = true
	}
	return greater, less, equal
}

func OneIsTrue(array []bool) bool {
	if len(array) == 0 {
		return false
	}
	for _, item := range array {
		if item {
			return true
		}
	}
	return false
}

func AllTrue(array []bool) bool {
	if len(array) == 0 {
		return false
	}
	for _, item := range array {
		if !item {
			return false
		}
	}
	return true
}

func AllTrueFloat64(array []float64) bool {
	if len(array) == 0 {
		return false
	}
	for _, item := range array {
		if item == 1 {
			return false
		}
	}
	return true
}

func Divide(array []float64) float64 {
	var out float64
	var length = len(array)
	if length == 0 {
		return 0
	}
	var firstValue = array[0]
	for i, value := range array {
		if i != 0 {
			firstValue = firstValue / value
			if i == length-1 {
				out = firstValue
			}
		}

	}
	return out
}

func Multiply(array []float64) float64 {
	var out float64
	var length = len(array)
	if length == 0 {
		return 0
	}
	var firstValue = array[0]
	for i, value := range array {
		if i != 0 {
			firstValue = firstValue * value
			if i == length-1 {
				out = firstValue
			}
		}

	}
	return out
}

func Add(array []float64) float64 {
	var out float64
	for _, num := range array {
		out = out + num
	}
	return out
}

func Subtract(array []float64) float64 {
	var out float64
	var length = len(array)
	if length == 0 {
		return 0
	}
	var firstValue = array[0]
	for i, value := range array {
		if i != 0 {
			firstValue = firstValue - value
			if i == length-1 {
				out = firstValue
			}
		}

	}
	return out
}

func OneIsTrueFloat64(array []float64) bool {
	if len(array) == 0 {
		return false
	}
	for _, item := range array {
		if item == 1 {
			return true
		}
	}
	return false
}

func MinMaxFloat64(array []float64) (min float64, max float64) {
	if len(array) == 0 {
		return 0, 0
	}
	max = array[0]
	min = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}

func MaxFloat64(array []float64) float64 {
	if len(array) == 0 {
		return 0
	}
	max := array[0]
	for _, item := range array {
		if item > max {
			max = item
		}
	}
	return max
}

func MinFloat64(array []float64) float64 {
	min := array[0]
	for _, item := range array {
		if item < min {
			min = item
		}
	}
	return min
}
