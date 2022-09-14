package nodes

import (
	"errors"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/node"
	"strings"
)

type NodesList struct {
	Nodes []*node.Schema `json:"nodes"`
}

func Encode(graph *flowctrl.Flow) (*NodesList, error) {
	var listSchema []*node.Schema
	nodeSchema := &node.Schema{}
	for _, _node := range graph.GetNodesSpec() { // we need to add each node that has one connection
		nodeType, err := setType(_node)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		nodeSchema = &node.Schema{
			Id:   _node.GetID(),
			Type: nodeType,
		}
		if len(_node.GetInputs()) > 0 {
			links := map[string]node.SchemaLinks{}
			// for a node we need its input and see if it has a connection, if so we need the uuid of the node its connection to
			for _, input := range _node.GetInputs() {
				// check the input has connection
				destOutputName := input.Connection.NodePort
				if destOutputName != "" {
					inputName := input.Name
					destNodeId := input.Connection.NodeID
					sourceNode := graph.GetNode(destNodeId)
					for _, output := range sourceNode.GetOutputs() {
						if output.Name == destOutputName {
							links[string(inputName)] = node.SchemaLinks{
								NodeId: destNodeId,
								Socket: destOutputName,
							}
						}
					}
				} else {
					if input.Connection.OverrideValue != nil {
						str := fmt.Sprintf("%v", input.Connection.OverrideValue)
						//linkValue = map[string]map[string]string{string(inputName): {"value": str}}
						links[string(input.Name)] = node.SchemaLinks{
							Value: str,
						}
					}
				}
			}
			nodeSchema.Metadata = &node.Metadata{
				PositionX: "271.5",
				PositionY: "-69",
			}
			nodeSchema.Inputs = node.SchemaInputs{Links: links} // when a connection is made
			listSchema = append(listSchema, nodeSchema)
		} else { // if a node has no input then add it here
			nodeSchema.Metadata = &node.Metadata{
				PositionX: "271.5",
				PositionY: "-69",
			}
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

func decodeType(n *node.Schema) (category, name string, err error) {
	if n == nil {
		return "", "", errors.New("node schema can not be empty")
	}
	if n.Type == "" {
		return "", "", errors.New("node type can not be empty")
	}
	parts := strings.Split(n.Type, "/")
	if len(parts) > 1 {
		return parts[0], parts[1], nil
	}
	return "", "", errors.New("failed to get category and name from node-type")

}
