package node

import "fmt"

func SetOutputHelp(help OutputHelp) *OutputOpts {
	return &OutputOpts{
		Help: help,
	}
}

type OutputOpts struct {
	Help OutputHelp
}

func outputOptions(opts ...*OutputOpts) *OutputOpts {
	if len(opts) > 0 {
		return opts[0]
	}
	return &OutputOpts{}
}

func outputHelp(opts ...*OutputOpts) (help OutputHelp) {
	return outputOptions(opts...).Help
}

func BuildOutput(portName OutputName, dataType DataTypes, fallback interface{}, outputs []*Output, opts ...*OutputOpts) *Output {
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
	port.Help = outputHelp(opts...)
	return port
}

// DynamicOutputs build n number of outputs -- out1, out2, out3, ..., outN
//	-p overrideNames[]string -> for example, we can pass in [a,b,c,d] or [less, grater, equal]
func DynamicOutputs(dataType DataTypes, fallback interface{}, count, minAllowed, maxAllowed int, outputs []*Output, overrideNames ...[]string) []*Output {
	if len(overrideNames) > 0 {
		if len(overrideNames[0]) < count {
			panic("build dynamic-outputs name length must match the count length of the required outputs")
		}
	}
	var out []*Output
	if count < minAllowed {
		count = minAllowed
	}
	if len(overrideNames) > 0 {
		for _, names := range overrideNames {
			for i, name := range names {
				if i < count {
					out = append(out, BuildOutput(OutputName(name), dataType, fallback, outputs))
				}

			}
		}
	} else {
		for i := 1; i <= count; i++ {
			name := fmt.Sprintf("%s%d", OutputNamePrefix, i)
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
