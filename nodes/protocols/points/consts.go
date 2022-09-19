package points

import "errors"

type ObjectID int
type ObjectType string
type IoNumber int  // 1, 2
type IoPort string // UI, UO eg: IoType:IoNumber -> UI1, UI2
type IoType string // digital

const AnalogInput ObjectType = "analogInput"
const AnalogOutput ObjectType = "analogOutput"
const AnalogVariable ObjectType = "analogVariable"

const BinaryInput ObjectType = "binaryInput"
const BinaryOutput ObjectType = "binaryOutput"
const BinaryVariable ObjectType = "binaryVariable"

const IoTypeTemp IoType = "thermistor_10k_type_2"
const IoTypeCurrent IoType = "current"
const IoTypeVolts IoType = "voltage_dc"
const IoTypeDigital IoType = "digital"

func ObjectSwitcher(o ObjectType) (string, error) {
	switch o {
	case AnalogInput:
		return "ai", nil
	case AnalogOutput:
		return "ao", nil
	case AnalogVariable:
		return "av", nil
	case BinaryInput:
		return "bi", nil
	case BinaryOutput:
		return "bo", nil
	case BinaryVariable:
		return "bv", nil
	}
	return "", errors.New("unable to find a valid object to convert the mqtt topic")
}
