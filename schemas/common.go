package schemas

// see react schema
// https://rjsf-team.github.io/react-jsonschema-form/

type Schema struct {
	Title      string      `json:"title"`
	Properties interface{} `json:"properties"`
	UiSchema   interface{} `json:"uiSchema"`
}

type String struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:""`
	Default  string `json:"default" default:""`
	Help     string `json:"help" default:""`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type StringLimits struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:""`
	Min      int    `json:"minLength" default:"0"`
	Max      int    `json:"maxLength" default:"0"`
	Default  string `json:"default" default:""`
	Help     string `json:"help" default:""`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type Number struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:""`
	Default  int    `json:"default" default:"0"`
	Help     string `json:"help" default:""`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

type NumberLimits struct {
	Type     string `json:"type" default:"number"`
	Title    string `json:"title" default:""`
	Min      int    `json:"minLength" default:"0"`
	Max      int    `json:"maxLength" default:"0"`
	Default  int    `json:"default" default:"0"`
	Help     string `json:"help" default:""`
	ReadOnly bool   `json:"readOnly" default:"false"`
}
