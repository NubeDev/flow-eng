package node

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/schema"
	"github.com/mitchellh/mapstructure"
)

type Settings struct {
	Type       string      `json:"type"`
	Title      string      `json:"title"`
	Properties interface{} `json:"properties"`
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

func (n *BaseNode) GetPropValueStr(name string) (string, error) {
	data := n.GetProperties(name)
	if data == nil {
		return "", errors.New(fmt.Sprintf("failed to to settings properties by name%s", name))
	}
	set := &schema.SettingBase{}
	err := mapstructure.Decode(n.GetProperties(name), set)
	if err != nil {
		return "", err
	}
	return set.DefaultValue, nil
}

func (n *BaseNode) DecodeProperties(name string, output interface{}) error {
	data := n.GetProperties(name)
	if data == nil {
		return errors.New(fmt.Sprintf("failed to to settings properties by name%s", name))
	}
	return mapstructure.Decode(n.GetProperties(name), output)
}

func (n *BaseNode) GetProperties(name string) interface{} {
	for _, setting := range n.Settings {
		if name == setting.Title {
			return setting.Properties
		}
	}
	return nil
}

func BuildSetting(propType, settingTitle string, body *BaseNode) (*Settings, error) {
	decode := schema.NewString(nil)
	err := body.DecodeProperties(settingTitle, decode)
	if err != nil {
		return nil, err
	}
	newSchema := schema.NewString(&schema.SettingBase{
		Title:        settingTitle,
		Min:          1,
		DefaultValue: decode.DefaultValue,
	})

	return &Settings{
		Type:       propType,
		Title:      settingTitle,
		Properties: newSchema,
	}, nil
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
