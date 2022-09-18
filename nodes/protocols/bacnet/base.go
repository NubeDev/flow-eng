package bacnet

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/mqttbase"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bstore"
)

const (
	category      = "bacnet"
	server        = "server"
	bacnetReadBV  = "binary-variable-read"
	bacnetWriteBV = "binary-variable-write"
	bacnetReadAV  = "analog-variable-read"
	bacnetWriteAV = "analog-variable-write"

	typeBV = "bv"
	typeAV = "av"
)

func getBacnetType(nodeName string) (bstore.ObjectType, error) {
	switch nodeName {
	case bacnetReadBV:
		return bstore.BinaryVariable, nil
	}
	return "", errors.New(fmt.Sprintf("bacnet add new point object type not found node: %s", nodeName))
}

func nodeDefault(body *node.Spec, nodeName, category string, application node.ApplicationName, opts interface{}) (*node.Spec, *mqttbase.Mqtt, error, *bstore.Point) {
	var err error
	body = node.Defaults(body, nodeName, category)

	objectType, err := getBacnetType(nodeName)

	pointName := node.BuildInput(node.Name, node.TypeString, nil, body.Inputs)
	objectID := node.BuildInput(node.ObjectId, node.TypeFloat, 1, body.Inputs)
	overrideInput := node.BuildInput(node.OverrideInput, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(pointName, objectID, overrideInput)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))

	parameters := &node.Parameters{
		Application: &node.Application{
			Application: applications.BACnet,
			IsChild:     true,
		},
	}
	body.Parameters = node.BuildParameters(parameters)
	body = node.BuildNode(body, inputs, outputs, nil)
	client := &mqttbase.Mqtt{}
	var ok bool
	client, ok = opts.(*mqttbase.Mqtt)
	fmt.Println(ok, 77777)

	objId, ok := objectID.GetValue().(float64)
	fmt.Println(objId, ok, 888888)

	point := &bstore.Point{
		Application: application,
		ObjectType:  objectType,
		ObjectID:    int(objId),
		IoType:      "",
		IoNumber:    0,
	}

	return body, client, err, point
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
func topicBuilder(objectType string, address int) string {
	return fmt.Sprintf("bacnet/%s/%d", objectType, address)
}

// TopicPresentValue bacnet/ao/1/pv
func TopicPresentValue(objectType string, address int) string {
	return fmt.Sprintf("%s/pv", topicBuilder(objectType, address))
}

// TopicPriority bacnet/ao/1/pri
func TopicPriority(objectType string, address int) string {
	return fmt.Sprintf("%s/pri", topicBuilder(objectType, address))
}

func GetObjectId(address float64) int {
	return int(address)
}
