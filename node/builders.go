package node

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
)

func BuildNodes(body ...Node) []Node {
	var out []Node
	for _, output := range body {
		out = append(out, output)
	}
	return out
}

func BuildInput(portName InputName, dataType DataTypes, fallback interface{}, inputs []*Input) *Input {
	port := &Input{
		Name:       portName,
		DataType:   dataType,
		Connection: &InputConnection{},
	}
	port = newInput(port)
	var addConnections bool
	if len(inputs) == 0 {
		inputs = []*Input{port}
	}
	for _, input := range inputs {
		if input.Connection.FallbackValue == nil {
			port.Connection.FallbackValue = fallback
		}
		if input.Name == portName {
			addConnections = true
			if input.Connection != nil { // this would be when the flow comes from json
				port.Connection = input.Connection
			} else {
				port.Connection = &InputConnection{}
			}
		}
	}
	if !addConnections {
		port.Connection = &InputConnection{}
	}
	return port
}

func BuildOutput(portName OutputName, dataType DataTypes, fallback interface{}, outputs []*Output) *Output {
	var connections []*OutputConnection
	port := &Output{
		Name:        portName,
		DataType:    dataType,
		Connections: connections,
	}
	port = newOutput(port)
	for _, output := range outputs {
		if output.Name == portName {
			for _, connection := range output.Connections {
				if connection.FallbackValue == nil {
					connection.FallbackValue = fallback
				}
				if connection.NodeID != "" && connection.NodePort != "" {
					connections = append(connections, connection)
				}
			}
		}
	}
	port.Connections = connections
	return port
}

// DynamicInputs build n number of inputs -- in1, in2, in3, ..., inN
func DynamicInputs(dataType DataTypes, fallback interface{}, n, minAllowed, maxAllowed int, inputs []*Input) []*Input {
	var out []*Input
	if n < minAllowed {
		n = minAllowed
	}
	for i := 1; i <= n; i++ {
		name := fmt.Sprintf("%s%d", InputNamePrefix, i)
		if i < maxAllowed {
			out = append(out, BuildInput(InputName(name), dataType, fallback, inputs))
		}
	}
	return out
}

// DynamicOutputs build n number of outputs -- out1, out2, out3, ..., outN
func DynamicOutputs(dataType DataTypes, fallback interface{}, n, maxAllowed int, outputs []*Output) []*Output {
	var out []*Output
	for i := 1; i <= n; i++ {
		name := fmt.Sprintf("%s%d", OutputNamePrefix, i+1)
		if i < maxAllowed {
			out = append(out, BuildOutput(OutputName(name), dataType, fallback, outputs))
		}
	}
	return out
}

func BuildNode(body *BaseNode, inputs []*Input, outputs []*Output, settings []*Settings) *BaseNode {
	body.Settings = settings
	body.Inputs = inputs
	body.Outputs = outputs
	return body
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

func Defaults(body *BaseNode, nodeName, category string) *BaseNode {
	if body == nil {
		body = &BaseNode{
			Info: Info{
				NodeName: helpers.ShortUUID(nodeName),
				NodeID:   "",
			},
		}
	}
	body.Info.Name = SetName(nodeName)
	body.Info.Category = SetName(category)
	body.Info.NodeID = SetUUID(body.Info.NodeID)
	return body
}

func SetUUID(uuid string) string {
	if uuid == "" {
		uuid = helpers.ShortUUID("node")
	}
	return uuid
}
