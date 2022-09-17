package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/mqttbase"
	"github.com/NubeDev/flow-eng/node"
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

func nodeDefault(body *node.Spec, nodeName, category string, opts interface{}) (*node.Spec, *mqttbase.Mqtt, error) {
	body = node.Defaults(body, nodeName, category)
	var inNames = []string{node.SetInputName(node.ObjectId), node.SetInputName(node.OverrideInput)}
	var inputsList []*node.Input
	pointName := node.BuildInput(node.In, node.TypeString, nil, body.Inputs)
	inputsList = append(inputsList, pointName)
	inputsList = append(inputsList, node.DynamicInputs(node.TypeFloat, nil, 2, 0, 0, body.Inputs, inNames)...)
	inputs := node.BuildInputs(inputsList...)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	client := &mqttbase.Mqtt{}
	var ok bool
	client, ok = opts.(*mqttbase.Mqtt)
	fmt.Println(ok, 9999999)

	//if len(opts) > 0 {
	//	for index, val := range opts {
	//		fmt.Println(index, val)
	//		client, ok = val.(*mqttbase.Mqtt)
	//		if !ok {
	//			fmt.Println("FUCK failed to make mqtt client")
	//		}
	//
	//	}
	//	//
	//	//client, ok = opts[0].(*mqttbase.Mqtt)
	//	//if !ok {
	//	//	fmt.Println("FUCK failed to make mqtt client")
	//	//}
	//}

	return body, client, nil
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
