package ping

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/notify"
	"time"
)

const (
	pingNode = "ping"
)

type result struct {
	Ip       string    `json:"ip"`
	Ok       bool      `json:"ok"`
	LastOk   time.Time `json:"lastOk"`
	LastFail time.Time `json:"lastFail"`
}

type Ping struct {
	*node.Spec
	triggered bool
	loopCount uint64
	lastOk    time.Time
	lastFail  time.Time
	lastPing  bool
	results   []result
}

func NewPing(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, pingNode, notify.Category)

	ip := node.BuildInput(node.Ip, node.TypeString, nil, body.Inputs, nil)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs, nil)
	body.Inputs = node.BuildInputs(ip, trigger)
	msg := node.BuildOutput(node.Outp, node.TypeString, nil, body.Outputs)
	body.Outputs = node.BuildOutputs(msg)
	var res []result
	return &Ping{body, false, 0, time.Now(), time.Now(), false, res}, nil
}

func (inst *Ping) ping(ipList []string) {
	inst.results = nil
	inst.triggered = true
	var r result
	for _, ip := range ipList {
		ok := helpers.CommandPing(ip)
		r.Ip = ip
		r.Ok = ok
		if ok {
			r.LastOk = time.Now()
		} else {
			r.LastFail = time.Now()
		}
		inst.results = append(inst.results, r)
	}
	inst.triggered = false
}

func (inst *Ping) Process() {
	firstLoop, trigger := inst.InputUpdated(node.TriggerInput)
	read := inst.ReadPin(node.Ip)
	var ipList []string
	err := json.Unmarshal([]byte(read.(string)), &ipList)
	if err != nil {
		ipList = append(ipList, read.(string))
	}

	if firstLoop || trigger {
		if !inst.triggered {
			go inst.ping(ipList)
		}
	}
	if len(inst.results) > 0 {
		out, _ := json.Marshal(inst.results)
		// value := json.ParseBytes(out)
		inst.WritePin(node.Outp, string(out))
	} else {
		inst.WritePinNull(node.Outp)
	}

}
