package flow

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
)

func (inst *Network) runner() {
	for _, net := range inst.pool.GetNetworksByConnection(inst.connectionUUID) {
		if net.UUID == inst.networkUUID {
			pprint.Print(net)
		}
	}

}
