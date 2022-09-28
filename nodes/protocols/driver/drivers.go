package driver

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/names"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
)

type Networks struct {
	Networks []*Network
}

type Network struct {
	UUID            string                `json:"uuid"`
	Name            string                `json:"name"`
	ConnectionUUID  string                `json:"connectionUUID"`
	MaxNetworkCount int                   `json:"maxNetworkCount"`
	MaxDeviceCount  int                   `json:"maxDeviceCount"`
	Application     names.ApplicationName `json:"application"`
	devices         []*Device
}

type Device struct {
	UUID        string `json:"uuid"`
	NetworkUUID string `json:"networkUUID"`
	Name        string `json:"name"`
	Points      []*Point
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
	AddNetwork(body *Network) *Network
	GetNetworksByConnection(connectionUUID string) []*Network
	GetNetworks() []*Network
	GetNetwork(uuid string) *Network
	GetNetworkByName(name string) *Network
	AddDevice(networkUUID string, body *Device) *Device
	GetDevice(deviceUUID string) *Device
	GetDeviceByNetworkName(networkName string) []*Device
	GetDeviceByName(networkName, deviceName string) *Device
	AddPoint(deviceUUID string, body *Point) *Point
}

func New(n *Networks) Driver {
	return n
}

func (inst *Networks) Get() *Networks {
	return inst
}

func (inst *Networks) AddNetwork(body *Network) *Network {
	inst.Networks = append(inst.Networks, body)
	return body
}

func (inst *Networks) GetNetworks() []*Network {
	return inst.Networks
}

func (inst *Networks) GetNetworksByConnection(connectionUUID string) []*Network {
	var out []*Network
	for _, network := range inst.Networks {
		if network.ConnectionUUID == connectionUUID {
			out = append(out, network)
		}
	}
	return out
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

func (inst *Networks) GetDevice(deviceUUID string) *Device {
	for _, network := range inst.GetNetworks() {
		for _, device := range network.devices {
			if device.UUID == deviceUUID {
				return device
			}
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

func (inst *Networks) AddDevice(networkUUID string, body *Device) *Device {
	for _, network := range inst.GetNetworks() {
		if network.UUID == networkUUID {
			body.NetworkUUID = networkUUID
			fmt.Println("!!!!!!!!!!!!!!!!!!!!", networkUUID)
			pprint.Print(body)
			network.devices = append(network.devices, body)
			return body
		}
	}
	return nil
}

func (inst *Networks) AddPoint(deviceUUID string, body *Point) *Point {
	device := inst.GetDevice(deviceUUID)
	body.DeviceUUID = deviceUUID
	device.Points = append(device.Points, body)
	return body
}
