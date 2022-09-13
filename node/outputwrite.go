package node

func (n *BaseNode) WritePin(name OutputName, value interface{}) {
	out := n.GetOutput(name)
	if out == nil {
		return
	}
	if name == out.Name {
		out.Write(value)
	}
}

func (n *BaseNode) WritePinNum(name OutputName, value float64) {
	n.WritePin(name, value)
}
