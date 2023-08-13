package functions

import (
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	Code schemas.String `json:"code"`
}

const eg = `
let in1 = Number(input.in1)
let in2 = Number(input.in2)

RQL.Result = in1+in2
`

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.Code.Title = "code"
	props.Code.Default = eg
	schema.Set(props)
	var uiSchema = map[string]map[string]string{
		"code": {
			"ui:widget": "textarea",
		},
	}
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "write function",
			Properties: props,
		},
		UiSchema: uiSchema,
	}
	return s
}
