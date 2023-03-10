package node

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/boolean"
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

func (n *Spec) ReadPayloadAsBool() (value, noPayload, nullPayload bool) {
	r := n.GetPayload()
	if r == nil {
		return false, true, false
	}
	if r.Any == nil {
		return false, true, false
	}
	m, ok := r.Any.(map[string]interface{})
	if ok {
		for _, v := range m {
			if v == nil {
				return false, false, true
			} else {
				return boolean.NonNil(boolean.ConvertInterfaceToBool(v)), false, false
			}
		}
	}
	return false, true, false
}

func (n *Spec) ReadPayloadAsFloat() (value float64, noPayload, nullPayload bool) {
	r := n.GetPayload()
	if r == nil {
		return 0, true, false
	}
	if r.Any == nil {
		return 0, true, false
	}
	m, ok := r.Any.(map[string]interface{})
	if ok {
		for _, v := range m {
			if v == nil {
				return 0, false, true
			} else {
				return conversions.GetFloat(v), false, false
			}
		}
	}
	return 0, true, false
}

func (n *Spec) ReadMQTTPayloadAsString() (value string, null bool) {
	r := n.GetPayload()
	if r == nil {
		return "", true
	}
	if r.Any == nil {
		return "", true
	}
	return fmt.Sprint(r.Any), false
}

func (n *Spec) ReadPayloadAsString() (value string, noPayload, nullPayload bool) {
	r := n.GetPayload()
	if r == nil {
		return "", true, false
	}
	if r.Any == nil {
		return "", true, false
	}
	m, ok := r.Any.(map[string]interface{})
	if ok {
		for _, v := range m {
			if v == nil || v == "" || v == "null" {
				return "", false, true
			} else {
				val, k := conversions.GetStringOk(v)
				if k {
					return val, false, false
				}
			}
		}
	}
	return "", true, false
}
