package nodes

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
)

func Decode(encodedNodes *NodesList) ([]*node.Spec, error) {
	var decodedNodes []*node.Spec
	for _, encodedNode := range encodedNodes.Nodes {
		var decodedNode *node.Spec
		_, getName, _ := decodeType(encodedNode)
		id := encodedNode.Id
		name := getName
		decodedNode = node.New(id, name, "", encodedNode.Metadata) // create a blank node
		newNode, err := Builder(decodedNode)
		if err != nil {
			return nil, err
		}
		for _, input := range newNode.GetInputs() { // add the input connections as required
			for inputName, links := range encodedNode.Inputs { // these would be the input connections
				if input.Name == node.InputName(inputName) {
					fmt.Println(name, links.Value)
					if links.Value != nil { // user has set a value and no input is connected
						input.Connection.OverrideValue = links.Value
					} else {
						for _, link := range links.Links {
							input.Connection.NodeID = link.NodeId
							input.Connection.NodePort = node.OutputName(link.Socket)
						}

					}
				}
			}
		}
		decodedNodes = append(decodedNodes, decodedNode)
	}
	pprint.PrintJOSN(decodedNodes)
	return decodedNodes, nil
}

// DecodeNonSubNodes the flow from the UI in to the node.Spec
func DecodeNonSubNodes(encodedNodes *NodesList) ([]*node.Spec, error) {
	var decodedNodes []*node.Spec
	for _, encodedNode := range encodedNodes.Nodes {
		var decodedNode *node.Spec
		_, getName, _ := decodeType(encodedNode)
		id := encodedNode.Id
		name := getName
		decodedNode = node.New(id, name, "", encodedNode.Metadata) // create a blank node
		newNode, err := Builder(decodedNode)
		if err != nil {
			return nil, err
		}
		for _, input := range newNode.GetInputs() { // add the input connections as required
			for inputName, links := range encodedNode.Inputs { // these would be the input connections
				if input.Name == node.InputName(inputName) {
					if links.Value != nil { // user has set a value and no input is connected
						str := fmt.Sprintf("%v", links.Value, input.Name)
						input.Connection.OverrideValue = float.StrToFloat(str) // TODO add in dataTypes later
					} else {
						for _, link := range links.Links {
							input.Connection.NodeID = link.NodeId
							input.Connection.NodePort = node.OutputName(link.Socket)
						}

					}
				}
			}
		}
		decodedNodes = append(decodedNodes, decodedNode)
	}
	pprint.PrintJOSN(decodedNodes)
	return decodedNodes, nil

}
