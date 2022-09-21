package bacnet

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	points2 "github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
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

func getBacnetType(nodeName string) (obj points2.ObjectType, isWriteable, isIO bool, err error) {
	switch nodeName {
	case bacnetAI:
		return points2.AnalogInput, false, true, nil
	case bacnetAO:
		return points2.AnalogOutput, true, true, nil
	case bacnetAV:
		return points2.AnalogVariable, true, false, nil
	case bacnetBI:
		return points2.BinaryInput, false, true, nil
	case bacnetBO:
		return points2.BinaryOutput, true, true, nil
	case bacnetBV:
		return points2.BinaryVariable, true, false, nil

	}
	return "", false, false, errors.New(fmt.Sprintf("bacnet add new point object type not found node: %s", nodeName))
}

func nodeDefault(body *node.Spec, nodeName, category string, application node.ApplicationName) (*node.Spec, error, *points2.Point) {
	var err error
	body = node.Defaults(body, nodeName, category)

	objectType, isWriteable, isIO, err := getBacnetType(nodeName)

	pointName := node.BuildInput(node.Name, node.TypeString, nil, body.Inputs)
	objectIDInput := node.BuildInput(node.ObjectId, node.TypeFloat, 1, body.Inputs)
	ioType := points2.IoTypeTemp // TODO make a setting
	if isWriteable {
		ioType = points2.IoTypeVolts
	}

	enable := true // TODO make a setting
	var inputs []*node.Input

	if isWriteable {
		overrideInput := node.BuildInput(node.In16, node.TypeFloat, nil, body.Inputs)
		inputs = node.BuildInputs(pointName, objectIDInput, overrideInput)
	} else {
		overrideInput := node.BuildInput(node.OverrideInput, node.TypeFloat, nil, body.Inputs)
		inputs = node.BuildInputs(pointName, objectIDInput, overrideInput)
	}

	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))

	parameters := &node.Parameters{
		Application: &node.Application{
			Application: applications.BACnet,
			IsChild:     true,
		},
	}
	body.Parameters = node.BuildParameters(parameters)
	body = node.BuildNode(body, inputs, outputs, nil)

	objectID, _ := objectIDInput.GetValue().(float64)
	if objectID == 0 {
		log.Errorf("bacnet-server object-id must be grater then 0 object-type:%s", objectType)
		objectID = 1
	}

	point := addPoint(application, ioType, objectType, points2.ObjectID(objectID), isWriteable, isIO, enable)
	store := getStore()
	point, err = store.AddPoint(point)
	log.Infof("bacnet-server add new point type:%s-%d", point.ObjectType, point.ObjectID)
	return body, err, point
}

func addPoint(application node.ApplicationName, ioType points2.IoType, objectType points2.ObjectType, id points2.ObjectID, isWriteable, isIO, enable bool) *points2.Point {
	point := &points2.Point{
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

// topicBuilder bacnet/ObjectType
func topicObjectBuilder(objectType string) string {
	return fmt.Sprintf("bacnet/%s", objectType)
}

// topicBuilder bacnet/ao/1
func topicBuilder(objectType string, address points2.ObjectID) string {
	return fmt.Sprintf("bacnet/%s/%d", objectType, address)
}

// TopicPresentValue bacnet/ao/1/pv
func TopicPresentValue(objectType string, address points2.ObjectID) string {
	return fmt.Sprintf("%s/pv", topicBuilder(objectType, address))
}

// TopicPriority bacnet/ao/1/pri
func TopicPriority(objectType string, address points2.ObjectID) string {
	return fmt.Sprintf("%s/pri", topicBuilder(objectType, address))
}
