package global

const DebugConnections bool = false

type Debug struct {
	NodeUUID    string `json:"nodeUUID"`
	NodeName    string `json:"nodeName"`
	FromOutput  string `json:"fromOutput"`
	ToInput     string `json:"toInput"`
	OutputValue any    `json:"outputValue"`
}
