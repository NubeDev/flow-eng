package gmail

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	FromAddress schemas.String `json:"fromAddress"`
	Password    schemas.String `json:"password"`
	ToAddress   schemas.String `json:"toAddress"`
}

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.FromAddress.Title = "From Address"
	props.FromAddress.Default = "noreply@nube-io.com"
	props.ToAddress.Title = "To Address"
	schema.Set(props)
	uiSchema := array.Map{
		"password": array.Map{
			"ui:widget": "password",
		},
	}
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "addresses and password",
			Properties: props,
		},
		UiSchema: uiSchema,
	}
	return s
}

type nodeSettings struct {
	FromAddress string `json:"fromAddress"`
	Password    string `json:"password"`
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
