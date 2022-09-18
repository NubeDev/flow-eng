package bacnet

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/modbusclient"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bstore"
	"github.com/NubeIO/nubeio-rubix-lib-modbus-go/modbus"
)

// processMessage are the messages from the bacnet-server via the mqtt-broker

func (inst *Server) mqttRunner() {
	go func() {
		msg, ok := inst.bus().Recv()
		if ok {
			payload := bstore.NewPayload()
			err := payload.NewMessage(msg)
			if err != nil {
				return
			}
			t, id := payload.GetObject()
			pnt := inst.db().GetPointByObject(t, id)
			inst.db().WritePointValue(pnt.UUID, payload.GetPresentValue())
		}
	}()
}

// edge28 & rubix-io input values will come from rest
func (inst *Server) edgeRunner() {

}

// modbus will come from polling
// this is to only work for the IO-16
func (inst *Server) modbusRunner() {

	cli := modbusclient.New(&modbusclient.Modbus{})

	for _, point := range inst.db().GetPointsByApplication(applications.Modbus) {
		if point.IsIO && point.IsWriteable {
			//slaveId := 0
			addr, t := cli.BuildInput(point.IoType, point.ObjectID)
			fmt.Println(addr, t)

			//slaveId = addr

		}

	}

}

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
