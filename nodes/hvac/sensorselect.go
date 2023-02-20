package hvac

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/conversions"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
)

type SensorSelect struct {
	*node.Spec
}

func NewSensorSelect(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, sensorSelect, category)
	settings := &SensorSelectSettings{}
	err := mapstructure.Decode(body.Settings, &settings)
	if err != nil {
		return nil, err
	}

	if settings == nil {
		body.Settings = map[string]interface{}{}
		body.Settings["sensor-count"] = 2
	} else if settings.SensorCount < 2 {
		body.Settings["sensor-count"] = 2
	}
	inputsCount := int(conversions.GetFloat(body.Settings["inputCount"]))

	fallback := node.BuildInput(node.Fallback, node.TypeFloat, nil, body.Inputs, false, false)
	dynamicInputs := dynamicInputs(inputsCount, 2, 20, body.Inputs)
	inputsArray := []*node.Input{fallback}
	inputsArray = append(inputsArray, dynamicInputs...)
	inputs := node.BuildInputs(inputsArray...)

	validOutput := node.BuildOutput(node.ValidOutput, node.TypeBool, nil, body.Outputs)
	min := node.BuildOutput(node.MinOutput, node.TypeFloat, nil, body.Outputs)
	max := node.BuildOutput(node.MaxOutput, node.TypeFloat, nil, body.Outputs)
	avg := node.BuildOutput(node.AvgOutput, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(validOutput, min, max, avg)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetDynamicInputs()
	n := &SensorSelect{body}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *SensorSelect) Process() {
	settings, _ := inst.getSettings(inst.GetSettings())
	count := settings.SensorCount
	sensorArray := inst.GetSensorArray(count)
	excludeArray := inst.GetExcludeArray(count)
	var validSensorArray []float64
	for i := 0; i < len(sensorArray); i++ {
		if !excludeArray[i] {
			sensorValue := sensorArray[i]
			if sensorValue != nil {
				validSensorArray = append(validSensorArray, *sensorValue)
			}
		}
	}
	if len(validSensorArray) <= 0 {
		fallback, fallbackNull := inst.ReadPinAsFloat(node.Fallback)
		if fallbackNull {
			inst.WritePinBool(node.ValidOutput, false)
			inst.WritePinNull(node.MinOutput)
			inst.WritePinNull(node.MaxOutput)
			inst.WritePinNull(node.AvgOutput)
		} else {
			inst.WritePinBool(node.ValidOutput, true)
			inst.WritePinFloat(node.MinOutput, fallback)
			inst.WritePinFloat(node.MaxOutput, fallback)
			inst.WritePinFloat(node.AvgOutput, fallback)
		}
		return
	}
	sum := 0.0
	var min, max float64
	for j := 0; j < len(validSensorArray); j++ {
		sum += validSensorArray[j]
		if j == 0 {
			min = validSensorArray[j]
			max = validSensorArray[j]
		} else {
			if min > validSensorArray[j] {
				min = validSensorArray[j]
			}
			if max < validSensorArray[j] {
				max = validSensorArray[j]
			}
		}
	}
	avg := sum / float64(len(validSensorArray))
	inst.WritePinBool(node.ValidOutput, true)
	inst.WritePinFloat(node.MinOutput, min)
	inst.WritePinFloat(node.MaxOutput, max)
	inst.WritePinFloat(node.AvgOutput, avg)
}

func dynamicInputs(count, minAllowed, maxAllowed int, inputs []*node.Input) []*node.Input {
	var out []*node.Input
	if count < minAllowed {
		count = minAllowed
	}

	for i := 1; i <= count; i++ {
		if i >= maxAllowed {
			break
		}
		sensorInputName := fmt.Sprintf("sensor%d", i)
		out = append(out, node.BuildInput(node.InputName(sensorInputName), node.TypeFloat, nil, inputs, false, false))
		excludeInputName := fmt.Sprintf("exclude%d", i)
		out = append(out, node.BuildInput(node.InputName(excludeInputName), node.TypeBool, nil, inputs, false, false))
	}
	return out
}

func (inst *SensorSelect) GetSensorArray(sensorCount int) []*float64 {
	var result []*float64
	for i := 1; i <= sensorCount; i++ {
		sensorInputName := fmt.Sprintf("sensor%d", i)
		input := inst.ReadPinAsFloatPointer(node.InputName(sensorInputName))
		result = append(result, input)
	}
	return result
}

func (inst *SensorSelect) GetExcludeArray(sensorCount int) []bool {
	var result []bool
	for i := 1; i <= sensorCount; i++ {
		excludeInputName := fmt.Sprintf("exclude%d", i)
		exclude, _ := inst.ReadPinAsBool(node.InputName(excludeInputName))
		result = append(result, exclude)
	}
	return result
}

// Custom Node Settings Schema

type SensorSelectSettingsSchema struct {
	SensorCount schemas.Integer `json:"sensor-count"`
}

type SensorSelectSettings struct {
	SensorCount int `json:"sensor-count"`
}

func (inst *SensorSelect) buildSchema() *schemas.Schema {
	props := &SensorSelectSettingsSchema{}
	// time selection
	props.SensorCount.Title = "Number of Primary Sensors"
	props.SensorCount.Default = 2
	props.SensorCount.Minimum = 2
	props.SensorCount.Maximum = 20

	schema.Set(props)

	uiSchema := array.Map{
		"ui:order": array.Slice{"sensor-count"},
	}
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "Node Settings",
			Properties: props,
		},
		UiSchema: uiSchema,
	}
	return s
}

func (inst *SensorSelect) getSettings(body map[string]interface{}) (*SensorSelectSettings, error) {
	settings := &SensorSelectSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
