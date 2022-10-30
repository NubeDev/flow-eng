package node

import "fmt"

type Status struct {
	SubTitle      string `json:"subTitle,omitempty"`
	InError       bool   `json:"inError,omitempty"`
	Message       string `json:"message,omitempty"`
	ErrorMessage  string `json:"errorMessage,omitempty"`
	WaringMessage string `json:"waringMessage,omitempty"`
	NotifyMessage string `json:"notifyMessage,omitempty"`
}

func (n *Spec) SetSubTitle(message string) {
	if n.Status == nil {
		n.Status = &Status{}
	}
	message = fmt.Sprintf("%s %s", n.Info.Category, message)
	n.Status.SubTitle = message
}

func (n *Spec) SetStatusError(message string) {
	if n.Status == nil {
		n.Status = &Status{}
	}
	n.Status.InError = true
	n.Status.ErrorMessage = message
}

func (n *Spec) SetNotifyMessage(message string) {
	if n.Status == nil {
		n.Status = &Status{}
	}
	n.Status.NotifyMessage = message
}

func (n *Spec) SetWaringMessage(message string) {
	if n.Status == nil {
		n.Status = &Status{}
	}
	n.Status.WaringMessage = message
}

func (n *Spec) SetStatusMessage(message string) {
	if n.Status == nil {
		n.Status = &Status{}
	}
	n.Status.Message = message
}

func (n *Spec) SetStatus(body *Status) {
	n.Status = body
}
func (n *Spec) GetStatus() *Status {
	return n.Status
}
