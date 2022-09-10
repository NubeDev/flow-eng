package schema

func NewString(body *SettingBase) *SettingBase {
	if body == nil {
		body = &SettingBase{}
	}
	return &SettingBase{
		Type:         PropString,
		Title:        body.Title,
		Min:          body.Min,
		Max:          body.Max,
		ReadOnly:     body.ReadOnly,
		DefaultValue: body.DefaultValue,
	}
}

type SettingBase struct {
	Type         string `json:"type" default:""`
	Title        string `json:"title" default:""`
	Min          int    `json:"minLength" default:"0"`
	Max          int    `json:"maxLength" default:"500"`
	ReadOnly     bool   `json:"readOnly"`
	DefaultValue string `json:"defaultValue"`
}
