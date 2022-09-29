package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/names"
	"time"
)

var runnerDelay = time.Duration(100)

type runnerStatus bool

var mqttSubLoop runnerStatus
var mqttPubLoop runnerStatus
var modbusLoop runnerStatus
var rubixIOLoop runnerStatus
var edgeIOLoop runnerStatus

func (inst *Server) protocolRunner() {
	gt := getApplication()
	if !mqttPubLoop {
		go inst.writeRunner()
		mqttPubLoop = true
	}
	if !modbusLoop {
		if gt == names.Modbus {
			go inst.modbusRunner()
			modbusLoop = true
		}
	}
	if !rubixIOLoop {
		if gt == names.RubixIO || gt == names.RubixIOAndModbus {
			go inst.rubixOutputsDispatch()
			rubixIOLoop = true
		}
	}
	if !edgeIOLoop {
		if gt == names.Edge {
			go inst.edge28OutputsDispatch()
			edgeIOLoop = true
		}
	}

}
