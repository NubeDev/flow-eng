package node

import "github.com/NubeDev/flow-eng/helpers"

// SetName so we can easily set all names as upper or lower when needed
func SetName(name string) string {
	return name
}

func ConvertToSpec(n Node) *Spec {
	return &Spec{
		Inputs:   n.GetInputs(),
		Outputs:  n.GetOutputs(),
		Info:     n.GetInfo(),
		Settings: n.GetSettings(),
		Metadata: n.GetMetadata(),
	}
}

func SetUUID(uuid string) string {
	if uuid == "" {
		uuid = helpers.ShortUUID("node")
	}
	return uuid
}
