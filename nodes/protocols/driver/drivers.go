package driver

import (
	"errors"
	"github.com/NubeDev/flow-eng/db"
	"github.com/NubeDev/flow-eng/helpers/names"
)

type Network struct {
	UUID            string                `json:"uuid"`
	Name            string                `json:"name"`
	MaxNetworkCount int                   `json:"maxNetworkCount"`
	MaxDeviceCount  int                   `json:"maxDeviceCount"`
	Application     names.ApplicationName `json:"application"`
	devices         []*Device
	Storage         db.DB
}

type Device struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	points []*Point
}

type Point struct {
	UUID     string    `json:"uuid"`
	Name     string    `json:"name"`
	Priority *Priority `json:"priority"`
}

type Networks interface {
	Get() *Network

	GetDevices() []*Device
	GetDevice(uuid string) *Device
	AddDevice(dev *Device)

	AddPoint(deviceUUID string, body *Point) error
	GetPoints() []*Point
}

func New(driver *Network) Networks {
	return driver
}

func (inst *Network) Get() *Network {
	return inst
}

func (inst *Network) GetDevices() []*Device {
	return inst.devices
}

func (inst *Network) AddDevice(body *Device) {
	inst.devices = append(inst.devices, body)
}

func (inst *Network) GetDevice(uuid string) *Device {
	for _, device := range inst.devices {
		if device.UUID == uuid {
			return device
		}
	}
	return nil
}

func (inst *Network) GetPoints() []*Point {
	var out []*Point
	for _, device := range inst.devices {
		out = append(out, device.points...)
	}
	return out
}

func (inst *Network) AddPoint(deviceUUID string, body *Point) error {
	if inst.MaxDeviceCount == 1 {
		devices := inst.GetDevices()
		if len(devices) > 0 {
			deviceUUID = devices[0].UUID
		}
	}
	dev := inst.GetDevice(deviceUUID)
	if dev == nil {
		return errors.New("failed to find device")
	}
	dev.points = append(dev.points, body)
	return nil
}
