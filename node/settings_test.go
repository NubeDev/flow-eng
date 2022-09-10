package node

import (
	"fmt"
	"testing"
)

type topic struct {
	Type     string `json:"type" default:"string"`
	Title    string `json:"title" default:"topic"`
	Min      int    `json:"minLength" default:"1"`
	ReadOnly bool   `json:"readOnly" default:"false"`
}

func TestBaseNode_GetSetting(t *testing.T) {
	body := &BaseNode{}
	body.Info.Name = "test"
	body.Info.Category = "test"
	body.Info.NodeID = "test"

	settings, _ := BuildSettings(BuildSetting("string", "topic", &topic{
		Type:     "string",
		Title:    "topic",
		Min:      1,
		ReadOnly: true,
	}))

	body.Settings = settings

	top := &topic{
		Type:     "string",
		Title:    "topic",
		Min:      1,
		ReadOnly: true,
	}
	fmt.Println(top)

	top_, ok := body.GetProperties("topi").(*topic)
	if ok {
		fmt.Println(top_.Title)
	}

}
