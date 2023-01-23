package schemas

// see react schema
// https://rjsf-team.github.io/react-jsonschema-form/

type SchemaBody struct {
	Title      string      `json:"title"`
	Properties interface{} `json:"properties"`
}

type Schema struct {
	Schema   SchemaBody  `json:"schema"`
	UiSchema interface{} `json:"uiSchema"`
}

type Boolean struct {
	Type      string   `json:"type" default:"boolean"`
	Title     string   `json:"title" default:""`
	Default   bool     `json:"default" default:"false"`
	Help      string   `json:"help" default:""`
	ReadOnly  bool     `json:"readOnly" default:"false"`
	EnumNames []string `json:"enumNames,omitempty"`
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
	Type     string  `json:"type" default:"number"`
	Title    string  `json:"title" default:""`
	Default  float64 `json:"default" default:"0"`
	Help     string  `json:"help" default:""`
	ReadOnly bool    `json:"readOnly" default:"false"`
	Minimum  float64 `json:"minimum" default:"0"`
	Maximum  float64 `json:"maximum" default:"10000"`
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

type EnumString struct {
	Type     string   `json:"type" default:"string"`
	Title    string   `json:"title" default:""`
	Default  string   `json:"default" default:""`
	Options  []string `json:"enum" default:"[]"`
	EnumName []string `json:"enumNames" default:"[]"`
}

type EnumInt struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:""`
	Default  int    `json:"default" default:""`
	Options  []int  `json:"enum" default:"[]"`
	EnumName []int  `json:"enumNames" default:"[]"`
}

type Integer struct {
	Type     string  `json:"type" default:"integer"`
	Title    string  `json:"title" default:""`
	Default  int     `json:"default" default:"0"`
	Help     string  `json:"help" default:""`
	ReadOnly bool    `json:"readOnly" default:"false"`
	Minimum  float64 `json:"minimum" default:"0"`
	Maximum  float64 `json:"maximum" default:"100"`
}
