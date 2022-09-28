package ffclient

import (
	"fmt"
	"github.com/NubeDev/flow-eng/services/clients/ffclient/nresty"
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
)

// AddDevice an object
func (inst *Client) AddDevice(device *model.Device) (*model.Device, error) {
	url := fmt.Sprintf("/api/devices")
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Device{}).
		SetBody(device).
		Post(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Device), nil
}

// GetFirstDevice first object
func (inst *Client) GetFirstDevice(withPoints ...bool) (*model.Device, error) {
	devices, err := inst.GetDevices(withPoints...)
	if err != nil {
		return nil, err
	}
	for _, device := range devices {
		return &device, err
	}
	return nil, err
}

// GetDevices all objects
func (inst *Client) GetDevices(withPoints ...bool) ([]model.Device, error) {
	url := fmt.Sprintf("/api/devices")
	if len(withPoints) > 0 {
		if withPoints[0] == true {
			url = fmt.Sprintf("/api/devices/?with_points=true")
		}
	}
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&[]model.Device{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	var out []model.Device
	out = *resp.Result().(*[]model.Device)
	return out, nil
}

// GetDevice an object
func (inst *Client) GetDevice(uuid string, withPoints ...bool) (*model.Device, error) {
	url := fmt.Sprintf("/api/devices/%s", uuid)
	if len(withPoints) > 0 {
		if withPoints[0] == true {
			url = fmt.Sprintf("/api/devices/%s?with_points=true", uuid)
		}
	}
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Device{}).
		Get(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Device), nil
}

// EditDevice edit an object
func (inst *Client) EditDevice(uuid string, device *model.Device) (*model.Device, error) {
	url := fmt.Sprintf("/api/devices/%s", uuid)
	resp, err := nresty.FormatRestyResponse(inst.client.R().
		SetResult(&model.Device{}).
		SetBody(device).
		Patch(url))
	if err != nil {
		return nil, err
	}
	return resp.Result().(*model.Device), nil
}

// DeleteDevice an object
func (inst *Client) DeleteDevice(uuid string) (bool, error) {
	_, err := nresty.FormatRestyResponse(inst.client.R().
		SetPathParams(map[string]string{"uuid": uuid}).
		Delete("/api/devices/{uuid}"))
	if err != nil {
		return false, err
	}
	return true, nil
}
