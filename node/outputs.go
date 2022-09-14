package node

import "fmt"

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

func BuildOutputs(body ...*Output) []*Output {
	var out []*Output
	for _, output := range body {
		out = append(out, output)
	}
	return out
}
