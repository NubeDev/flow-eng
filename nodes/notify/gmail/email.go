package gmail

import (
	"fmt"
	"net/smtp"

	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/notify"
	"github.com/jordan-wright/email"
)

const (
	gmailNode = "gmail"
)

type Gmail struct {
	*node.Spec
}

func NewGmail(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, gmailNode, notify.Category)
	to := node.BuildInput(node.To, node.TypeString, nil, body.Inputs)
	subject := node.BuildInput(node.Subject, node.TypeString, nil, body.Inputs)
	message := node.BuildInput(node.Message, node.TypeString, nil, body.Inputs)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs)
	inputs := node.BuildInputs(to, subject, message, trigger)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	return &Gmail{body}, nil
}

func (inst *Gmail) sendEmail(ed map[string]string) {
	connection, err := inst.GetDB().GetConnection("")
	pprint.Print(connection)
	e := email.NewEmail()

	e.From = ed["from"]
	e.To = []string{ed["to"]}
	e.Subject = ed["subject"]
	e.Text = []byte(ed["message"])
	// e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	err = e.Send("smtp.gmail.com:587", smtp.PlainAuth("", ed["from"], ed["password"], "smtp.gmail.com"))
	fmt.Println(err)
	if err != nil {
		return
	}

}

// instructions to generate gmail application password
// 1. Go to your Google Account.
// 2. Select Security.
// 3. Under "Signing in to Google," select App Passwords. You may need to sign in.
// 4. At the bottom, choose Select app and choose the app you using and then Select device and choose the device you’re using and then Generate.
// refer to this page if unclear: https://support.google.com/accounts/answer/185833?visit_id=638089001343155683-2129707965&p=InvalidSecondFactor&rd=1
func (inst *Gmail) Process() {
	to := inst.GetInput(node.To).GetValue()
	subject := inst.GetInput(node.Subject).GetValue()
	message := inst.GetInput(node.Message).GetValue()

	s, _ := getSettings(inst.GetSettings())

	var ed map[string]string = make(map[string]string)
	ed["from"] = s.FromAddress
	ed["password"] = s.Password
	ed["subject"] = subject.(string)
	ed["message"] = message.(string)
	if to != nil {
		ed["to"] = to.(string)
	} else {
		ed["to"] = s.ToAddress
	}

	_, cov := inst.InputUpdated(node.TriggerInput)
	if cov {
		fmt.Println("TRIGGER EMAIL")
		inst.sendEmail(ed)
	}

}
