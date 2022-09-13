package node

func (n *BaseNode) ReadPin(name InputName) interface{} {
	input := n.GetInput(name)
	if input == nil {
		return nil
	}
	if input.Connection.OverrideValue != nil { // this would be that the user wrote a value to the input directly
		return input.Connection.OverrideValue
	}
	if input.Connection.FallbackValue != nil { // this would be that the user wrote a value to the input directly
		return input.Connection.FallbackValue
	}
	return input.GetValue()
}

func (n *BaseNode) ReadMultiple(count int) []interface{} {
	var out []interface{}
	for i, input := range n.GetInputs() {
		if i < count {
			out = append(out, n.ReadPin(input.Name))
		}
	}
	return out
}
