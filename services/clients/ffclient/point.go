package ffclient

import (
	"github.com/NubeDev/flow-eng/services/clients/ffclient/nresty"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

// AddPoint an object
func (inst *FlowClient) AddPoint(body *model.Point) (*model.Point, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Point{}).
		SetBody(body).
		Post("/api/points"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Point), nil
}

// GetPoints an objects
func (inst *FlowClient) GetPoints() ([]model.Point, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&[]model.Point{}).
		Get("/api/points"))
	if err != nil {
		return nil, err
	}
	var out []model.Point
	out = *resp.Result().(*[]model.Point)
	return out, nil
}

// GetPoint an object
func (inst *FlowClient) GetPoint(uuid string) (*model.Point, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Point{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Get("/api/points/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Point), nil
}

// DeletePoint an object
func (inst *FlowClient) DeletePoint(uuid string) (bool, error) {
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetPathParams(map[string]string{"uuid": uuid}).
		Delete("/api/points/{uuid}"))
	if err != nil {
		return false, err
	}
	return true, nil
}

// EditPoint an object
func (inst *FlowClient) EditPoint(uuid string, body *model.Point) (*model.Point, error) {
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		SetResult(&model.Point{}).
		SetPathParams(map[string]string{"uuid": uuid}).
		Patch("/api/points/{uuid}"))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Point), nil
}
