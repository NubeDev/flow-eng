package rest

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	Address schemas.Enum `json:"method"`
}

func buildSchema() *schemas.Schema {
	m := &nodeSchema{}
	m.Address.Title = "method"
	m.Address.Options = []string{"GET", "POST"}
	m.Address.EnumName = []string{"GET", "POST"}
	pprint.PrintJOSN(m)
	schema.Set(m)

	s := &schemas.Schema{
		Title:      "gmail",
		Properties: m,
		UiSchema:   nil,
	}
	return s
}
