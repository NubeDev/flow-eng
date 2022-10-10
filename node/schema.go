package node

// Schema is needed for the flow-ui
type Schema struct {
	Id       string                  `json:"id"`                 // node uuid
	Type     string                  `json:"type"`               // math/add
	Metadata *Metadata               `json:"metadata,omitempty"` // positions on the editor
	Inputs   map[string]SchemaInputs `json:"inputs,omitempty"`
	Settings map[string]interface{}  `json:"settings,omitempty"`
}

type SchemaInputs struct {
	Value interface{}   `json:"value,omitempty"`
	Links []SchemaLinks `json:"links,omitempty"`
}

// SchemaLinks node links
type SchemaLinks struct {
	NodeId string `json:"nodeId"` // from node uuid
	Socket string `json:"socket"` // this is the port/pin name
}
