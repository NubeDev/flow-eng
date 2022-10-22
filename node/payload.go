package node

import (
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/helpers/float"
)

type Payload struct {
	ValueFloat *float64 `json:"float,omitempty"`
	BoolFloat  *bool    `json:"boolean,omitempty"`
	String     *string  `json:"string,omitempty"`
	Any        any      `json:"any,omitempty"`
}

func (n *Spec) SetPayload(body *Payload) {
	n.Payload = body
}
func (n *Spec) GetPayload() *Payload {
	return n.Payload
}

func (n *Spec) ReadPayloadAsFloat() (value float64, null bool) {
	r := n.GetPayload()
	if r == nil {
		return 0, true
	}
	if r.ValueFloat == nil {
		return 0, true
	}
	return conversions.GetFloat(float.NonNil(r.ValueFloat)), false
}
