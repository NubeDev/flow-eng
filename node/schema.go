package node

// Schema is needed for the flow-ui
type Schema struct {
	Id       string                  `json:"id"`                 // node uuid
	Type     string                  `json:"type"`               // math/add
	Metadata *Metadata               `json:"metadata,omitempty"` // positions on the editor
	Inputs   map[string]SchemaInputs `json:"inputs"`
}

type SchemaInputs struct {
	Value interface{} `json:"value,omitempty"`
	Links []struct {
		NodeId string `json:"nodeId"`
		Socket string `json:"socket"`
	} `json:"links,omitempty"`
}

// SchemaLinks node links
type SchemaLinks struct {
	NodeId string      `json:"nodeId,omitempty"` // from node uuid
	Socket OutputName  `json:"socket,omitempty"` // this is the port/pin name
	Value  interface{} `json:"value"`
}
