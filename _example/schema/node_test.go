package node

import (
	"errors"
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/helpers/float"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes"
	"github.com/NubeDev/flow-eng/nodes/math"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"strings"
	"testing"
)

type NodesList struct {
	Nodes []*node.Schema `json:"nodes"`
}

var encodedNodes NodesList

func TestSpec_NodeConnectionEncode(t *testing.T) {

	const1, err := math.NewConst(nil) // new node
	if err != nil {
		log.Errorln(err)
		return
	}

	err = const1.OverrideInputValue(node.In1, 22)
	if err != nil {
		log.Errorln(err)
		return
	}
	const2, err := math.NewConst(nil)
	if err != nil {
		log.Errorln(err)
		return
	}
	err = const2.OverrideInputValue(node.In1, nil)
	if err != nil {
		log.Errorln(err)
		return
	}
	add, err := math.NewAdd(nil) // new math (add) node
	if err != nil {
		log.Errorln(err)
		return
	}

	add2, err := math.NewAdd(nil) // new math (add) node
	if err != nil {
		log.Errorln(err)
		return
	}

	err = add2.OverrideInputValue(node.In1, 22)
	if err != nil {
		log.Errorln(err)
		return
	}
	err = add2.OverrideInputValue(node.In2, 44)
	if err != nil {
		log.Errorln(err)
		return
	}

	graph := flowctrl.New() // init the flow engine

	graph.AddNode(const1) // add the nodes to the runtime
	graph.AddNode(const2)
	graph.AddNode(add)
	graph.AddNode(add2)
	//err = graph.ManualNodeConnector(const1, add, node.Out1, node.In1) // connect const-1 and 2 to the add node
	//if err != nil {
	//	log.Errorln(err)
	//	return
	//}
	//
	//err = graph.ManualNodeConnector(const2, add, node.Out1, node.In2)
	//if err != nil {
	//	log.Errorln(err)
	//	return
	//}

	graph.ReBuildFlow(true)

	var listSchema []*node.Schema
	nodeSchema := &node.Schema{}

	for _, Spec := range graph.GetNodesSpec() { // we need to add each node that has one connection
		nodeType, err := setType(Spec)
		if err != nil {
			fmt.Println(err)
			return
		}
		nodeSchema = &node.Schema{
			Id:   Spec.GetID(),
			Type: nodeType,
		}
		if len(Spec.GetInputs()) > 0 {
			links := map[string]node.SchemaLinks{}
			// for a node we need its input and see if it has a connection, if so we need the uuid of the node its connection to
			for _, input := range Spec.GetInputs() {
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

			//nodeSchema.Inputs = linkValue                 // as a value when no input is connected
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
	encodedNodes = NodesList{
		Nodes: listSchema,
	}

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

func TestSpec_Decode(t *testing.T) {
	var decodedNodes []*node.Spec
	var decodedNode *node.Spec
	for _, encodedNode := range encodedNodes.Nodes {
		ins := &node.SchemaInputs{}
		err := mapstructure.Decode(encodedNode.Inputs, ins)
		if err != nil {
			fmt.Println(err)
			//return
		}
		_, getName, _ := decodeType(encodedNode)
		id := encodedNode.Id
		name := getName
		decodedNode = node.New(id, name, "", encodedNode.Metadata) // create a blank node
		newNode, err := nodes.Builder(decodedNode)                 // make the new node as per its type
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

}

//func TestSpec_NodeNonConnection(t *testing.T) {
//	var list []*node.Schema
//	var value = map[string]map[string]string{"duration": map[string]string{"value": "22"}}
//	s1 := &node.Schema{
//		Id:   "2",
//		Type: "time/delay",
//		Metadata: &node.Metadata{
//			PositionX: "271.5",
//			PositionY: "-69",
//		},
//		Inputs: value,
//	}
//
//	list = append(list, s1)
//	a := NodesList{
//		Nodes: list,
//	}
//
//	pprint.PrintJOSN(a)
//
//}
