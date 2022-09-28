package flow

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
)

func (inst *Network) runner() {

	for _, net := range inst.pool.GetNetworks() {
		pprint.Print(net)

		//if net ==
	}

}
