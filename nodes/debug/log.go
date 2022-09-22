package debugging

import (
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
)

type Log struct {
	*node.Spec
}

const (
	category = "debug"
	logNode  = "log"
)

func NewLog(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, logNode, category)
	inputs := node.BuildInputs(node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs))
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Log{body}, nil
}

func (inst *Log) Process() {
	in1 := inst.ReadPin(node.In)
	if in1 != nil {
		log.Infof("log: node:%s value: %v", inst.Info.Name, in1)
	} else {
		log.Infof("log: node:%s no value", inst.Info.Name)
	}

}

func (inst *Log) Cleanup() {}