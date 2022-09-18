package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"strconv"
	"strings"
)

type PriArray struct {
	P1  *float64 `json:"_1"`
	P2  *float64 `json:"_2"`
	P3  *float64 `json:"_3"`
	P4  *float64 `json:"_4"`
	P5  *float64 `json:"_5"`
	P6  *float64 `json:"_6"`
	P7  *float64 `json:"_7"`
	P8  *float64 `json:"_8"`
	P9  *float64 `json:"_9"`
	P10 *float64 `json:"_10"`
	P11 *float64 `json:"_11"`
	P12 *float64 `json:"_12"`
	P13 *float64 `json:"_13"`
	P14 *float64 `json:"_14"`
	P15 *float64 `json:"_15"`
	P16 *float64 `json:"_16"`
}

func set(part string) *float64 {
	if part == "Null" {
		return nil
	} else {
		f, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil
		}
		return float.New(f)
	}
}

type PriAndValue struct {
	Number int
	Value  float64
}

func getHighest(num int, val *float64) *PriAndValue {
	return &PriAndValue{
		Number: num,
		Value:  float.NonNil(val),
	}
}

func GetHighest(payload *PriArray) *PriAndValue {
	if payload.P1 != nil {
		return getHighest(1, payload.P1)
	}
	if payload.P2 != nil {
		return getHighest(2, payload.P2)
	}
	if payload.P3 != nil {
		return getHighest(3, payload.P3)
	}
	return nil

}

func CleanArray(payload string) *PriArray {
	payload = strings.ReplaceAll(payload, "{", "")
	payload = strings.ReplaceAll(payload, "}", "")
	parts := strings.Split(payload, ",")
	if len(parts) != 16 {
		return nil
	}

	arr := &PriArray{
		P1:  set(parts[0]),
		P2:  set(parts[1]),
		P3:  set(parts[2]),
		P4:  set(parts[3]),
		P5:  set(parts[4]),
		P6:  set(parts[5]),
		P7:  set(parts[6]),
		P8:  set(parts[7]),
		P9:  set(parts[8]),
		P10: set(parts[9]),
		P11: set(parts[10]),
		P12: set(parts[11]),
		P13: set(parts[12]),
		P14: set(parts[13]),
		P15: set(parts[14]),
		P16: set(parts[15]),
	}
	return arr
}
