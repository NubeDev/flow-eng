package node

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/conversions"
)

type Payload struct {
	Any any `json:"any,omitempty"`
}

func (n *Spec) SetPayload(body *Payload) {
	n.Payload = body
}
func (n *Spec) GetPayload() *Payload {
	return n.Payload
}

func (n *Spec) GetPayloadNull() (value any, null bool) {
	if n.Payload == nil {
		return nil, true
	}
	return n.Payload.Any, false
}

func (n *Spec) ReadPayloadAsFloat() (value float64, null bool) {
	r := n.GetPayload()
	if r == nil {
		return 0, true
	}
	if r.Any == nil {
		return 0, true
	}
	return conversions.GetFloat(r.Any), false
}

func (n *Spec) ReadPayloadAsString() (value string, null bool) {
	r := n.GetPayload()
	if r == nil {
		return "", true
	}
	if r.Any == nil {
		return "", true
	}
	return fmt.Sprint(r.Any), false
}
