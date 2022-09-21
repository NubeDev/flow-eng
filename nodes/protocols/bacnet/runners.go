package bacnet

import (
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	log "github.com/sirupsen/logrus"
	"time"
)

var runnerDelay = time.Duration(100)

type runnerStatus bool

var mqttSubLoop runnerStatus
var mqttPubLoop runnerStatus
var modbusLoop runnerStatus
var rubixIOLoop runnerStatus

func (inst *Server) protocolRunner() {
	if !mqttPubLoop {
		go inst.writeRunner()
		mqttPubLoop = true
	}
	if !modbusLoop {
		if getRunnerType() == applications.Modbus {
			go inst.modbusRunner()
			modbusLoop = true
		}
	} else {
		if getRunnerType() == applications.Modbus {
			log.Infof("SKIP Modbus as the current poll is not finished")
		}
	}

	if !rubixIOLoop {
		if getRunnerType() == applications.RubixIO {
			go inst.rubixOutputsRunner()
			rubixIOLoop = true
		}
	}

}
