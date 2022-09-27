package debugging

import (
	"fmt"
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
	comment := node.BuildInput(node.Comment, node.TypeString, nil, body.Inputs)
	num := node.BuildInput(node.InNumber, node.TypeFloat, nil, body.Inputs)
	str := node.BuildInput(node.InString, node.TypeString, nil, body.Inputs)
	inputs := node.BuildInputs(comment, num, str)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Log{body}, nil
}

func (inst *Log) Process() {
	comment := inst.ReadPin(node.Comment)
	inNum := inst.ReadPin(node.InNumber)
	inStr := inst.ReadPin(node.InString)

	if inNum != nil {
		log.Infof("log: comment: %v number: %v", comment, inNum)
	}
	str := fmt.Sprintf("%v", inStr)
	if str != "" {
		log.Infof("log: comment: %v string: %s", comment, str)
	}

}

func (inst *Log) Cleanup() {}
