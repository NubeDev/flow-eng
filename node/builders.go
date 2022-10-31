package node

import (
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/enescakir/emoji"
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

func BuildNode(body *Spec, inputs []*Input, outputs []*Output, settings map[string]interface{}) *Spec {
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

const noParent = "please add node to a flow-network node"

func SetNoParent(body *Spec) *Spec {
	if body.ParentId == "" {
		body.SetStatusError(noParent)
		body.SetErrorIcon(string(emoji.RedCircle))
	}
	return body
}

func SetError(body *Spec, message string) *Spec {
	body.SetStatusError(message)
	return body
}

func SetStatus(body *Spec, message string) *Spec {
	body.SetStatusMessage(message)
	return body
}
