package node

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/boolean"
	"github.com/NubeDev/flow-eng/helpers/integer"
	"github.com/mitchellh/mapstructure"
)

const (
	InputCount = "Inputs Count"
)

type SettingOptions struct {
	Type  PropType
	Title Title
	Min   int
	Max   int
	Value interface{}
}

func NewSetting(body *Spec, opts *SettingOptions) (base *PropertyBase, setting *Settings, value interface{}, err error) {
	if opts == nil {
		opts = &SettingOptions{}
	}
	var min = opts.Min
	var max = opts.Max
	var dataType PropType
	var title = opts.Title
	if title == "" {
		return nil, nil, 0, errors.New("title can not be empty")
	}
	if min == 0 {
		min = 1
	}
	if max == 0 {
		max = 1
	}
	if opts.Type == "" {
		panic("new-setting: dataType can not be empty")
	} else {
		dataType = opts.Type
	}
	if dataType == String {

	} else if dataType == Number {
		var getValue = min
		getValue = body.GetPropValueInt(opts.Title, min)
		base = &PropertyBase{
			Min: integer.New(min),
			Max: integer.New(max),
		}
		setting, err = Setting(dataType, title, base)
		if err != nil {
			return nil, nil, 0, err
		}
		return base, setting, getValue, err
	} else if dataType == Array {
		getArray, err := body.GetPropValue(opts.Title)
		base = &PropertyBase{}
		setting, err = Setting(dataType, title, base)
		if err != nil {
			return nil, nil, 0, err
		}
		return base, setting, getArray, err
	}

	return nil, nil, 0, errors.New("new-setting: no valid settings type")
}

type PropertyBase struct {
	Type     PropType    `json:"type" default:""`
	Title    Title       `json:"title" default:""`
	Min      *int        `json:"min,omitempty" default:"0"`
	Max      *int        `json:"max,omitempty" default:"500"`
	ReadOnly *bool       `json:"readOnly,omitempty"`
	Value    interface{} `json:"value"`
}

func settingInt() {

}

type Title string
type PropType string

const (
	Array   PropType = "array"
	String  PropType = "string"
	Number  PropType = "number"
	Boolean PropType = "boolean"
)

func NewProperty(args *PropertyBase) *PropertyBase {
	if args == nil {
		args = &PropertyBase{}
	}
	if args.Type == "" {
		args.Type = String
	}
	//if args.Max == 0 {
	//	args.Max = 200
	//}
	if boolean.IsNil(args.ReadOnly) {
		args.ReadOnly = boolean.NewFalse()
	}
	return &PropertyBase{
		Type:     args.Type,
		Title:    args.Title,
		Min:      args.Min,
		Max:      args.Max,
		ReadOnly: args.ReadOnly,
		Value:    args.Value,
	}
}

func Setting(propType PropType, settingTitle Title, body *PropertyBase) (*Settings, error) {
	return &Settings{
		Type:       propType,
		Title:      settingTitle,
		Properties: body,
	}, nil
}

type Settings struct {
	Type       PropType    `json:"type"`
	Title      Title       `json:"title"`
	Properties interface{} `json:"properties"` // PropertyBase
}

func (n *Spec) GetSettings() []*Settings {
	return n.Settings
}

func (n *Spec) GetSetting(name Title) *Settings {
	for _, setting := range n.Settings {
		if name == setting.Title {
			return setting
		}
	}
	return nil
}

func (n *Spec) SetPropertiesValue(name Title, value interface{}) error {
	data := n.GetProperties(name)
	if data == nil {
		return errors.New(fmt.Sprintf("failed to to settings properties by name%s", name))
	}

	properties := &PropertyBase{}
	err := n.DecodeProperties(name, properties)
	if err != nil {
		return err
	}
	setting := n.GetSetting(name)
	properties.Value = value

	setting.Properties = properties
	return nil
}

func (n *Spec) GetPropValue(name Title) (interface{}, error) {
	data := n.GetProperties(name)
	if data == nil {
		return "", errors.New(fmt.Sprintf("failed to to settings properties by name%s", name))
	}
	set := &PropertyBase{}
	err := mapstructure.Decode(data, set)
	if err != nil {
		return "", err
	}
	return set.Value, nil
}

// GetPropValueInt if there was an existing value then try and get it (would be used when node is created from json)
func (n *Spec) GetPropValueInt(name Title, fallbackValue int) int {
	data, err := n.GetPropValue(name)
	if err != nil {
		return 0
	}
	i, ok := data.(int)
	if !ok {
		return fallbackValue
	}
	return i
}

func (n *Spec) GetPropValueStr(name Title) (string, error) {
	data, err := n.GetPropValue(name)
	if err != nil {
		return "", err
	}
	toStr := fmt.Sprintf("%v", data)
	return toStr, nil
}

func (n *Spec) DecodeProperties(name Title, output interface{}) error {
	data := n.GetProperties(name)
	if data == nil {
		return errors.New(fmt.Sprintf("failed to find settings properties by name:%s", name))
	}
	err := mapstructure.Decode(data, output)
	if err != nil {
		return errors.New(fmt.Sprintf("mapstructure.Decode err:%s", err))
	}
	return nil
}

func (n *Spec) GetProperties(name Title) interface{} {
	for _, setting := range n.Settings {
		if name == setting.Title {
			return setting.Properties
		}
	}
	return nil
}

func BuildSettings(props ...*Settings) ([]*Settings, error) {
	var out []*Settings
	var names []Title
	for _, output := range props {
		out = append(out, output)
		names = append(names, output.Title)
	}
	if len(unique(names)) != len(out) { // quick hack sure there is a better way
		return nil, errors.New("the setting title must be unique")
	}
	return out, nil
}

func contains(e []Title, c Title) bool {
	for _, s := range e {
		if s == c {
			return true
		}
	}
	return false
}

func unique(e []Title) []Title {
	var r []Title
	for _, s := range e {
		if !contains(r[:], s) {
			r = append(r, s)
		}
	}
	return r
}
