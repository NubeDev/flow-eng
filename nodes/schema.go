package nodes

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/schemas"
)

func GetSchema(nodeName string) (*schemas.Schema, error) {
	for _, spec := range All() {
		if nodeName == spec.GetName() {
			return spec.GetSchema(), nil
		}
	}
	return nil, errors.New(fmt.Sprintf("no node found by name: %s", nodeName))
}
