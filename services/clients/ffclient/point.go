package ffclient

import (
	"fmt"
	"github.com/NubeDev/flow-eng/services/clients/ffclient/nresty"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

type Priority struct {
	P1  *float64 `json:"_1,omitempty"`
	P2  *float64 `json:"_2,omitempty"`
	P3  *float64 `json:"_3,omitempty"`
	P4  *float64 `json:"_4,omitempty"`
	P5  *float64 `json:"_5,omitempty"`
	P6  *float64 `json:"_6,omitempty"`
	P7  *float64 `json:"_7,omitempty"`
	P8  *float64 `json:"_8,omitempty"`
	P9  *float64 `json:"_9,omitempty"`
	P10 *float64 `json:"_10,omitempty"`
	P11 *float64 `json:"_11,omitempty"`
	P12 *float64 `json:"_12,omitempty"`
	P13 *float64 `json:"_13,omitempty"`
	P14 *float64 `json:"_14,omitempty"`
	P15 *float64 `json:"_15,omitempty"`
	P16 *float64 `json:"_16,omitempty"`
}

// PointWriteByName an object /api/points/name?network_name=modbus&device_name=modbus&point_name=modbus
func (inst *FlowClient) PointWriteByName(networkName string, deviceName string, pointName string, pri *Priority) (*model.Point, error) {
	url := fmt.Sprintf("/api/points/name/?network_name=%s&device_name=%s&point_name=%s", networkName, deviceName, pointName)
	body := map[string]interface{}{
		"priority": pri,
	}
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		SetResult(&model.Point{}).
		Patch(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Point), nil
}

// PointWrite write a point by its uuid
func (inst *FlowClient) PointWrite(uuid string, pri *Priority) (*model.Point, error) {
	url := fmt.Sprintf("/api/points/write/%s", uuid)
	body := map[string]interface{}{
		"priority": pri,
	}
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetBody(body).
		SetResult(&model.Point{}).
		Patch(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Point), nil
}

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
