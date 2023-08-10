package nodes

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/pallet"
	log "github.com/sirupsen/logrus"
)

type NodesList struct {
	Nodes []*node.Schema `json:"nodes"`
}

// Encode the flow from the flow-eng in correct format for react-flow
func Encode(specs []*node.Spec) (*NodesList, error) {
	var listSchema []*node.Schema
	for _, _node := range specs { // we need to add each node that has one link
		nodeSchema := &node.Schema{}
		nodeType, err := pallet.SetType(_node)
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
				inputsLinks.DefaultValue = input.Connection.DefaultValue
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
