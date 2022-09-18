package node

import (
	"github.com/NubeDev/flow-eng/helpers"
	"strings"
)

func SetInputName(name InputName) string {
	return strings.ToLower(string(name))
}

func SetOutputName(name OutputName) string {
	return strings.ToLower(string(name))
}

var ABCs = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}

// SetName so we can easily set all names as upper or lower when needed
func SetName(name string) string {
	return name
}

func ConvertToSpec(n Node) *Spec {
	if n == nil {
		return nil
	}
	return &Spec{
		Inputs:     n.GetInputs(),
		Outputs:    n.GetOutputs(),
		Info:       n.GetInfo(),
		Settings:   n.GetSettings(),
		Metadata:   n.GetMetadata(),
		Parameters: n.GetParameters(),
		SubFlow:    n.GetSubFlow(),
	}
}

func SetUUID(uuid string) string {
	if uuid == "" {
		uuid = helpers.ShortUUID("node")
	}
	return uuid
}
