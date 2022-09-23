package bacnet

import (
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"time"
)

var runnerDelay = time.Duration(100)

type runnerStatus bool

var mqttSubLoop runnerStatus
var mqttPubLoop runnerStatus
var modbusLoop runnerStatus
var rubixIOLoop runnerStatus

func (inst *Server) protocolRunner() {
	gt := getApplication()
	if !mqttPubLoop {
		go inst.writeRunner()
		mqttPubLoop = true
	}
	if !modbusLoop {
		if gt == applications.Modbus {
			go inst.modbusRunner()
			modbusLoop = true
		}
	}
	if !rubixIOLoop {
		if gt == applications.RubixIO || gt == applications.RubixIOAndModbus {
			go inst.rubixDispatch()
			rubixIOLoop = true
		}
	}

}
