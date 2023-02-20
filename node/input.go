package node

import (
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/uuid"
)

type Input struct {
	Name             InputName        `json:"name"` // in1
	DataType         DataTypes        `json:"type"` // int
	Connection       *InputConnection `json:"link,omitempty"`
	value            interface{}
	updated          bool // if the input updated or node
	values           array.ArrStore
	uuid             uuid.Value
	direction        Direction
	connector        *Connector
	Help             InputHelp `json:"help"`
	FolderExport     bool      `json:"folderExport"`
	HideInput        bool      `json:"hideInput"`
	Position         int       `json:"position"`
	OverridePosition bool      `json:"overridePosition"`
	SettingName      string    `json:"settingName"`
	PreventOverride  bool      `json:"preventOverride"`
}

func newInput(body *Input) *Input {
	var values array.ArrStore
	return &Input{
		body.Name,
		body.DataType,
		body.Connection,
		nil,
		false,
		values,
		uuid.New(),
		DirectionInput,
		nil,
		"",
		false,
		false,
		0,
		false,
		body.SettingName,
		false,
	}
}

func (p *Input) UUID() uuid.Value {
	return p.uuid
}

func (p *Input) Direction() Direction {
	return p.direction
}

func (p *Input) Connectors() []*Connector {
	if p.connector == nil {
		return []*Connector{}
	}
	return []*Connector{p.connector}
}

func (p *Input) GetValue() interface{} {
	if p.value == nil {
		if p.Connection.NodeID == "" { // don't use default or override values if there is a link connected
			if p.Connection.OverrideValue != nil {
				return p.Connection.OverrideValue
			}
			if p.Connection.DefaultValue != nil {
				return p.Connection.DefaultValue
			}
		}
	}
	return p.value
}

func (p *Input) SetValue(value interface{}) {
	p.value = value
}
