package node

// SetName so we can easily set all names as upper or lower when needed
func SetName(name string) string {
	return name
}

func ConvertToBase(n Node) *BaseNode {
	return &BaseNode{
		Inputs:   n.GetInputs(),
		Outputs:  n.GetOutputs(),
		Info:     n.GetInfo(),
		Settings: n.GetSettings(),
		Metadata: n.GetMetadata(),
	}
}
