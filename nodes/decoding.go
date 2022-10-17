package nodes

import (
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
)

func Decode(encodedNodes *NodesList) ([]*node.Spec, error) {
	var decodedNodes []*node.Spec
	for _, encodedNode := range encodedNodes.Nodes {
		var decodedNode *node.Spec
		_, getName, _ := decodeType(encodedNode.Type)
		id := encodedNode.Id
		name := getName
		decodedNode = node.New(id, name, "", encodedNode.Metadata, encodedNode.Settings) // create a blank node
		decodedNode.IsParent = encodedNode.IsParent
		decodedNode.ParentId = encodedNode.ParentId
		newNode, err := Builder(decodedNode, nil)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		for _, input := range newNode.GetInputs() { // add the input connections as required
			for inputName, links := range encodedNode.Inputs { // these would be the input connections
				if input.Name == node.InputName(inputName) {
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
	return decodedNodes, nil
}
