package node

type nodeValue struct {
	Pin      string      `json:"pin"`
	DataType DataTypes   `json:"dataType"`
	Value    interface{} `json:"value"`
}

type Values struct {
	NodeName string                 `json:"name"`
	NodeID   string                 `json:"nodeId"`
	Settings map[string]interface{} `json:"settings,omitempty"`
	Outputs  []*nodeValue           `json:"outputs"`
	Inputs   []*nodeValue           `json:"inputs"`
}

// NodeValues get the node current values
func (n *Spec) NodeValues() *Values {
	var out = &Values{
		NodeName: n.GetName(),
		NodeID:   n.GetID(),
		Settings: n.GetSettings(),
	}
	for _, output := range n.Outputs {
		v := &nodeValue{
			Pin:      string(output.Name),
			DataType: output.DataType,
			Value:    output.GetValue(),
		}
		out.Outputs = append(out.Outputs, v)
	}
	for _, input := range n.Inputs {
		v := &nodeValue{
			Pin:      string(input.Name),
			DataType: input.DataType,
			Value:    input.GetValue(),
		}
		out.Inputs = append(out.Inputs, v)
	}
	return out
}
