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
}

func NewPing(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, pingNode, notify.Category)

	ip := node.BuildInput(node.Ip, node.TypeString, nil, body.Inputs)
	t := node.BuildInput(node.Time, node.TypeString, nil, body.Inputs)
	trigger := node.BuildInput(node.TriggerInput, node.TypeFloat, nil, body.Inputs)
	body.Inputs = node.BuildInputs(ip, t, trigger)

	msg := node.BuildOutput(node.Result, node.TypeString, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(msg)
	return &Ping{body, false, false, 0, time.Now()}, nil
}

func (inst *Ping) ping() {
	if !inst.firstLoop {
		//inst.setEmailClient()
	}
	trigger := inst.ReadPinBool(node.TriggerInput)
	if trigger && !inst.triggered {
		inst.triggered = true
		ip := inst.ReadPinAsString(node.Ip)
		fmt.Println("TRIGGER PING", ip)
		ok := helpers.CommandPing(ip)
		if ok {
			inst.lastOk = time.Now()
			ping := helpers.PingMessage(ip, ok, inst.lastOk)
			fmt.Println(ping)
			inst.triggered = false
		} else {
			ping := helpers.PingMessage(ip, ok, inst.lastOk)
			fmt.Println(ping)
			inst.triggered = false
		}

	}
	if !trigger && inst.triggered {
		fmt.Println("RESET")
		inst.triggered = false
	}
}

func (inst *Ping) Process() {
	inst.loopCount++
	go inst.ping()

}

func (inst *Ping) Cleanup() {}
