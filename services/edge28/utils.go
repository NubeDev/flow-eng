package edge28lib

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/structs"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube/edge28"
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/nube/thermistor"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"reflect"
	"strconv"
)

// GetGPIOValueForUOByType converts the point value to the correct edge28 UO GPIO value based on the IoType
func getValueUO(ioType string, value float64) (float64, error) {
	var err error
	var out float64
	switch ioType {
	case UOTypes.VOLTSDC:
		out = edge28.VoltageToGPIOValue(value)
	default:
		err = errors.New("UO IoType is not a recognized type")
	}
	return out, err

}

// getValueUI converts the GPIO value to the scaled UI value based on the IoType
func getValueUI(ioType string, value float64) (float64, error) {
	var err error
	var result float64

	if !structs.ExistsInStrut(UITypes, ioType) {
		err = errors.New(fmt.Sprintf("skipping IoType %s not recognized.", ioType))
		return 0, err
	}
	switch ioType {
	case UITypes.RAW:
		result = value
	case UITypes.DIGITAL:
		result = edge28.GPIOValueToDigital(value)
	case UITypes.VOLTSDC:
		result = edge28.GPIOValueToVoltage(value)
	case UITypes.MILLIAMPS:
		result = edge28.ScaleGPIOValueTo420ma(value)
	case UITypes.RESISTANCE:
		result = edge28.ScaleGPIOValueToResistance(value)
	case UITypes.THERMISTOR10KT2:
		resistance := edge28.ScaleGPIOValueToResistance(value)
		result, err = thermistor.ResistanceToTemperature(resistance, thermistor.T210K)
	default:
		err = errors.New("UI IoType is not a recognized type")
		return 0, err
	}
	return result, nil
}

// convertDigital converts true/false values (all basic types allowed) to BBB GPIO 0/1 ON/OFF (FOR DOs/Relays) and to 100/0 (FOR UOs).  Note that the GPIO value for digital points is inverted.
func convertDigital(input interface{}, isUO bool) (float64, error) {
	var inputAsBool bool
	var err error = nil
	switch input.(type) {
	case string:
		inputAsBool, err = strconv.ParseBool(reflect.ValueOf(input).String())
	case int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
		inputAsBool = reflect.ValueOf(input).Int() != 0
	case float32, float64:
		inputAsBool = reflect.ValueOf(input).Float() != float64(0)
	case bool:
		inputAsBool = reflect.ValueOf(input).Bool()
	default:
		err = errors.New("edge28-polling: input is not a recognized type")
	}
	if err != nil {
		return 0, err
	} else if inputAsBool {
		if isUO {
			return 0, nil // 0 is the 12vdc/ON GPIO value for UOs
		} else {
			return 1, nil // 1 is the 12vdc/ON GPIO value for DOs/Relays
		}
	} else {
		if isUO {
			return 100, nil // 100 is the 0vdc/OF GPIO value for UOs
		} else {
			return 0, nil // 0 is the 0vdc/OFF GPIO value for DOs/Relays
		}
	}
}

func limitValue(ioType string, inputVal float64) (outputVal float64) {
	inputValFloat := inputVal
	switch ioType {
	case UOTypes.DIGITAL, UITypes.DIGITAL:
		if inputValFloat <= 0 {
			outputVal = 0
		} else {
			outputVal = 1
		}
	case UOTypes.VOLTSDC, UITypes.VOLTSDC:
		if inputValFloat <= 0 {
			outputVal = 0
		} else if inputValFloat >= 10 {
			outputVal = 10
		} else {
			outputVal = inputVal
		}
	default:
		outputVal = inputVal
	}
	return outputVal
}

func limitPriorityArrayByEdge28Type(ioType string, priority *model.PointWriter) *map[string]*float64 {
	priorityMap := map[string]*float64{}
	for key, val := range *priority.Priority {
		var outputVal *float64
		if val == nil {
			outputVal = nil
		} else {
			switch ioType {
			case UOTypes.DIGITAL, UITypes.DIGITAL:
				if *val <= 0 {
					outputVal = float.New(0)
				} else {
					outputVal = float.New(1)
				}

			case UOTypes.VOLTSDC, UITypes.VOLTSDC:
				if *val <= 0 {
					outputVal = float.New(0)
				} else if *val >= 10 {
					outputVal = float.New(10)
				} else {
					outputVal = float.New(*val)
				}
			default:
				outputVal = float.New(*val)
			}
		}
		priorityMap[key] = outputVal
	}
	return &priorityMap
}
