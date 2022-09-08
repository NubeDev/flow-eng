package node

import (
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/helpers"
)

func BuildInput(portName PortName, dataType DataTypes, inputs []*Input) *Input {
	out := &Input{}
	port := &InputPort{
		Name:     portName,
		DataType: dataType,
	}
	_dataType := buffer.String
	port = newInputPort(_dataType, port)
	out.InputPort = port
	out.Value = adapter.NewString(port)
	var addConnections bool
	for _, input := range inputs {
		if input.Name == portName {
			addConnections = true
			if input.Connection != nil { // this would be when the flow comes from json
				out.Connection = input.Connection
			} else {
				out.Connection = &InputConnection{}
			}
		}
	}
	if !addConnections {
		out.Connection = &InputConnection{}
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
	_dataType := buffer.String
	port = newOutputPort(_dataType, port)
	out.OutputPort = port
	out.Value = adapter.NewString(port)
	var connections []*Connection
	for _, output := range outputs {
		if output.Name == portName {
			for _, connection := range output.Connections {
				if connection.NodeID != "" && connection.NodePort != "" {
					connections = append(connections, connection)
				}
			}
		}
	}
	if out.Connections == nil {
		out.Connections = connections
	}
	return out
}

func BuildInputs(body ...*Input) []*Input {
	var out []*Input
	for _, input := range body {
		out = append(out, input)
	}
	return out
}

func BuildOutputs(body ...*Output) []*Output {
	var out []*Output
	for _, output := range body {
		out = append(out, output)
	}
	return out
}

func EmptyNode(body *BaseNode) *BaseNode {
	if body == nil {
		body = &BaseNode{
			Info: Info{
				NodeID: "",
			},
		}
	}
	return body
}

func SetUUID(uuid string) string {
	if uuid == "" {
		uuid = helpers.ShortUUID("node")
	}
	return uuid
}
