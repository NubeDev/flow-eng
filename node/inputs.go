package node

import "fmt"

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

func BuildInputs(body ...*Input) []*Input {
	var out []*Input
	for _, input := range body {
		out = append(out, input)
	}
	return out
}

// DynamicInputs build n number of inputs -- in1, in2, in3, ..., inN
//	-p overrideNames[]string -> for example, we can pass in [a,b,c,d] or [temp, humidity]
func DynamicInputs(dataType DataTypes, fallback interface{}, count, minAllowed, maxAllowed int, inputs []*Input, overrideNames ...[]string) []*Input {
	var out []*Input
	if count < minAllowed {
		count = minAllowed
	}
	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("%s%d", InputNamePrefix, i)
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
				out = append(out, BuildInput(InputName(n), dataType, fallback, inputs))
			}

		} else {
			if i < maxAllowed {
				out = append(out, BuildInput(InputName(name), dataType, fallback, inputs))
			}
		}

	}
	return out
}