package node

type Status struct {
	Message      string `json:"message,omitempty"`
	InError      bool   `json:"inError,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func (n *Spec) SetStatus(body *Status) {

}
func (n *Spec) GetStatus() *Status {

	return nil
}
