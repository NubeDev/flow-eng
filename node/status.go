package node

type Status struct {
	Message      string `json:"message,omitempty"`
	InError      bool   `json:"inError,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func (n *Spec) SetStatus(body *Status) {
	n.Status = body
}
func (n *Spec) GetStatus() *Status {
	return n.Status
}
