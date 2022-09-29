package edge28lib

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
)

func (inst *Edge28) WriteUO(point *points.Point) (*points.Point, error) {
	pri := points.GetHighest(point.WriteValue)
	if pri != nil {
		ioNum := uoSelection(point.ObjectID)
		isD := isDO(point.IoType)
		ioType := points.IoTypeVolts
		if isD {
			ioType = points.IoTypeDigital
		}
		writeValue, err := processOutput(pri.Value, string(ioType), isD)
		fmt.Println(pri.Value, writeValue, ioNum, ioType, isD)
		if err != nil {
			return nil, err
		}
		_, err = inst.client.WriteUO(ioNum, writeValue)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func processOutput(value float64, ioType string, isDo bool) (float64, error) {
	value = limitValue(ioType, value)
	var err error
	var wv float64
	if isDo {
		if value == 0 {
			return 100, nil
		}
		if value == 1 {
			return 0, nil
		}
		return 0, errors.New("invalid value for the DO write")
	} else {
		wv, err = getValueUO(ioType, value)
	}
	return wv, err
}

func isDO(t points.IoType) bool {
	if t == points.IoTypeDigital {
		return true
	}
	return false
}

func uoSelection(num points.ObjectID) string {
	if num == 1 {
		return fmt.Sprintf("UO%d", 1)
	}
	if num == 2 {
		return fmt.Sprintf("UO%d", 2)
	}
	if num == 3 {
		return fmt.Sprintf("UO%d", 3)
	}
	if num == 4 {
		return fmt.Sprintf("UO%d", 4)
	}
	if num == 5 {
		return fmt.Sprintf("UO%d", 5)
	}
	if num == 6 {
		return fmt.Sprintf("UO%d", 6)
	}
	if num == 7 {
		return fmt.Sprintf("UO%d", 7)
	}
	return ""
}

func (inst *Edge28) WriteDO(point *points.Point) (*points.Point, error) {
	pri := points.GetHighest(point.WriteValue)
	if pri != nil {
		ioNum := doSelection(point.ObjectID)
		_, err := inst.client.WriteDO(ioNum, pri.Value)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func doSelection(num points.ObjectID) string {
	if num == 1 {
		return fmt.Sprintf("DO%d", 1)
	}
	if num == 2 {
		return fmt.Sprintf("DO%d", 2)
	}
	if num == 3 {
		return fmt.Sprintf("DO%d", 3)
	}
	if num == 4 {
		return fmt.Sprintf("DO%d", 4)
	}
	if num == 5 {
		return fmt.Sprintf("DO%d", 5)
	}
	if num == 6 {
		return fmt.Sprintf("R%d", 1)
	}
	if num == 7 {
		return fmt.Sprintf("R%d", 2)
	}
	return ""
}
