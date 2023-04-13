package bacnetio

import (
	"time"
)

var runnerDelay = time.Duration(500)
var modbusDelay = time.Duration(50)

func (inst *Server) protocolRunner() {
	go inst.modbusRunner(inst.GetSettings())
}
