package gmail

import (
	"encoding/json"

	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	FromAddress schemas.String `json:"fromAddress"`
	Token       schemas.String `json:"token"`
	ToAddress   schemas.String `json:"toAddress"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.FromAddress.Title = "From Address"
	props.FromAddress.Default = "noreply@nube-io.com"
	props.ToAddress.Title = "To Address"
	schema.Set(props)
	uiSchema := array.Map{
		"token": array.Map{
			"ui:widget": "password",
		},
	}
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "addresses and token",
			Properties: props,
		},
		UiSchema: uiSchema,
	}
	return s
}

type nodeSettings struct {
	FromAddress string `json:"fromAddress"`
	Token       string `json:"token"`
	ToAddress   string `json:"toAddress"`
}

func getSettings(body map[string]interface{}) (*nodeSettings, error) {
	settings := &nodeSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
