package node

import (
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/math"
	log "github.com/sirupsen/logrus"
	"testing"
)

type NodesList struct {
	Nodes interface{} `json:"nodes"`
}

func TestBaseNode_NodeConnection(t *testing.T) {

	const1, err := math.NewConst(nil) // new node
	if err != nil {
		log.Errorln(err)
		return
	}

	err = const1.OverrideInputValue(node.In1, nil)
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

	graph := flowctrl.New() // init the flow engine

	graph.AddNode(const1) // add the nodes to the runtime
	graph.AddNode(const2)
	graph.AddNode(add)
	// graph.AddNode(mqttSub)
	// graph.AddNode(mqttPub)

	err = graph.ManualNodeConnector(const1, add, node.Out1, node.In1) // connect const-1 and 2 to the add node
	if err != nil {
		log.Errorln(err)
		return
	}

	err = graph.ManualNodeConnector(const2, add, node.Out1, node.In2)
	if err != nil {
		log.Errorln(err)
		return
	}

	graph.ReBuildFlow(true)

	var listSchema []*node.Schema
	nodeSchema := &node.Schema{}

	for _, baseNode := range graph.GetNodesBase() { // we need to add each node that has one connection
		nodeSchema = &node.Schema{
			Id:   baseNode.GetID(),
			Type: setType(baseNode),
		}
		//var link map[string]map[string][]node.Links
		if len(baseNode.GetInputs()) > 0 {
			links := map[string]node.Links{}

			// for a node we need its input and see if it has a connection, if so we need the uuid of the node its connection to
			for _, input := range baseNode.GetInputs() {
				// check the input has connection
				destOutputName := input.Connection.NodePort
				if destOutputName != "" {
					inputName := input.Name
					destNodeId := input.Connection.NodeID
					sourceNode := graph.GetNode(destNodeId)
					for _, output := range sourceNode.GetOutputs() {
						if output.Name == destOutputName {
							links[string(inputName)] = node.Links{
								NodeId: destNodeId,
								Socket: string(destOutputName),
							}
						}
					}

				}
			}
			nodeSchema.Metadata = &node.Metadata{
				PositionX: "271.5",
				PositionY: "-69",
			}
			nodeSchema.Inputs = map[string]map[string]node.Links{"links": links}
			listSchema = append(listSchema, nodeSchema)
		} else { // if a node has no input then add it here
			nodeSchema.Metadata = &node.Metadata{
				PositionX: "271.5",
				PositionY: "-69",
			}
			listSchema = append(listSchema, nodeSchema)
		}

	}

	a := NodesList{
		Nodes: listSchema,
	}

	pprint.PrintJOSN(a)

}

func setType(n *node.BaseNode) string {
	if n == nil {
		return ""
	}
	return fmt.Sprintf("%s/%s", n.Info.Category, n.Info.Name)

}

func TestBaseNode_NodeConnection3(t *testing.T) {

	const1, err := math.NewConst(nil) // new node
	if err != nil {
		log.Errorln(err)
		return
	}

	err = const1.OverrideInputValue(node.In1, nil)
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

	graph := flowctrl.New() // init the flow engine

	graph.AddNode(const1) // add the nodes to the runtime
	graph.AddNode(const2)
	graph.AddNode(add)
	// graph.AddNode(mqttSub)
	// graph.AddNode(mqttPub)

	err = graph.ManualNodeConnector(const1, add, node.Out1, node.In1) // connect const-1 and 2 to the add node
	if err != nil {
		log.Errorln(err)
		return
	}

	err = graph.ManualNodeConnector(const2, add, node.Out1, node.In2)
	if err != nil {
		log.Errorln(err)
		return
	}

	graph.ReBuildFlow(true)

	var listSchema []*node.Schema
	nodeSchema := &node.Schema{}

	for _, baseNode := range graph.GetNodesBase() { // we need to add each node that has one connection
		nodeSchema = &node.Schema{
			Id:   baseNode.GetID(),
			Type: setType(baseNode),
		}
		//var link map[string]map[string][]node.Links
		if len(baseNode.GetInputs()) > 0 {
			var links map[string]node.Links
			// for a node we need its input and see if it has a connection, if so we need the uuid of the node its connection to
			for _, input := range baseNode.GetInputs() {
				// check the input has connection
				destOutputName := input.Connection.NodePort
				if destOutputName != "" {
					inputName := input.Name
					destNodeId := input.Connection.NodeID
					sourceNode := graph.GetNode(destNodeId)
					for _, output := range sourceNode.GetOutputs() {
						if output.Name == destOutputName {
							link := map[string]map[string][]node.Links{string(inputName): map[string][]node.Links{"links": []node.Links{
								node.Links{
									NodeId: destNodeId,
									Socket: string(destOutputName),
								},
							}}}
							fmt.Println(link)
						}

						//aa := links[string(inputName)].NodeId

					}

				}
			}
			nodeSchema.Metadata = &node.Metadata{
				PositionX: "271.5",
				PositionY: "-69",
			}
			nodeSchema.Inputs = links
			listSchema = append(listSchema, nodeSchema)
		} else { // if a node has no input then add it here
			nodeSchema.Metadata = &node.Metadata{
				PositionX: "271.5",
				PositionY: "-69",
			}
			listSchema = append(listSchema, nodeSchema)
		}

	}

	a := NodesList{
		Nodes: listSchema,
	}

	pprint.PrintJOSN(a)

}

