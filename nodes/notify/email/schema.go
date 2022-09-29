package email

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	Address schemas.String `json:"address"`
}

func buildSchema() *schemas.Schema {
	m := &nodeSchema{}
	m.Address.Title = "address"
	pprint.PrintJOSN(m)

	schema.Set(m)

	s := &schemas.Schema{
		Title:      "email",
		Properties: m,
		UiSchema:   nil,
	}

	return s
}
