package ping

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/notify"
	"time"
)

const (
	pingNode = "ping"
)

type Ping struct {
	*node.Spec
	firstLoop bool
	triggered bool
	loopCount uint64
	lastOk    time.Time
	lastFail  time.Time
	lastPing  bool
}

func NewPing(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, pingNode, notify.Category)

	ip := node.BuildInput(node.Ip, node.TypeString, nil, body.Inputs)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs)
	body.Inputs = node.BuildInputs(ip, trigger)

	ok := node.BuildOutput(node.Ok, node.TypeBool, nil, body.Outputs)
	msg := node.BuildOutput(node.Msg, node.TypeString, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(ok, msg)
	return &Ping{body, false, false, 0, time.Now(), time.Now(), false}, nil
}

func (inst *Ping) ping(ip string, trigger bool) {

	if trigger && !inst.triggered {
		inst.triggered = true
		ok := helpers.CommandPing(ip)
		if ok {
			inst.lastOk = time.Now()
			inst.triggered = false
			inst.lastPing = true
		} else {
			inst.lastFail = time.Now()
			inst.triggered = false
			inst.lastPing = false

		}
	}
	if !trigger && inst.triggered {
		inst.triggered = false
	}
}

func (inst *Ping) Process() {
	ip := inst.ReadPinAsString(node.Ip)
	firstLoop, trigger := inst.InputUpdated(node.TriggerInput)
	if firstLoop || trigger {
		go inst.ping(ip, true)
	}
	pingMsg := helpers.PingMessage(ip, inst.lastPing, inst.lastOk)
	pingFailMsg := helpers.PingMessage(ip, inst.lastPing, inst.lastFail)
	fmt.Println(inst.lastOk, inst.lastFail)
	if inst.lastPing {
		inst.WritePin(node.Msg, pingMsg)
	} else {
		inst.WritePin(node.Msg, pingFailMsg)
	}
	inst.WritePin(node.Ok, inst.lastPing)

}

func (inst *Ping) Cleanup() {}
