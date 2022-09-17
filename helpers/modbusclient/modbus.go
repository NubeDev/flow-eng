package modbusclient

import (
	"fmt"
	"github.com/NubeIO/nubeio-rubix-lib-modbus-go/modbus"
)

func modbusInit() (*modbus.Client, error) {

	mbClient := &modbus.Client{
		HostIP:   "192.168.15.202",
		HostPort: 502,
	}
	mbClient, err := mbClient.New()
	if err != nil {
		return nil, err
	}
	mbClient.TCPClientHandler.Address = fmt.Sprintf("%s:%d", "192.168.15.202", 502)
	mbClient.TCPClientHandler.SlaveID = byte(1)

	return mbClient, nil
}
