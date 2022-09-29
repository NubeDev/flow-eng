package gmail

import (
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	Address schemas.String `json:"address"`
}

func buildSchema() *schemas.Schema {
	m := &nodeSchema{}
	m.Address.Title = "address"
	m.Address.Default = "noreply@nube-io.com"
	schema.Set(m)
	s := &schemas.Schema{
		Title:      "gmail",
		Properties: m,
		UiSchema:   nil,
	}
	return s
}
