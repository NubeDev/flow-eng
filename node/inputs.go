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

func BuildInput(portName InputName, dataType DataTypes, defaultValue interface{}, inputs []*Input, hasSetting bool, preventOverride bool, opts ...*InputsOpts) *Input {
	var settingName string
	if hasSetting {
		settingName = string(portName)
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
		if input.Name == portName {
			if input.Connection.DefaultValue == nil {
				port.Connection.DefaultValue = defaultValue
			}
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
func DynamicInputs(dataType DataTypes, defaultVal interface{}, count, minAllowed, maxAllowed int, inputs []*Input, preventOverrides bool, overrideNames ...[]string) []*Input {
	if len(overrideNames) > 0 {
		if len(overrideNames[0]) < count {
			panic("build dynamic-inputs name length must match the count length of the required inputs")
		}
	}
	var out []*Input
	if count < minAllowed {
		count = minAllowed
	}
	maxAllowed += 1 // this is needed otherwise the last input will not be added
	if len(overrideNames) > 0 {
		for _, names := range overrideNames {
			for i, name := range names {
				if i < count {
					out = append(out, BuildInput(InputName(name), dataType, defaultVal, inputs, false, preventOverrides))
				}

			}
		}
	} else {
		for i := 1; i <= count; i++ {
			name := fmt.Sprintf("%s%d", InputNamePrefix, i)
			if i < maxAllowed {
				out = append(out, BuildInput(InputName(name), dataType, defaultVal, inputs, false, preventOverrides))
			}
		}
	}
	return out
}
