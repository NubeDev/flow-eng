package node

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"time"
)

type Payload struct {
	Any        any `json:"any,omitempty"`
	LastUpdate time.Time
}

func (n *Spec) SetPayload(body *Payload) {
	body.LastUpdate = time.Now()
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
	m, ok := r.Any.(map[string]interface{})
	if ok {
		for _, v := range m {
			return conversions.GetFloat(v), false
		}
	}
	return 0, true
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
