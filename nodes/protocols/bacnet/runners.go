package bacnet

import (
	"github.com/NubeDev/flow-eng/helpers/names"
	"time"
)

var runnerDelay = time.Duration(100)

func (inst *Server) protocolRunner(application names.ApplicationName) {
	go inst.writeRunner()
	if application == names.Modbus {
		go inst.modbusRunner()
	}
	if application == names.RubixIO || application == names.RubixIOAndModbus {
		go inst.rubixOutputsDispatch()
	}
	if application == names.Edge {
		go inst.edge28OutputsDispatch()
		go inst.edge28InputsRunner()
	}

}
