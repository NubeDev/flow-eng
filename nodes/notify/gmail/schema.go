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
	var password = map[string]map[string]string{
		"address": {
			"ui:widget":      "password",
			"ui:title":       "Age of person",
			"ui:description": "(earthian year)",
		},
	}

	s := &schemas.Schema{
		Title:      "gmail",
		Properties: m,
		UiSchema:   password,
	}
	return s
}
