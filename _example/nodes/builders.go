package nodes

import (
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/node"
)

func BuildInputs(body *Node) []*node.TypeInput {
	var out []*node.TypeInput
	for _, input := range body.InputList {
		out = append(out, buildInputFloat(input.PortCommon.Name, input.PortCommon.Connection))
	}
	return out
}

func BuildOutputs(body *Node) []*node.TypeOutput {
	var out []*node.TypeOutput
	for _, output := range body.OutputList {
		out = append(out, buildOutputFloat(output.PortCommonOut.Name, output.PortCommonOut.Connections))
	}
	return out
}

func BuildInputFloat(portName node.PortName) *node.TypeInput {
	var dataType buffer.Type
	dataType = buffer.Float64
	var port = node.NewInputPort(dataType)
	return &node.TypeInput{
		PortCommon: &node.PortCommon{
			Name:       portName,
			Type:       node.TypeFloat64,
			Connection: nil,
		},
		InputPort:    port,
		ValueFloat64: adapter.NewFloat64(port),
	}
}

func BuildOutputFloat(portName node.PortName) *node.TypeOutput {
	var dataType buffer.Type
	dataType = buffer.Float64
	var port = node.NewOutputPort(dataType)
	return &node.TypeOutput{
		PortCommonOut: &node.PortCommonOut{
			Name:        portName,
			Type:        node.TypeFloat64,
			Connections: nil,
		},
		OutputPort:   port,
		ValueFloat64: adapter.NewFloat64(port),
	}
}

func buildInputFloat(inputName node.PortName, conn *node.Connection) *node.TypeInput {
	var dataType buffer.Type
	dataType = buffer.Float64
	var port = node.NewInputPort(dataType)
	return &node.TypeInput{
		PortCommon: &node.PortCommon{
			Name: inputName,
			Type: node.TypeFloat64,
			Connection: &node.Connection{
				NodeID: conn.NodeID,
			},
		},
		InputPort:    port,
		ValueFloat64: adapter.NewFloat64(port),
	}
}

func buildOutputFloat(outputName node.PortName, conn []*node.Connection) *node.TypeOutput {
	var dataType buffer.Type
	dataType = buffer.Float64
	var port = node.NewOutputPort(dataType)
	return &node.TypeOutput{
		PortCommonOut: &node.PortCommonOut{
			Name:        outputName,
			Type:        node.TypeFloat64,
			Connections: conn,
		},
		OutputPort:   port,
		ValueFloat64: adapter.NewFloat64(port),
	}
}

func GetInput(body *Node, num int) *node.TypeInput {
	for i, input := range body.InputList {
		if i == num {
			return input
		}
	}
	return nil
}

func GetOutConnections(body *Node, num int) []*node.Connection {
	for i, output := range body.OutputList {
		if i == num {
			return output.Connections
		}
	}
	return nil
}
