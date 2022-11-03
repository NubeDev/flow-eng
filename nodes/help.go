package nodes

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/node"
)

func NodeHelp() []*node.Help {
	var all []*node.Help
	for _, spec := range All() {
		all = append(all, spec.NodeHelp())
	}
	return all
}

func NodeHelpByName(nodeName string) (*node.Help, error) {
	var s *node.Help
	var found bool
	for _, spec := range All() {
		if nodeName == spec.GetName() {
			s = spec.NodeHelp()
			found = true
		}
	}
	if !found {
		return nil, errors.New(fmt.Sprintf("no node found by name: %s", nodeName))
	}
	return s, nil
}
