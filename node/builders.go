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

func BuildInput(portName PortName, dataType DataTypes, fallback interface{}, inputs []*InputPort) *InputPort {
	port := &InputPort{
		Name:       portName,
		DataType:   dataType,
		Connection: &InputConnection{}}
	port = newInputPort(port)
	var addConnections bool
	if len(inputs) == 0 {
		inputs = []*InputPort{port}
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

func BuildOutput(portName PortName, dataType DataTypes, fallback interface{}, outputs []*OutputPort) *OutputPort {
	var connections []*OutputConnection
	port := &OutputPort{
		Name:        portName,
		DataType:    dataType,
		Connections: connections,
	}
	port = newOutputPort(port)
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

// DynamicInputs build n number of inputs
// startOfName eg: in would make in1, in2, in3
func DynamicInputs(startOfName PortName, dataType DataTypes, fallback interface{}, count, minAllowed, maxAllowed int, inputs []*InputPort) []*InputPort {
	var out []*InputPort
	if count < minAllowed {
		count = minAllowed
	}
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("%s%d", startOfName, i+1)
		if i < maxAllowed {
			out = append(out, BuildInput(PortName(name), dataType, fallback, inputs))
		}
	}
	return out
}

// DynamicOutputs build n number of outputs
// startOfName eg: in would make out1, out2, and so on
func DynamicOutputs(startOfName PortName, dataType DataTypes, fallback interface{}, count, maxAllowed int, outputs []*OutputPort) []*OutputPort {
	var out []*OutputPort
	for i := 0; i < count; i++ {
		name := fmt.Sprintf("%s%d", startOfName, i+1)
		if i < maxAllowed {
			out = append(out, BuildOutput(PortName(name), dataType, fallback, outputs))
		}
	}
	return out
}

func BuildNode(body *BaseNode, inputs []*InputPort, outputs []*OutputPort, settings []*Settings) *BaseNode {
	body.Settings = settings
	body.Inputs = inputs
	body.Outputs = outputs
	return body
}

func BuildInputs(body ...*InputPort) []*InputPort {
	var out []*InputPort
	for _, input := range body {
		out = append(out, input)
	}
	return out
}

func BuildOutputs(body ...*OutputPort) []*OutputPort {
	var out []*OutputPort
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
