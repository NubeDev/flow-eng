package bacnet

import (
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	log "github.com/sirupsen/logrus"
)

type runnerStatus bool

var mqttSubLoop runnerStatus
var mqttPubLoop runnerStatus
var modbusLoop runnerStatus

func (inst *Server) protocolRunner() {
	if !mqttSubLoop {
		go inst.mqttSubRunner()
		mqttSubLoop = true
	}
	if !mqttPubLoop {
		go inst.mqttPubRunner()
		mqttPubLoop = true
	}
	if !modbusLoop {
		if getRunnerType() == applications.Modbus {
			go inst.modbusRunner()
			modbusLoop = true
		}
	} else {
		log.Infof("SKIP Modbus as the current poll is not finished")
	}

}

// edge28 & rubix-io input values will come from rest
func (inst *Server) edgeRunner() {

}
