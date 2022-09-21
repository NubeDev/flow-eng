package rubixcli

import (
	"github.com/NubeDev/flow-eng/services/clients/rubixcli/nresty"
)

type Ping struct {
	Ok     bool `json:"ok"`
	PigPIO bool `json:"pigio_is_running"`
}

type BulkResponse struct {
	Ok bool `json:"ok"`
}

func (inst *Client) Ping() (*Ping, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&Ping{}).
		Get("/api/system/ping"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*Ping), nil
}
