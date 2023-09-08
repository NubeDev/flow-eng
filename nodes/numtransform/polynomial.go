package numtransform

import (
	"math"

	"github.com/NubeDev/flow-eng/node"
)

type Polynomial struct {
	*node.Spec
}

func NewPolynomial(body *node.Spec, _ ...any) (node.Node, error) {
	body = node.Defaults(body, polynomial, Category)
	x := node.BuildInput(node.X, node.TypeFloat, nil, body.Inputs, false, false)
	x0 := node.BuildInput(node.X0, node.TypeFloat, nil, body.Inputs, false, false)
	x1 := node.BuildInput(node.X1, node.TypeFloat, nil, body.Inputs, false, false)
	x2 := node.BuildInput(node.X2, node.TypeFloat, nil, body.Inputs, false, false)
	x3 := node.BuildInput(node.X3, node.TypeFloat, nil, body.Inputs, false, false)
	x4 := node.BuildInput(node.X4, node.TypeFloat, nil, body.Inputs, false, false)
	x5 := node.BuildInput(node.X5, node.TypeFloat, nil, body.Inputs, false, false)
	// x6 := node.BuildInput(node.X6, node.TypeFloat, nil, body.Inputs, false, false)
	// x7 := node.BuildInput(node.X7, node.TypeFloat, nil, body.Inputs, false, false)
	// x8 := node.BuildInput(node.X8, node.TypeFloat, nil, body.Inputs, false, false)
	// x9 := node.BuildInput(node.X9, node.TypeFloat, nil, body.Inputs, false, false)
	// x10 := node.BuildInput(node.X10, node.TypeFloat, nil, body.Inputs, false, false)
	// inputs := node.BuildInputs(x, x0, x1, x2, x3, x4, x5, x6, x7, x8, x9, x10)

	inputs := node.BuildInputs(x, x0, x1, x2, x3, x4, x5)

	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Polynomial{body}, nil
}

func (inst *Polynomial) Process() {
	x, _ := inst.ReadPinAsFloat(node.X)
	x0, _ := inst.ReadPinAsFloat(node.X0)
	x1, _ := inst.ReadPinAsFloat(node.X1)
	x2, _ := inst.ReadPinAsFloat(node.X2)
	x3, _ := inst.ReadPinAsFloat(node.X3)
	x4, _ := inst.ReadPinAsFloat(node.X4)
	x5, _ := inst.ReadPinAsFloat(node.X5)
	// x6, _ := inst.ReadPinAsFloat(node.X6)
	// x7, _ := inst.ReadPinAsFloat(node.X7)
	// x8, _ := inst.ReadPinAsFloat(node.X8)
	// x9, _ := inst.ReadPinAsFloat(node.X9)
	// x10, _ := inst.ReadPinAsFloat(node.X10)

	output := (x5 * math.Pow(x, 5)) + (x4 * math.Pow(x, 4)) + (x3 * math.Pow(x, 3)) + (x2 * math.Pow(x, 2)) + (x1 * math.Pow(x, 1)) + (x0)
	// output := (x10 * math.Pow(x, 10)) + (x9 * math.Pow(x, 9)) + (x8 * math.Pow(x, 8)) + (x7 * math.Pow(x, 7)) + (x6 * math.Pow(x, 6)) + (x5 * math.Pow(x, 5)) + (x4 * math.Pow(x, 4)) + (x3 * math.Pow(x, 3)) + (x2 * math.Pow(x, 2)) + (x1 * math.Pow(x, 1)) + (x0)
	inst.WritePinFloat(node.Out, output)
}
