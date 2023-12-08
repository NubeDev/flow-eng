package ffclient

import (
	"errors"
	"github.com/NubeDev/flow-eng/services/clients/ffclient/nresty"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/model"
)

func (inst *Client) Ping() error {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Network{}).
		Get("/api/system/ping"))
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return errors.New("failed to ping server")
	}
	return nil
}
