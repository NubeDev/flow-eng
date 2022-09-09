package node

type SettingName string

type Settings struct {
	Name     SettingName `json:"name"`
	DataType DataTypes   `json:"type"`
	Value    interface{} `json:"value"`
	Config   interface{} `json:"config"`
}
