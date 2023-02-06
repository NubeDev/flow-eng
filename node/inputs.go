package node

import (
	"fmt"
)

type InputsOpts struct {
	Help InputHelp
}

func SetInputHelp(help InputHelp) *InputsOpts {
	return &InputsOpts{
		Help: help,
	}
}

func inputsOptions(opts ...*InputsOpts) *InputsOpts {
	if len(opts) > 0 {
		return opts[0]
	}
	return &InputsOpts{}
}

func inputHelp(opts ...*InputsOpts) (help InputHelp) {
	return inputsOptions(opts...).Help
}

func BuildInput(portName InputName, dataType DataTypes, fallback interface{}, inputs []*Input, settingName *string, opts ...*InputsOpts) *Input {
	if settingName != nil {
		portName = InputName(fmt.Sprintf("[%s]", portName))
	}
	port := &Input{
		Name:        portName,
		DataType:    dataType,
		Connection:  &InputConnection{},
		SettingName: settingName,
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
			port.SettingName = settingName
		}
	}
	if !addConnections {
		port.Connection = &InputConnection{}
	}
	port.Help = inputHelp(opts...)
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
	if len(overrideNames) > 0 {
		if len(overrideNames[0]) < count {
			panic("build dynamic-inputs name length must match the count length of the required inputs")
		}
	}
	var out []*Input
	if count < minAllowed {
		count = minAllowed
	}
	if len(overrideNames) > 0 {
		for _, names := range overrideNames {
			for i, name := range names {
				if i < count {
					out = append(out, BuildInput(InputName(name), dataType, fallback, inputs, nil))
				}

			}
		}
	} else {
		for i := 1; i <= count; i++ {
			name := fmt.Sprintf("%s%d", InputNamePrefix, i)
			if i < maxAllowed {
				out = append(out, BuildInput(InputName(name), dataType, fallback, inputs, nil))
			}
		}
	}
	return out
}
