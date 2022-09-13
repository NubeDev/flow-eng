package node

import (
	"errors"
	"fmt"
)

func (n *BaseNode) OverrideInputValue(name InputName, value interface{}) error {
	in := n.GetInput(name)
	if in == nil {
		return errors.New(fmt.Sprintf("failed to find port%s", name))
	}
	if in.Connection != nil {
		in.Connection.OverrideValue = value
	} else {
		return errors.New(fmt.Sprintf("this node has no inputs"))
	}
	return nil
}
