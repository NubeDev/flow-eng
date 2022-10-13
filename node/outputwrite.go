package node

import "github.com/NubeDev/flow-eng/helpers/conversions"

func (n *Spec) WritePin(name OutputName, value interface{}) {
	out := n.GetOutput(name)
	if out == nil {
		return
	}
	if name == out.Name {
		out.Write(value)
	}
}

func (n *Spec) WritePinNull(name OutputName) {
	out := n.GetOutput(name)
	if out == nil {
		return
	}
	if name == out.Name {
		out.Write(nil)
	}
}

// WritePinFloat write a float64 pointer
func (n *Spec) WritePinFloat(name OutputName, value float64) {
	out := n.GetOutput(name)
	if out == nil {
		return
	}
	if name == out.Name {
		out.Write(conversions.FloatToFixed(value, 2))
	}
}

func (n *Spec) WritePinFalse(name OutputName) {
	out := n.GetOutput(name)
	if out == nil {
		return
	}
	if name == out.Name {
		out.Write(false)
	}
}

func (n *Spec) WritePinTrue(name OutputName) {
	out := n.GetOutput(name)
	if out == nil {
		return
	}
	if name == out.Name {
		out.Write(true)
	}
}

func (n *Spec) WritePinBool(name OutputName, value bool) {
	out := n.GetOutput(name)
	if out == nil {
		return
	}
	if name == out.Name {
		out.Write(value)
	}
}
