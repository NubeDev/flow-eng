package node

type Status struct {
	SubTitle      string `json:"subTitle,omitempty"`
	ActiveMessage bool   `json:"activeMessage,omitempty"`
	Message       string `json:"message,omitempty"`
	ErrorMessage  string `json:"errorMessage,omitempty"`
	ErrorIcon     string `json:"errorIcon,omitempty"`
	WaringMessage string `json:"waringMessage,omitempty"`
	WaringIcon    string `json:"waringIcon,omitempty"`
	NotifyMessage string `json:"notifyMessage,omitempty"`
	NotifyIcon    string `json:"notifyIcon,omitempty"`
}

func (n *Spec) SetErrorIcon(icon string) {
	if n.Status == nil {
		n.Status = &Status{}
	}
	n.Status.ActiveMessage = true
	n.Status.ErrorIcon = icon
}

func (n *Spec) SetNotifyIcon(icon string) {
	if n.Status == nil {
		n.Status = &Status{}
	}
	n.Status.ActiveMessage = true
	n.Status.NotifyIcon = icon
}

func (n *Spec) SetWaringIcon(icon string) {
	if n.Status == nil {
		n.Status = &Status{}
	}
	n.Status.ActiveMessage = true
	n.Status.WaringIcon = icon
}

func (n *Spec) SetSubTitle(message string) {
	if n.Status == nil {
		n.Status = &Status{}
	}
	if n.GetNodeName() == "" && n.Info.Category != "bacnet" {
		n.SetNodeName(n.GetName())
	}
	n.Status.SubTitle = message
}

func (n *Spec) SetStatusError(message string) {
	if n.Status == nil {
		n.Status = &Status{}
	}
	n.Status.ActiveMessage = true
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
	n.Status.ActiveMessage = true
	n.Status.WaringMessage = message
}

func (n *Spec) SetStatusMessage(message string) {
	if n.Status == nil {
		n.Status = &Status{}
	}
	n.Status.ActiveMessage = true
	n.Status.Message = message
}

func (n *Spec) SetStatus(body *Status) {
	n.Status = body
}
func (n *Spec) GetStatus() *Status {
	return n.Status
}
