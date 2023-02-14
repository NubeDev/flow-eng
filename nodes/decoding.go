package nodes

import (
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
)

// Decode the flow from react-flow to the flow-eng in correct format
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
		decodedNode.Info.Icon = encodedNode.Icon
		decodedNode.Info.NodeName = encodedNode.NodeName
		if encodedNode.Payload != nil {
			decodedNode.Payload = &node.Payload{Any: encodedNode.Payload}
		}
		newNode, err := Builder(decodedNode, nil, nil)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		for _, input := range newNode.GetInputs() { // add the input connections as required
			for inputName, links := range encodedNode.Inputs { // these would be the input connections
				input.Position = links.Position
				input.OverridePosition = links.OverridePosition
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

type NodeFilter string

const (
	FilterIsParent      NodeFilter = "FilterIsParent"
	FilterIsParentChild NodeFilter = "FilterIsParentChild" // when a node is a parent under another parent
	FilterIsChild       NodeFilter = "FilterIsChild"
	FilterNonContainer  NodeFilter = "FilterNonContainer" // when a node is not in a container node
)

// FilterNodes filter nodes that are a container node or inside a container node
// 	- parentID:string if you want to filter nodes that are inside a container pass in the node container nodeID
func FilterNodes(nodesList []*node.Spec, filter NodeFilter, parentID string) []*node.Spec {
	var out []*node.Spec
	if filter == FilterIsParent {
		for _, n := range nodesList {
			if n.IsParent {
				if n.ParentId == "" {
					out = append(out, n)
				}
			}
		}
	}
	if filter == FilterIsParentChild {
		for _, n := range nodesList {
			if n.IsParent {
				if n.ParentId != "" {
					out = append(out, n)
				}
			}
		}
	}

	if filter == FilterIsChild {
		for _, n := range nodesList {
			if !n.IsParent {
				if n.ParentId != "" {
					out = append(out, n)
				}
			}
		}
	}

	if filter == FilterNonContainer {
		for _, n := range nodesList {
			if !n.IsParent {
				if n.ParentId == "" {
					out = append(out, n)
				}
			}
		}
	}

	if parentID != "" { // get all the nodes that are inside a container
		for _, n := range out {
			if n.ParentId == parentID {
				out = append(out, n)
			}
		}
		return out
	}
	return out
}
