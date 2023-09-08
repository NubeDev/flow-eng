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
	Category = "debug"
	logNode  = "log"
)

func NewLog(body *node.Spec, _ ...interface{}) (node.Node, error) {
	body = node.Defaults(body, logNode, Category)
	comment := node.BuildInput(node.Comment, node.TypeString, nil, body.Inputs, false, false)
	num := node.BuildInput(node.InNumber, node.TypeFloat, nil, body.Inputs, false, false)
	str := node.BuildInput(node.InString, node.TypeString, nil, body.Inputs, false, false)
	b := node.BuildInput(node.InBoolean, node.TypeBool, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(comment, num, str, b)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Log{body}, nil
}

func (inst *Log) Process() {
	comment, _ := inst.ReadPinAsString(node.Comment)
	inNum, nullNum := inst.ReadPinAsFloat(node.InNumber)
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
		log.Infof("log: comment: %s boolean: %t", comment, inBool)
	}
	inst.WritePinFloat(node.Out, 0)
}
