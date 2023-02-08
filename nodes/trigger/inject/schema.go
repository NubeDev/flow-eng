package inject

import (
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	schema.Set(props)

	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "Inject node settings",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}
