package bacnetio

import (
	"github.com/NubeDev/flow-eng/helpers/names"
	"time"
)

var runnerDelay = time.Duration(100)

func (inst *Server) protocolRunner(application names.ApplicationName) {
	go inst.writeRunner()
	if application == names.RubixIOAndModbus || application == names.Modbus {
		go inst.modbusRunner(inst.GetSettings())
	}
	if application == names.RubixIO || application == names.RubixIOAndModbus {
		go inst.rubixOutputsDispatch()
	}
	if application == names.Edge {
		go inst.edge28OutputsDispatch()
		go inst.edge28InputsRunner()
	}

}
