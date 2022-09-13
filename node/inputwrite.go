package node

import (
	"errors"
	"fmt"
)

func (n *BaseNode) OverrideInputValue(name InputName, value interface{}) error {
	input := n.GetInput(name)
	if input == nil {
		return errors.New(fmt.Sprintf("failed to find port %s", name))
	}
	if input.Connection != nil {
		input.Connection.OverrideValue = value
	} else {
		return errors.New(fmt.Sprintf("this node has no inputs"))
	}
	return nil
}
