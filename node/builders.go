package node

import (
	"github.com/NubeDev/flow-eng/helpers"
)

func BuildNodes(body ...*Spec) []*Spec {
	var out []*Spec
	for _, output := range body {
		if output != nil {
			out = append(out, output)
		}
	}
	return out
}

func BuildNode(body *Spec, inputs []*Input, outputs []*Output, settings []*Settings) *Spec {
	body.Settings = settings
	body.Inputs = inputs
	body.Outputs = outputs
	return body
}

func Defaults(body *Spec, nodeName, category string) *Spec {
	if body == nil {
		body = &Spec{
			Info: Info{
				NodeName: helpers.ShortUUID(nodeName),
				NodeID:   "",
			},
		}
	}
	body.Info.Name = SetName(nodeName)
	body.Info.Category = SetName(category)
	body.Info.NodeID = SetUUID(body.Info.NodeID)
	return body
}
