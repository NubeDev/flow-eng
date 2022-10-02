package mathematics

import (
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
)

type nodeSchema struct {
	Sch schemas.EnumString `json:"function"`
}

// NODEs will single in/out
const (
	acos  = "acos"
	asin  = "asin"
	atan  = "atan"
	cbrt  = "cbrt"
	cos   = "cos"
	exp   = "exp"
	log   = "log"
	log10 = "log10"
	sin   = "sin"
	sqrt  = "sqrt"
	tan   = "tan"
)

func buildSchema() *schemas.Schema {
	props := &nodeSchema{}
	props.Sch.Title = "function"
	props.Sch.Default = acos
	props.Sch.Options = []string{acos, asin, atan, cbrt, cos, exp, log, log10, sin, sqrt, tan}
	props.Sch.EnumName = []string{acos, asin, atan, cbrt, cos, exp, log, log10, sin, sqrt, tan}
	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "function",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}
