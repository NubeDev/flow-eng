package nodes

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/mitchellh/mapstructure"
)

// Decode the flow from the UI in to the node.Spec
func Decode(encodedNodes *NodesList) ([]*node.Spec, error) {
	var decodedNodes []*node.Spec
	var decodedNode *node.Spec
	for _, encodedNode := range encodedNodes.Nodes {
		ins := &node.SchemaInputs{}
		err := mapstructure.Decode(encodedNode.Inputs, ins)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		_, getName, _ := decodeType(encodedNode)
		id := encodedNode.Id
		name := getName
		decodedNode = node.New(id, name, "", encodedNode.Metadata) // create a blank node
		newNode, err := Builder(decodedNode)                       // make the new node as per its type
		for _, input := range newNode.GetInputs() {                // add the input connections as required
			for inputName, links := range ins.Links { // these would be the input connections
				if input.Name == node.InputName(inputName) {
					if links.Value != nil { // user has set a value and no input is connected
						str := fmt.Sprintf("%v", links.Value)
						input.Connection.OverrideValue = float.StrToFloat(str) // TODO add in dataTypes later
					} else {
						input.Connection.NodeID = links.NodeId
						input.Connection.NodePort = links.Socket
					}
				}
			}
		}
		decodedNodes = append(decodedNodes, decodedNode)

	}
	pprint.PrintJOSN(decodedNodes)
	return decodedNodes, nil

}
