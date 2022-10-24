package integer

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func GetUnit64(in interface{}) (val uint64) {
	switch i := in.(type) {
	case bool:
		if i {
			return 1
		}
		return 0
	case int:
		val = uint64(i)
	case float64:
		val = uint64(i)
	case float32:
		val = uint64(i)
	case int64:
		val = uint64(i)
	case uint64:
		val = i
	case string:
		if s, err := strconv.ParseUint(fmt.Sprintf("%v", in), 10, 64); err == nil {
			val = s
		} else {
			val = 0
		}
	default:
		return 0
	}
	return val
}

func New(value int) *int {
	return &value
}

func NewUint(value uint) *uint {
	return &value
}

func NonNil(b *int) int {
	if b == nil {
		return 0
	} else {
		return *b
	}
}

func IsNil(b *int) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

func IsUnitNil(b *uint) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

func IsUnit32Nil(b *uint32) bool {
	if b == nil {
		return true
	} else {
		return false
	}
}

// RandInt returns a random int within the specified range.
func RandInt(range1, range2 int) int {
	if range1 == range2 {
		return range1
	}
	var min, max int
	if range1 > range2 {
		max = range1
		min = range2
	} else {
		max = range2
		min = range1
	}
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}
