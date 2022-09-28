package driver

import (
	"github.com/NubeDev/flow-eng/helpers/names"
)

type Networks struct {
	Networks []*Network
}

type Network struct {
	UUID            string                `json:"uuid"`
	Name            string                `json:"name"`
	MaxNetworkCount int                   `json:"maxNetworkCount"`
	MaxDeviceCount  int                   `json:"maxDeviceCount"`
	Application     names.ApplicationName `json:"application"`
	devices         []*Device
}

type Device struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Points []*Point
}

type Point struct {
	UUID       string    `json:"uuid"`
	DeviceUUID string    `json:"deviceUUID"`
	Name       string    `json:"name"`
	ReadOnly   *bool     `json:"readOnly"`
	Priority   *Priority `json:"priority"`
}

type Driver interface {
	Get() *Networks
	AddNetwork(body *Network)
	GetNetworks() []*Network
	GetNetwork(uuid string) *Network
	GetNetworkByName(name string) *Network
	GetDeviceByNetworkName(networkName string) []*Device
	GetDeviceByName(networkName, deviceName string) *Device
}

func New(n *Networks) Driver {
	return n
}

func (inst *Networks) Get() *Networks {
	return inst
}

func (inst *Networks) AddNetwork(body *Network) {
	inst.Networks = append(inst.Networks, body)
}

func (inst *Networks) GetNetworks() []*Network {
	return inst.Networks
}

func (inst *Networks) GetNetwork(uuid string) *Network {
	for _, network := range inst.Networks {
		if network.UUID == uuid {
			return network
		}
	}
	return nil
}

func (inst *Networks) GetNetworkByName(name string) *Network {
	for _, network := range inst.Networks {
		if network.Name == name {
			return network
		}
	}
	return nil
}

func (inst *Networks) GetDeviceByNetworkName(networkName string) []*Device {
	for _, network := range inst.GetNetworks() {
		if network.Name == networkName {
			return network.devices
		}
	}
	return nil
}

func (inst *Networks) GetDeviceByName(networkName, deviceName string) *Device {
	for _, network := range inst.GetNetworks() {
		if network.Name == networkName {
			for _, device := range network.devices {
				if device.Name == deviceName {
					return device
				}
			}
		}
	}
	return nil
}
