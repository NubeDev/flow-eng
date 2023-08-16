package rules

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"math"
	"math/rand"
	"time"
)

func (inst *RQL) JSONFilter(body string, filter string) (interface{}, error) {
	// var m interface{}
	// err := mapstructure.Decode(body, &m)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// gjson.Parse(body)

	// if err := json.Unmarshal([]byte(body), &m); err != nil {
	// 	// return nil, err
	// 	fmt.Println(err)
	// }
	//
	// fmt.Println(m)

	value := gjson.Get(body, filter)

	return value.Value(), nil
}

func (inst *RQL) Parse(body string) (interface{}, error) {
	// var m interface{}
	// err := mapstructure.Decode(body, &m)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// gjson.Parse(body)

	// if err := json.Unmarshal([]byte(body), &m); err != nil {
	// 	// return nil, err
	// 	fmt.Println(err)
	// }
	//
	// fmt.Println(m)

	m, ok := gjson.Parse(body).Value().(interface{})
	fmt.Println(m)
	fmt.Println(ok)
	if !ok {

		// not a map
	}

	return m, nil
}

func (inst *RQL) Stringify(body interface{}) (string, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return "", nil
	}
	return string(b), nil
}

// RandInt returns a random int within the specified range.
func (inst *RQL) RandInt(range1, range2 int) int {
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

// RandFloat returns a random float64 within the specified range.
func (inst *RQL) RandFloat(range1, range2 float64) float64 {
	if range1 == range2 {
		return range1
	}
	var min, max float64
	if range1 > range2 {
		max = range1
		min = range2
	} else {
		max = range2
		min = range1
	}
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}

// LimitToRange returns the input value clamped within the specified range
func (inst *RQL) LimitToRange(value float64, range1 float64, range2 float64) float64 {
	if range1 == range2 {
		return range1
	}
	var min, max float64
	if range1 > range2 {
		max = range1
		min = range2
	} else {
		max = range2
		min = range1
	}
	return math.Min(math.Max(value, min), max)
}

// RoundTo returns the input value rounded to the specified number of decimal places.
func (inst *RQL) RoundTo(value float64, decimals uint32) float64 {
	if decimals < 0 {
		return value
	}
	return math.Round(value*math.Pow(10, float64(decimals))) / math.Pow(10, float64(decimals))
}

// Scale returns the (float64) input value (between inputMin and inputMax) scaled to a value between outputMin and outputMax
func (inst *RQL) Scale(value, inMin, inMax, outMin, outMax float64) float64 {
	if inMin == inMax || outMin == outMax {
		return value
	}
	scaled := ((value-inMin)/(inMax-inMin))*(outMax-outMin) + outMin
	if scaled > math.Max(outMin, outMax) {
		return math.Max(outMin, outMax)
	} else if scaled < math.Min(outMin, outMax) {
		return math.Min(outMin, outMax)
	} else {
		return scaled
	}
}
