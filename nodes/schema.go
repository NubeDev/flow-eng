package nodes

import (
	"github.com/NubeDev/flow-eng/schemas"
)

func GetSchema(category string, name string) *schemas.Schema {
	for _, spec := range All() {
		if category == spec.GetInfo().Category && name == spec.GetName() {
			return spec.GetSchema()
		}
	}
	return nil
}
