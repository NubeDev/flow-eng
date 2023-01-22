package node

// Schema is needed for the flow-ui
type Schema struct {
	Id       string                   `json:"id"`   // node uuid
	Type     string                   `json:"type"` // math/add
	NodeName string                   `json:"nodeName,omitempty"`
	Icon     string                   `json:"icon,omitempty"`
	Metadata *Metadata                `json:"metadata,omitempty"` // positions on the editor
	Inputs   map[string]SchemaInputs  `json:"inputs"`
	Outputs  map[string]SchemaOutputs `json:"outputs"`
	Settings map[string]interface{}   `json:"settings,omitempty"`
	IsParent bool                     `json:"isParent"`
	ParentId string                   `json:"parentId,omitempty"`
	Payload  *Payload                 `json:"payload,omitempty"`
}

type SchemaOutputs struct {
	Position         int  `json:"position"`
	OverridePosition bool `json:"overridePosition"`
}

type SchemaInputs struct {
	Value            interface{}   `json:"value,omitempty"`
	Links            []SchemaLinks `json:"links,omitempty"`
	Position         int           `json:"position"`
	OverridePosition bool          `json:"overridePosition"`
}

// SchemaLinks node links
type SchemaLinks struct {
	NodeId string `json:"nodeId"` // from node uuid
	Socket string `json:"socket"` // this is the port/pin name
}
