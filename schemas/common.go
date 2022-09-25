package schemas

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
