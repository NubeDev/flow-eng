package nodes

import (
	"errors"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/node"
	log "github.com/sirupsen/logrus"
	"strings"
)

type NodesList struct {
	Nodes []*node.Schema `json:"nodes"`
}

// Encode the flow from the flow-eng in correct format for react-flow
func Encode(graph *flowctrl.Flow) (*NodesList, error) {
	var listSchema []*node.Schema
	for _, _node := range graph.GetNodesSpec() { // we need to add each node that has one link
		nodeSchema := &node.Schema{}
		nodeType, err := setType(_node)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		nodeSchema = &node.Schema{
			Id:       _node.GetID(),
			Type:     nodeType,
			NodeName: _node.GetNodeName(),
			Metadata: _node.GetMetadata(),
			Settings: _node.GetSettings(),
			IsParent: _node.IsParent,
			ParentId: _node.ParentId,
			Payload:  _node.Payload,
			Icon:     _node.GetIcon(),
		}
		if len(_node.GetInputs()) > 0 {
			links := map[string]node.SchemaInputs{}
			// for a node we need its input and see if it has a link, if so we need the uuid of the node its link to
			for _, input := range _node.GetInputs() {
				inputsLinks := node.SchemaInputs{}
				inputsLinks.Position = input.Position
				inputsLinks.OverridePosition = input.OverridePosition
				// check the input has links
				destOutputName := input.Connection.NodePort
				if destOutputName != "" {
					link := node.SchemaLinks{}
					link.Socket = string(input.Connection.NodePort)
					link.NodeId = input.Connection.NodeID
					inputsLinks.Links = append(inputsLinks.Links, link)
					links[string(input.Name)] = inputsLinks
					nodeSchema.Inputs = links
				} else if input.Connection.OverrideValue != nil {
					inputsLinks.Value = input.Connection.OverrideValue
					links[string(input.Name)] = inputsLinks
					nodeSchema.Inputs = links
				} else {
					links[string(input.Name)] = inputsLinks
					nodeSchema.Inputs = links
				}
			}
			// fmt.Println(links)
			listSchema = append(listSchema, nodeSchema)
		} else { // if a node has no input then add it here
			listSchema = append(listSchema, nodeSchema)
		}
	}
	encodedNodes := &NodesList{
		Nodes: listSchema,
	}

	return encodedNodes, nil
}

func setType(n *node.Spec) (string, error) {
	if n == nil {
		return "", errors.New("node info can not be empty")
	}
	if n.Info.Name == "" {
		return "", errors.New("node name can not be empty")
	}
	if n.Info.Category == "" {
		return "", errors.New("node category can not be empty")
	}
	return fmt.Sprintf("%s/%s", n.Info.Category, n.Info.Name), nil

}

func decodeType(nodeType string) (category, name string, err error) {
	parts := strings.Split(nodeType, "/")
	if len(parts) > 1 {
		return parts[0], parts[1], nil
	}
	return "", "", errors.New("failed to get category and name from node-type")
}
