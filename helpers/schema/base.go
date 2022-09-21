package schema

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
)

type UiOptions struct {
	Inline bool `json:"inline"`
}

type radio struct {
	UiWidget  string    `json:"ui:widget"`
	UiOptions UiOptions `json:"ui:options"`
}

type setting struct {
	Type       string      `json:"type"`
	Title      string      `json:"title"`
	Properties interface{} `json:"properties"`
}

func test() {
	p := radio{
		UiWidget: "radio",
		UiOptions: UiOptions{
			Inline: true,
		},
	}
	s := setting{
		Type:       "aaaa",
		Title:      "addfsdf",
		Properties: p,
	}
	pprint.PrintJOSN(s)

}
