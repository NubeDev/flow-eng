package rubixIO

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-rest-go/pkg/nube/rubixio"
)

type Output struct {
	IoNumber string
	Value    int
	Enable   bool
}

type RubixIO struct{}

// DecodeInputs decode the mqtt data
func (inst *RubixIO) DecodeInputs(body []byte) (*rubixio.Inputs, error) {
	data := &rubixio.Inputs{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to decode rubix-io mqtt payload err:%s", err.Error()))
	}
	return data, nil
}

func (inst *RubixIO) syncOutputs(dev []*Output) (bulk []rubixio.BulkWrite) {
	for _, output := range dev {
		bulkWrite := rubixio.BulkWrite{
			IONum: output.IoNumber,
			Value: output.Value,
		}
		if output.Enable {
			bulk = append(bulk, bulkWrite)
		}
	}
	return
}
