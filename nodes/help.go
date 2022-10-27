package nodes

import "github.com/NubeDev/flow-eng/node"

func NodeHelp() []*node.Help {
	var all []*node.Help
	for _, spec := range All() {
		all = append(all, spec.NodeHelp())
	}
	return all
}