func TestBaseNode_NodeConnection2(t *testing.T) {

	const1, err := math.NewConst(nil) // new node
	if err != nil {
		log.Errorln(err)
		return
	}

	err = const1.OverrideInputValue(node.In1, 11.0)
	if err != nil {
		log.Errorln(err)
		return
	}
	const2, err := math.NewConst(nil)
	if err != nil {
		log.Errorln(err)
		return
	}
	err = const2.OverrideInputValue(node.In1, 11.0)
	if err != nil {
		log.Errorln(err)
		return
	}
	add, err := math.NewAdd(nil) // new math (add) node
	if err != nil {
		log.Errorln(err)
		return
	}

	graph := flowctrl.New() // init the flow engine

	graph.AddNode(const1) // add the nodes to the runtime
	graph.AddNode(const2)
	graph.AddNode(add)
	// graph.AddNode(mqttSub)
	// graph.AddNode(mqttPub)

	err = graph.ManualNodeConnector(const1, add, node.Out1, node.In1) // connect const-1 and 2 to the add node
	if err != nil {
		log.Errorln(err)
		return
	}

	err = graph.ManualNodeConnector(const2, add, node.Out1, node.In2)
	if err != nil {
		log.Errorln(err)
		return
	}

	graph.ReBuildFlow(true)

	var listSchema []*node.Schema
	nodeSchema := &node.Schema{}

	for _, baseNode := range graph.GetNodesBase() {
		nodeSchema = &node.Schema{
			Id:   baseNode.GetID(),
			Type: setType(baseNode),
		}

		// for a node we need its input and see if it has a connection, if so we need the uuid of the node its connection to
		for _, input := range baseNode.GetInputs() {
			inputName := input.Name
			destOutputName := input.Connection.NodePort
			destNodeId := input.Connection.NodeID
			if destNodeId != "" {
				sourceNode := graph.GetNode(destNodeId)
				for _, output := range sourceNode.GetOutputs() {
					if output.Name == destOutputName {
						var (
							value = map[string]map[string][]node.Links{string(inputName): map[string][]node.Links{"links": []node.Links{
								node.Links{
									NodeId: destNodeId,
									Socket: string(destOutputName),
								},
							}}}
						)

						nodeSchema.Metadata = &node.Metadata{
							PositionX: "271.5",
							PositionY: "-69",
						}
						nodeSchema.Inputs = value
						listSchema = append(listSchema, nodeSchema)
					}
				}
			}

		}
	}

	a := NodesList{
		Nodes: listSchema,
	}

	pprint.PrintJOSN(a)

}

func TestBaseNode_NodeNonConnection(t *testing.T) {
	var list []*node.Schema
	var value = map[string]map[string]string{"duration": map[string]string{"value": "22"}}
	s1 := &node.Schema{
		Id:   "2",
		Type: "time/delay",
		Metadata: &node.Metadata{
			PositionX: "271.5",
			PositionY: "-69",
		},
		Inputs: value,
	}

	list = append(list, s1)
	a := NodesList{
		Nodes: list,
	}

	pprint.PrintJOSN(a)

}
