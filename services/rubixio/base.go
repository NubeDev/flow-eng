package rubixIO

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	"github.com/NubeDev/flow-eng/services/clients/rubixcli"
	"github.com/NubeIO/nubeio-rubix-lib-rest-go/pkg/nube/rubixio"
)

type RubixIO struct {
	rest *rubixcli.Client
}

func New() *RubixIO {
	rest := rubixcli.New(&rubixcli.Connection{Ip: "192.168.15.191"})
	return &RubixIO{
		rest: rest,
	}
}

func (inst *RubixIO) BulkWrite(point []*points.Point) error {
	outs := inst.bulkWrite(point)
	if len(outs) > 0 {
		_, err := inst.rest.BulkWrite(inst.rest.BulkWriteBuilder(outs...))
		if err != nil {
			return err
		}
	}
	return nil
}

func (inst *RubixIO) bulkWrite(point []*points.Point) []*rubixcli.Output {
	var outs []*rubixcli.Output
	for _, p := range point {
		ioName, err := inst.uoIoNum(p)
		out := &rubixcli.Output{
			IoNumber: ioName,
			Value:    int(p.WriteValue),
		}
		if p.Enable && p.IsWriteable && err == nil {
			outs = append(outs, out)
		}
	}
	return outs
}

func (inst *RubixIO) uoIoNum(point *points.Point) (string, error) {
	if point.IsWriteable && point.ObjectType == points.AnalogOutput {
		switch point.ObjectID {
		case 1:
			return "UO1", nil
		case 2:
			return "UO2", nil
		case 3:
			return "UO3", nil
		case 4:
			return "UO4", nil
		case 5:
			return "UO5", nil
		case 6:
			return "UO6", nil
		}
	}
	if point.IsWriteable && point.ObjectType == points.BinaryOutput {
		switch point.ObjectID {
		case 1:
			return "DO1", nil
		case 2:
			return "DO2", nil
		}
	}

	return "", errors.New("rubix-io input object-id was not found")
}

func (inst *RubixIO) GetInputValue(point *points.Point, inputs *rubixio.Inputs) (float64, error) {
	if point == nil {
		return 0, errors.New("rubix-io point can not be empty")
	}
	if inputs == nil {
		return 0, errors.New("rubix-io inputs can not be empty")
	}
	ioNum, err := inst.getInputIONum(point)
	if err != nil {
		return 0, err
	}
	return inst.getInputValue(ioNum, point.IoType, inputs)

}

func (inst *RubixIO) getInputIONum(point *points.Point) (string, error) {
	if !point.IsWriteable {
		switch point.ObjectID {
		case 1:
			return "UI1", nil
		case 2:
			return "UI2", nil
		case 3:
			return "UI3", nil
		case 4:
			return "UI4", nil
		case 5:
			return "UI5", nil
		case 6:
			return "UI6", nil
		case 7:
			return "UI7", nil
		case 8:
			return "UI8", nil
		}

	}
	return "", errors.New("rubix-io input object-id was not found")
}

func (inst *RubixIO) getInputValue(ioNum string, ioType points.IoType, inputs *rubixio.Inputs) (float64, error) {
	found, temp, voltage, current, _, digital, err := rubixio.GetInputValues(inputs, ioNum)
	if err != nil {
		return 0, errors.New("rubix-io inputs can not be empty")
	}
	if found {
		switch ioType {
		case points.IoTypeTemp:
			return temp, nil
		case points.IoTypeDigital:
			return float64(digital), nil
		case points.IoTypeVolts:
			return voltage, nil
		case points.IoTypeCurrent:
			return current, nil
		}
	}
	return 0, errors.New("rubix-io input io-type was not found")

}

// DecodeInputs decode the mqtt data
func (inst *RubixIO) DecodeInputs(body []byte) (*rubixio.Inputs, error) {
	data := &rubixio.Inputs{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to decode rubix-io mqtt payload err:%s", err.Error()))
	}
	return data, nil
}
