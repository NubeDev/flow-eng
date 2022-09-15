package node

func (n *Spec) WritePin(name OutputName, value interface{}) {
	out := n.GetOutput(name)
	if out == nil {
		return
	}
	if name == out.Name {
		out.Write(value)
	}
}
