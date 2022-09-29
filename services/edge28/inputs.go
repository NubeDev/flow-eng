package edge28lib

import (
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/clients/edgerest"
	log "github.com/sirupsen/logrus"
)

func (inst *Edge28) GetDIs(p ...*points.Point) ([]*points.Point, error) {
	diValues, err := inst.client.GetDIs()
	if err != nil {
		return nil, err
	}
	var out []*points.Point
	for _, point := range p {
		if point.ObjectType == points.BinaryInput {
			for i := 1; i < 8; i++ {
				if point.ObjectID == points.ObjectID(i) {
					pntReturn := diValue(i, diValues)
					if err != nil {
						log.Error(err)
					}
					point.ValueFromRead = pntReturn
					out = append(out, point)
				}
			}
		}
	}
	return out, nil
}

func diValue(num int, values *edgerest.DI) float64 {
	if num == 1 {
		return values.Val.DI1.Val
	}
	if num == 2 {
		return values.Val.DI2.Val
	}
	if num == 3 {
		return values.Val.DI3.Val
	}
	if num == 4 {
		return values.Val.DI4.Val
	}
	if num == 5 {
		return values.Val.DI5.Val
	}
	if num == 6 {
		return values.Val.DI6.Val
	}
	if num == 7 {
		return values.Val.DI7.Val
	}
	return 0
}

func processInput(val float64, ioType string) (float64, error) {
	return getValueUI(ioType, val)
}

func (inst *Edge28) GetUIs(p ...*points.Point) ([]*points.Point, error) {
	uiValues, err := inst.client.GetUIs()
	if err != nil {
		return nil, err
	}
	var out []*points.Point
	for _, point := range p {
		if point.ObjectType == points.AnalogInput {
			for i := 1; i < 8; i++ {
				if point.ObjectID == points.ObjectID(i) {
					pntReturn, _ := inputSelection(inputValueSelection(i, uiValues), point)
					if err != nil {
						log.Error(err)
					}
					out = append(out, pntReturn)
				}
			}
		}
	}
	return out, nil
}

func inputValueSelection(num int, uiValues *edgerest.UI) float64 {
	if num == 1 {
		return uiValues.Val.UI1.Val
	}
	if num == 2 {
		return uiValues.Val.UI2.Val
	}
	if num == 3 {
		return uiValues.Val.UI3.Val
	}
	if num == 4 {
		return uiValues.Val.UI4.Val
	}
	if num == 5 {
		return uiValues.Val.UI5.Val
	}
	if num == 6 {
		return uiValues.Val.UI6.Val
	}
	if num == 7 {
		return uiValues.Val.UI7.Val
	}
	return 0
}

func inputSelection(value float64, point *points.Point) (*points.Point, error) {
	if point.IoType == points.IoTypeTemp {
		input, err := processInput(value, UITypes.THERMISTOR10KT2)
		point.ValueFromRead = input
		return point, err
	}
	if point.IoType == points.IoTypeVolts {
		input, err := processInput(value, UITypes.VOLTSDC)
		point.ValueFromRead = input
		return point, err
	}
	if point.IoType == points.IoTypeCurrent {
		input, err := processInput(value, UITypes.MILLIAMPS)
		point.ValueFromRead = input
		return point, err
	}
	if point.IoType == points.IoTypeDigital {
		input, err := processInput(value, UITypes.DIGITAL)
		point.ValueFromRead = input
		return point, err
	}
	return nil, nil
}
