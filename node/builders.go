package node

import (
	"errors"
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/helpers"
)

func Builder(name string, body *Node) (*Node, error) {
	switch name {
	case nodeA:
		return SpecNodeA(body)
	case nodeB:
		return SpecNodeA(body)
	}
	return nil, errors.New("node not found")
}

func buildInput(portName PortName, dataType DataTypes, inputs []*Input) *Input {
	out := &Input{}
	port := &InputPort{
		Name:     portName,
		DataType: dataType,
	}
	var _dataType buffer.Type
	if dataType == TypeFloat64 {
		_dataType = buffer.Float64
	}
	port = newInputPort(_dataType, port)
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

func buildOutput(portName PortName, dataType DataTypes, outputs []*Output) *Output {
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
	port = newOutputPort(_dataType, port)
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

func buildInputs(body ...*Input) []*Input {
	var out []*Input
	for _, input := range body {
		out = append(out, input)
	}
	return out
}

func buildOutputs(body ...*Output) []*Output {
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
