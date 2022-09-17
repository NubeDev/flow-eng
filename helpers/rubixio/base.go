package rubix

import (
	"github.com/NubeIO/nubeio-rubix-lib-rest-go/pkg/nube/rubixio"
)

type Instance struct{}

type Output struct {
	IoNumber string
	Value    int
	Enable   bool
}

func (inst *Instance) syncOutputs(dev []*Output) (bulk []rubixio.BulkWrite) {
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
