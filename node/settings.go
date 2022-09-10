package node

import (
	"errors"
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

func (n *BaseNode) GetProperties(name string) interface{} {
	for _, setting := range n.Settings {
		if name == setting.Title {
			return setting.Properties
		}
	}
	return nil
}

func BuildSetting(propType, title string, schema interface{}) *Settings {
	return &Settings{
		Type:       propType,
		Title:      title,
		Properties: schema,
	}

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
