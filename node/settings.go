package node

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/boolean"
	"github.com/mitchellh/mapstructure"
)

type PropertyBase struct {
	Type     Prop        `json:"type" default:""`
	Title    string      `json:"title" default:""`
	Min      int         `json:"minLength" default:"0"`
	Max      int         `json:"maxLength" default:"500"`
	ReadOnly *bool       `json:"readOnly"`
	Value    interface{} `json:"value"`
}

type Prop string

const (
	String  Prop = "string"
	Number  Prop = "number"
	Boolean Prop = "boolean"
)

func NewProperty(args *PropertyBase) *PropertyBase {
	if args == nil {
		args = &PropertyBase{}
	}
	if args.Type == "" {
		args.Type = String
	}
	if args.Max == 0 {
		args.Max = 200
	}
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

func NewSetting(propType Prop, settingTitle string, body *PropertyBase) (*Settings, error) {
	return &Settings{
		Type:       propType,
		Title:      settingTitle,
		Properties: body,
	}, nil
}

type Settings struct {
	Type       Prop        `json:"type"`
	Title      string      `json:"title"`
	Properties interface{} `json:"properties"` // PropertyBase
}

func (n *BaseNode) GetSettings() []*Settings {
	return n.Settings
}

func (n *BaseNode) GetSetting(name string) *Settings {
	for _, setting := range n.Settings {
		if name == setting.Title {
			return setting
		}
	}
	return nil
}

func (n *BaseNode) GetPropValue(name string) (interface{}, error) {
	data := n.GetProperties(name)
	if data == nil {
		return "", errors.New(fmt.Sprintf("failed to to settings properties by name%s", name))
	}
	set := &PropertyBase{}
	err := mapstructure.Decode(n.GetProperties(name), set)
	if err != nil {
		return "", err
	}
	return set.Value, nil
}

//GetPropValueInt if there was an existing value then try and get it (would be used when node is created from json)
func (n *BaseNode) GetPropValueInt(name string, fallbackValue int) int {
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

func (n *BaseNode) GetPropValueStr(name string) (string, error) {
	data, err := n.GetPropValue(name)
	if err != nil {
		return "", err
	}
	toStr := fmt.Sprintf("%v", data)
	return toStr, nil
}

func (n *BaseNode) DecodeProperties(name string, output interface{}) error {
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

func (n *BaseNode) GetProperties(name string) interface{} {
	for _, setting := range n.Settings {
		if name == setting.Title {
			return setting.Properties
		}
	}
	return nil
}

func BuildSettings(props ...*Settings) ([]*Settings, error) {
	var out []*Settings
	var names []string
	for _, output := range props {
		out = append(out, output)
		names = append(names, output.Title)
	}
	if len(unique(names)) != len(out) { // quick hack sure there is a better way
		return nil, errors.New("the setting title must be unique")
	}
	return out, nil
}

func contains(e []string, c string) bool {
	for _, s := range e {
		if s == c {
			return true
		}
	}
	return false
}

func unique(e []string) []string {
	var r []string
	for _, s := range e {
		if !contains(r[:], s) {
			r = append(r, s)
		}
	}
	return r
}
