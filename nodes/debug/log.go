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
	b := node.BuildInput(node.InBoolean, node.TypeBool, nil, body.Inputs)
	inputs := node.BuildInputs(comment, num, str, b)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	return &Log{body}, nil
}

func (inst *Log) Process() {
	comment := inst.ReadPinAsString(node.Comment)
	inNum, nullNum := inst.ReadPinAsFloatOk(node.InNumber)
	inStr := inst.ReadPin(node.InString)
	inBool := inst.ReadPin(node.InBoolean)
	if nullNum {
		log.Infof("log: comment: %v number: null", comment)
	} else {
		log.Infof("log: comment: %v number: %v", comment, inNum)
	}
	if inStr != nil {
		str := fmt.Sprintf("%v", inStr)
		log.Infof("log: comment: %v string: %s", comment, str)
	}
	if inst.InputHasConnection(node.InBoolean) {
		log.Infof("log: comment: %s bool: %t", comment, inBool)
	}
}

func (inst *Log) Cleanup() {}
