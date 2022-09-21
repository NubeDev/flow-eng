package rubixcli

import (
	"github.com/NubeDev/flow-eng/services/clients/rubixcli/nresty"
	"github.com/NubeIO/nubeio-rubix-lib-rest-go/pkg/nube/rubixio"
)

type Output struct {
	IoNumber string
	Value    int
}

func (inst *Client) BulkWriteBuilder(dev ...*Output) (bulk []rubixio.BulkWrite) {
	for _, output := range dev {
		bulkWrite := rubixio.BulkWrite{
			IONum: output.IoNumber,
			Value: output.Value,
		}
		bulk = append(bulk, bulkWrite)
	}
	return
}

func (inst *Client) BulkWrite(body []rubixio.BulkWrite) (*BulkResponse, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&BulkResponse{}).
		SetBody(body).
		Post("/api/outputs/bulk"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*BulkResponse), nil
}
