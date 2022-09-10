package schema

func NewNumber(body *SettingBase) *SettingBase {
	if body == nil {
		body = &SettingBase{}
	}
	return &SettingBase{
		Type:         PropNum,
		Title:        body.Title,
		Min:          body.Min,
		Max:          body.Max,
		ReadOnly:     body.ReadOnly,
		DefaultValue: body.DefaultValue,
	}
}
