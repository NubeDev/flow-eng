package node

import (
	"github.com/NubeDev/flow-eng/helpers/conversions"
	log "github.com/sirupsen/logrus"
)

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
func (n *Spec) WritePinFloat(name OutputName, value float64, precision ...int) {
	out := n.GetOutput(name)
	if out == nil {
		return
	}
	if name == out.Name {
		var p = 2
		if len(precision) > 0 {
			p = precision[0]
		}
		out.Write(conversions.FloatToFixed(value, p))
	}
}

func (n *Spec) WritePinFalse(name OutputName) {
	out := n.GetOutput(name)
	if out == nil {
		log.Errorf("failed to find node to write oputput value FALSE node: %s", n.GetName())
		return
	}
	if name == out.Name {
		out.Write(false)
	}
}

func (n *Spec) WritePinTrue(name OutputName) {
	out := n.GetOutput(name)
	if out == nil {
		log.Errorf("failed to find node to write oputput value TRUE node: %s", n.GetName())
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
