package bacnetio

import (
	"time"
)

var runnerDelay = time.Duration(100)

func (inst *Server) protocolRunner() {
	go inst.writeRunner()
	go inst.modbusRunner(inst.GetSettings())
}
