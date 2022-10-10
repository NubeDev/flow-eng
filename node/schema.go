package node

// Schema is needed for the flow-ui
type Schema struct {
	Id       string                  `json:"id"`                 // node uuid
	Type     string                  `json:"type"`               // math/add
	Metadata *Metadata               `json:"metadata,omitempty"` // positions on the editor
	Inputs   map[string]SchemaInputs `json:"inputs"`
	Settings map[string]interface{}  `json:"settings"`
}

type SchemaInputs struct {
	Value interface{}   `json:"value,omitempty"`
	Links []SchemaLinks `json:"links"`
}

// SchemaLinks node links
type SchemaLinks struct {
	NodeId string `json:"nodeId,omitempty"` // from node uuid
	Socket string `json:"socket,omitempty"` // this is the port/pin name
}
