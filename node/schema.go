package node

// Schema is needed for the flow-ui
type Schema struct {
	Id       string      `json:"id"`                 // node uuid
	Type     string      `json:"type"`               // math/add
	Metadata *Metadata   `json:"metadata,omitempty"` // positions on the editor
	Inputs   interface{} `json:"inputs"`
}

type Inputs struct {
	Links map[string]Links `json:"links,omitempty"`
}

// Links node links
type Links struct {
	NodeId string      `json:"nodeId,omitempty"` // from node uuid
	Socket string      `json:"socket,omitempty"` // this is the port/pin name
	Value  interface{} `json:"value"`
}
