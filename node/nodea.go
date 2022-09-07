package node

import (
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/helpers"
)

type NodeA struct {
	*Node
}

func BuildInput(portName PortName, dataType DataTypes, inputs []*Input) *Input {
	out := &Input{}
	port := &InputPort{
		Name:     portName,
		DataType: dataType,
	}
	var _dataType buffer.Type
	if dataType == TypeFloat64 {
		_dataType = buffer.Float64
	}
	port = NewInputPort(_dataType, port)
	out.InputPort = port
	out.ValueFloat64 = adapter.NewFloat64(port)
	var addConnections bool
	for _, input := range inputs {
		if input.Name == portName {
			addConnections = true
			if input.Connection != nil { // this would be when the flow comes from json
				out.Connection = input.Connection
			} else {
				out.Connection = &Connection{}
			}
		}
	}
	if !addConnections {
		out.Connection = &Connection{}
	}
	return out
}

func BuildOutput(portName PortName, dataType DataTypes, outputs []*Output) *Output {
	out := &Output{}
	port := &OutputPort{
		Name:        portName,
		DataType:    dataType,
		Connections: nil,
	}
	var _dataType buffer.Type
	if dataType == TypeFloat64 {
		_dataType = buffer.Float64
	}
	port = NewOutputPort(_dataType, port)
	out.OutputPort = port
	out.ValueFloat64 = adapter.NewFloat64(port)

	for _, output := range outputs {
		if output.Name == portName {
			for _, connection := range output.Connections {
				out.Connections = []*Connection{connection}
			}
		}
	}
	if out.Connections == nil {
		out.Connections = []*Connection{&Connection{}}
	}
	return out
}

func BuildOutputConnection(body *Output) *Output {

	body.Connections = []*Connection{&Connection{}}
	return body
}

func BuildInputs2(body ...*Input) []*Input {
	var out []*Input
	for _, input := range body {
		out = append(out, input)
	}
	return out
}

func BuildOutputs2(body ...*Output) []*Output {
	var out []*Output
	for _, output := range body {
		out = append(out, output)
	}
	return out
}

func emptyNode(body *Node) *Node {
	if body == nil {
		body = &Node{
			Info: Info{
				NodeID: "",
			},
		}
	}
	return body
}

func setUUID(uuid string) string {
	if uuid == "" {
		uuid = helpers.ShortUUID("node")
	}
	return uuid
}

func SpecNodeA(body *Node) *Node {
	body = emptyNode(body)
	body.Info.Name = "nodeA"
	body.Info.NodeID = setUUID(body.Info.NodeID)
	body.Inputs = BuildInputs2(BuildInput("in1", TypeFloat64, body.Inputs), BuildInput("in2", TypeFloat64, body.Inputs))
	body.Outputs = BuildOutputs2(BuildOutput("out1", TypeFloat64, body.Outputs))
	// then update the connections
	return body
}

func (n *NodeA) Process() {

	for _, out := range n.GetOutputs() {
		out.ValueFloat64.Set(22)
	}
}

func (n *NodeA) Cleanup() {}
