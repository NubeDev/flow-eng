package node

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/helpers"
)

const (
	nodeType    = "pass"
	inputCount  = 1
	outputCount = 1
)

func New(body *Node) (*Node, error) {
	body, err := Check(body, Spec{nodeType, inputCount, outputCount})
	if err != nil {
		return nil, err
	}
	return &Node{
		Inputs:  BuildInputs(body),
		Outputs: BuildOutputs(body),
		Info: Info{
			NodeID:      body.Info.NodeID,
			Name:        body.Info.Name,
			Description: "desc",
			Version:     "1",
		},
	}, nil
}

type Spec struct {
	Name        string
	InputCount  int
	OutputCount int
}

func Check(body *Node, nodeSpec Spec) (*Node, error) {
	if body == nil {
		return nil, errors.New("node body can not be empty")
	}
	if body.Info.Name == "" {
		return nil, errors.New("node name can not be empty, try AND, OR")
	}
	if body.Info.NodeID == "" {
		body.Info.NodeID = helpers.ShortUUID(nodeSpec.Name)
	}
	if len(body.Inputs) != nodeSpec.InputCount {
		return nil, errors.New(fmt.Sprintf("input count is incorrect required:%d provided:%d", nodeSpec.InputCount, len(body.Inputs)))
	}
	if len(body.Outputs) != nodeSpec.OutputCount {
		return nil, errors.New(fmt.Sprintf("output count is incorrect required:%d provided:%d", nodeSpec.OutputCount, len(body.Outputs)))
	}
	return body, nil
}

func BuildInputs(body *Node) []*Input {
	var out []*Input
	for _, input := range body.Inputs {
		out = append(out, BuildInputFloat(input.Name, input.Connection))
	}
	return out
}

func BuildOutputs(body *Node) []*Output {
	var out []*Output
	for _, output := range body.Outputs {
		out = append(out, BuildOutputFloat(output.Name, output.Connections))
	}
	return out
}

func BuildInputFloat(portName PortName, conn *Connection) *Input {
	out := &Input{}
	port := &InputPort{
		Name:       portName,
		DataType:   "",
		Connection: conn,
		Const:      nil,
	}
	var dataType buffer.Type
	dataType = buffer.Float64
	port = NewInputPort(dataType, port)
	out.InputPort = port
	out.ValueFloat64 = adapter.NewFloat64(port)
	return out
}

func BuildOutputFloat(portName PortName, conn []*Connection) *Output {
	out := &Output{}
	port := &OutputPort{
		Name:        portName,
		DataType:    "",
		Connections: conn,
	}
	var dataType buffer.Type
	dataType = buffer.Float64
	port = NewOutputPort(dataType, port)
	out.OutputPort = port
	out.ValueFloat64 = adapter.NewFloat64(port)
	return out
}

func GetInput(body *Node, num int) *Input {
	for i, input := range body.Inputs {
		if i == num {
			return input
		}
	}
	return nil
}

func GetOutConnections(body *Node, num int) []*Connection {
	for i, output := range body.Outputs {
		if i == num {
			return output.Connections
		}
	}
	return nil
}
