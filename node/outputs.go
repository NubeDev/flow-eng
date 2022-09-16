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
			}
		}
	}
	port.Connections = connections
	return port
}

// DynamicOutputs build n number of outputs -- out1, out2, out3, ..., outN
//	-p overrideNames[]string -> for example, we can pass in [a,b,c,d] or [less, grater, equal]
func DynamicOutputs(dataType DataTypes, fallback interface{}, count, minAllowed, maxAllowed int, outputs []*Output, overrideNames ...[]string) []*Output {
	var out []*Output
	if count < minAllowed {
		count = minAllowed
	}
	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("%s%d", OutputNamePrefix, i)
		if len(overrideNames) > 0 { // for example, we can pass in [a,b,c,d] or [temp, humidity]
			var n string
			overrideName := overrideNames[0]
			if len(overrideName) >= i {
				n = overrideName[i-1]
			}
			if n == "" { // if count in wrong then use in1, in2 and so on
				n = name
			}
			if i < maxAllowed {
				out = append(out, BuildOutput(OutputName(name), dataType, fallback, outputs))
			}
		} else {
			if i < maxAllowed {
				out = append(out, BuildOutput(OutputName(name), dataType, fallback, outputs))
			}
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
