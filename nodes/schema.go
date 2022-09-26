package nodes

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/schemas"
)

func GetSchema(nodeName string) (interface{}, error) {
	s := &schemas.Schema{}
	for _, spec := range All() {
		if nodeName == spec.GetName() {
			s = spec.GetSchema()
		}
	}
	res := map[string]interface{}{
		"schema": s,
	}
	if s == nil {
		return nil, errors.New(fmt.Sprintf("no node found by name: %s", nodeName))
	}
	return res, nil

}
