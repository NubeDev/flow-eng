package bacnet

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bstore"
)

const (
	category = "bacnet"
	server   = "server"
	bacnetBV = "binary-variable"
	bacnetAV = "analog-variable"

	typeBV = "bv"
	typeAV = "av"
)

func getBacnetType(nodeName string) (obj bstore.ObjectType, isWriteable bool, err error) {
	switch nodeName {
	case bacnetBV:
		return bstore.BinaryVariable, true, nil
	}
	return "", false, errors.New(fmt.Sprintf("bacnet add new point object type not found node: %s", nodeName))
}

func nodeDefault(body *node.Spec, nodeName, category string, application node.ApplicationName) (*node.Spec, error, *bstore.Point) {
	var err error
	body = node.Defaults(body, nodeName, category)

	objectType, isWriteable, err := getBacnetType(nodeName)

	pointName := node.BuildInput(node.Name, node.TypeString, nil, body.Inputs)
	objectID := node.BuildInput(node.ObjectId, node.TypeFloat, 1, body.Inputs)

	var inputs []*node.Input

	if isWriteable {
		overrideInput := node.BuildInput(node.In16, node.TypeFloat, nil, body.Inputs)
		inputs = node.BuildInputs(pointName, objectID, overrideInput)
	} else {
		overrideInput := node.BuildInput(node.OverrideInput, node.TypeFloat, nil, body.Inputs)
		inputs = node.BuildInputs(pointName, objectID, overrideInput)
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

	objId, _ := objectID.GetValue().(float64)

	point := &bstore.Point{
		Application: application,
		ObjectType:  objectType,
		ObjectID:    bstore.ObjectID(int(objId)),
	}
	store := GetStore()
	point, err = store.AddPoint(point)
	return body, err, point
}

//func process(body node.Node) {
//	nodeName := body.GetName()
//
//	objectId := body.ReadPin(node.ObjectId)
//	overrideInput := body.ReadPin(node.OverrideInput)
//
//	body.WritePin(node.Out, output)
//}

// topicBuilder bacnet/ao/1
func topicBuilder(objectType string, address bstore.ObjectID) string {
	return fmt.Sprintf("bacnet/%s/%d", objectType, address)
}

// TopicPresentValue bacnet/ao/1/pv
func TopicPresentValue(objectType string, address bstore.ObjectID) string {
	return fmt.Sprintf("%s/pv", topicBuilder(objectType, address))
}

// TopicPriority bacnet/ao/1/pri
func TopicPriority(objectType string, address bstore.ObjectID) string {
	return fmt.Sprintf("%s/pri", topicBuilder(objectType, address))
}
