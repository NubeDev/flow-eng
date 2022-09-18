package main

import (
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes"
	bacnet "github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	log "github.com/sirupsen/logrus"
)

func main() {

	readBody := &node.Spec{
		Info: node.Info{
			NodeID: "read111111",
		},
		SubFlow: &node.SubFlow{
			ParentID: "a123",
		},
	}

	read, err := bacnet.NewBacnetBVRead(readBody, nil)
	if err != nil {
		return
	}

	read2Body := &node.Spec{
		Info: node.Info{
			NodeID: "read22222",
		},
		SubFlow: &node.SubFlow{
			ParentID: "a123",
		},
	}

	read2, err := bacnet.NewBacnetBVRead(read2Body, nil)
	if err != nil {
		return
	}
	//pprint.PrintJOSN(read2)

	bacBody := &node.Spec{
		Info: node.Info{
			NodeID: "a123",
		},
		SubFlow: &node.SubFlow{
			ParentID: "",
			Nodes:    nil,
		},
	}

	bac, err := bacnet.NewServer(bacBody, node.ConvertToSpec(read), node.ConvertToSpec(read2))
	if err != nil {
		log.Errorln(err)
		return
	}
	graph := flowctrl.New() // init the flow engine
	graph.AddNodes(bac)     // bacnetServer node must be added first to inti the store

	for _, spec := range bac.GetSubFlowNodes() { // add the server sub-flow nodes
		node_, err := nodes.Builder(spec, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		graph.AddNode(node_)
	}

	nodesList := graph.GetNodes()

	pprint.PrintJOSN(nodesList)

}
