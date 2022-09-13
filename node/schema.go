package node

// Schema is needed for the flow-ui
type Schema struct {
	Id       string     `json:"id"`                 // node uuid
	Type     string     `json:"type"`               // math/add
	Metadata *Metadata  `json:"metadata,omitempty"` // positions on the editor
	Inputs   *InputsMap `json:"inputs,omitempty"`
}

// Links node links
type Links struct {
	NodeId string      `json:"nodeId,omitempty"` // from node uuid
	Socket string      `json:"socket,omitempty"` // this is the port/pin name
	Value  interface{} `json:"value,omitempty"`  // when there is no link but a value is set by the user
}

type InputsMap struct {
	Links map[string]Links `json:"links,omitempty"`
}
