package bacnet

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/NubeDev/flow-eng/node"
	points "github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
)

const (
	category = "bacnet"
	server   = "server"
	bacnetBI = "binary-input"
	bacnetBO = "binary-output"
	bacnetBV = "binary-variable"
	bacnetAV = "analog-variable"
	bacnetAO = "analog-output"
	bacnetAI = "analog-input"

	typeAI = "ai"
	typeAO = "ao"
	typeAV = "av"
	typeBV = "bv"
)

func getBacnetType(nodeName string) (obj points.ObjectType, isWriteable, isIO bool, err error) {
	switch nodeName {
	case bacnetAI:
		return points.AnalogInput, false, true, nil
	case bacnetAO:
		return points.AnalogOutput, true, true, nil
	case bacnetAV:
		return points.AnalogVariable, true, false, nil
	case bacnetBI:
		return points.BinaryInput, false, true, nil
	case bacnetBO:
		return points.BinaryOutput, true, true, nil
	case bacnetBV:
		return points.BinaryVariable, true, false, nil

	}
	return "", false, false, errors.New(fmt.Sprintf("bacnet add new point object type not found node: %s", nodeName))
}

func nodeDefault(body *node.Spec, nodeName, category string, application names.ApplicationName) (*node.Spec, error) {
	var err error
	body = node.Defaults(body, nodeName, category)
	_, isWriteable, _, err := getBacnetType(nodeName)
	pointName := node.BuildInput(node.Name, node.TypeString, nil, body.Inputs)
	objectIDInput := node.BuildInput(node.ObjectId, node.TypeFloat, 0, body.Inputs)
	var inputs []*node.Input
	if isWriteable {
		in14 := node.BuildInput(node.In14, node.TypeFloat, nil, body.Inputs)
		in15 := node.BuildInput(node.In15, node.TypeFloat, nil, body.Inputs)
		inputs = node.BuildInputs(pointName, objectIDInput, in14, in15)
	} else {
		overrideInput := node.BuildInput(node.OverrideInput, node.TypeFloat, nil, body.Inputs)
		inputs = node.BuildInputs(pointName, objectIDInput, overrideInput)
	}

	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))

	parameters := &node.Parameters{
		Application: &node.Application{
			Application: names.BACnet,
			IsChild:     true,
		},
	}
	body.Parameters = node.BuildParameters(parameters)
	body = node.BuildNode(body, inputs, outputs, nil)
	return body, err
}

func addPoint(application names.ApplicationName, ioType points.IoType, objectType points.ObjectType, id points.ObjectID, isWriteable, isIO, enable bool) *points.Point {
	point := &points.Point{
		Application: application,
		ObjectType:  objectType,
		ObjectID:    id,
		IoType:      ioType,
		IsIO:        isIO,
		IsWriteable: isWriteable,
		Enable:      enable,
	}
	return point

}

func nodeDefaultOld(body *node.Spec, nodeName, category string, application names.ApplicationName) (*node.Spec, error, *points.Point) {
	var err error
	body = node.Defaults(body, nodeName, category)
	objectType, isWriteable, isIO, err := getBacnetType(nodeName)
	pointName := node.BuildInput(node.Name, node.TypeString, nil, body.Inputs)
	objectIDInput := node.BuildInput(node.ObjectId, node.TypeFloat, 0, body.Inputs)
	//fmt.Println(11111, objectType)
	//pprint.PrintJOSN(objectIDInput)
	ioType := points.IoTypeDigital // TODO make a setting
	if isWriteable {
		ioType = points.IoTypeDigital
	}

	enable := true // TODO make a setting
	var inputs []*node.Input

	if isWriteable {
		in14 := node.BuildInput(node.In14, node.TypeFloat, nil, body.Inputs)
		in15 := node.BuildInput(node.In15, node.TypeFloat, nil, body.Inputs)
		inputs = node.BuildInputs(pointName, objectIDInput, in14, in15)
	} else {
		overrideInput := node.BuildInput(node.OverrideInput, node.TypeFloat, nil, body.Inputs)
		inputs = node.BuildInputs(pointName, objectIDInput, overrideInput)
	}

	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))

	parameters := &node.Parameters{
		Application: &node.Application{
			Application: names.BACnet,
			IsChild:     true,
		},
	}
	body.Parameters = node.BuildParameters(parameters)
	body = node.BuildNode(body, inputs, outputs, nil)

	objectID := conversions.GetInt(objectIDInput.GetValue())
	if objectID == 0 {
		log.Errorf("bacnet-server object-id must be grater then 0 object-type:%s", objectType)
		objectID = 1
	}
	point := addPoint(application, ioType, objectType, points.ObjectID(objectID), isWriteable, isIO, enable)
	store := getStore()
	point, err = store.AddPoint(point, true)
	if err != nil {
		log.Errorf("bacnet-server add new point type:%s-%d", objectType, points.ObjectID(objectID))
	}
	return body, err, point
}

// topicBuilder bacnet/ObjectType
func topicObjectBuilder(objectType string) string {
	return fmt.Sprintf("bacnet/%s", objectType)
}

// topicBuilder bacnet/ao/1
func topicBuilder(objectType string, address points.ObjectID) string {
	return fmt.Sprintf("bacnet/%s/%d", objectType, address)
}

// TopicPresentValue bacnet/ao/1/pv
func TopicPresentValue(objectType string, address points.ObjectID) string {
	return fmt.Sprintf("%s/pv", topicBuilder(objectType, address))
}

// TopicPriority bacnet/ao/1/pri
func TopicPriority(objectType string, address points.ObjectID) string {
	return fmt.Sprintf("%s/pri", topicBuilder(objectType, address))
}
