package math

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
)

const (
	category = "math"
	add      = "add"
	sub      = "subtract"
	multiply = "multiply"
)

func Process(body *node.BaseNode) {
	_, in1Val, in1Not := body.ReadPinNum(node.In1)
	_, in2Val, in2Not := body.ReadPinNum(node.In2)
	if in1Not && in2Not {
		f, err := mathOperation(body.GetName(), in1Val, in2Val)
		if err != nil {
			return
		}
		fmt.Println("MATH", f, body.GetName())
		body.WritePinNum(node.Out1, f)
		return
	}
	if in1Not {
		body.WritePinNum(node.Out1, in1Val)
		return
	}
	if in2Not {
		body.WritePinNum(node.Out1, in2Val)
		return
	}
}

func mathOperation(operation string, firstValue float64, values ...float64) (float64, error) {
	switch operation {
	case add:
		var out float64
		var length = len(values)
		for i, value := range values {
			firstValue = firstValue + value
			if i == length-1 {
				out = firstValue
			}
		}
		return out, nil
	case sub:
		var out float64
		var length = len(values)
		for i, value := range values {
			firstValue = firstValue - value
			if i == length-1 {
				out = firstValue
			}
		}
		return out, nil
	case multiply:
		var out float64
		var length = len(values)
		for i, value := range values {
			firstValue = firstValue * value
			if i == length-1 {
				out = firstValue
			}
		}
		return out, nil
	}
	return 0, errors.New("invalid operation")
}
