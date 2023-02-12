package modbuscli

import (
	"errors"
	"github.com/NubeIO/nubeio-rubix-lib-modbus-go/modbus"
	"time"
)

type Modbus struct {
	IsSerial             bool
	Address              string
	Port, Slave, Timeout int
	Serial               *modbus.Serial
	client               *modbus.Client
}

func (inst *Modbus) Init(opts *Modbus) (*Modbus, error) {
	if opts.Port == 0 {
		opts.Port = 502
	}
	if opts.Timeout <= 100 {
		opts.Timeout = 100
	}
	mbClient := &modbus.Client{
		HostIP:   opts.Address,
		HostPort: opts.Port,
		IsSerial: opts.IsSerial,
		Serial:   opts.Serial,
	}

	mbClient, err := mbClient.New()
	if err != nil {
		return nil, err
	}
	mbClient.RTUClientHandler.Timeout = time.Duration(opts.Timeout) * time.Millisecond
	m := &Modbus{
		IsSerial: opts.IsSerial,
		Address:  opts.Address,
		Port:     opts.Port,
		Slave:    opts.Slave,
		client:   mbClient,
	}
	return m, nil
}

func (inst *Modbus) SetSlave(slave int) error {
	if slave == 0 {
		return errors.New("no modbus slave address was passed in")
	}
	if inst.IsSerial {
		inst.client.RTUClientHandler.SetSlave(byte(slave))
	} else {
		inst.client.TCPClientHandler.SetSlave(byte(slave))
	}
	return nil
}
