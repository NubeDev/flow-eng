package gmail

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/smtp"

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
	body.SetHelp(fmt.Sprintln("Please format the message in HTML. Instructions to generate gmail application token. \n 1. Go to your Google Account. \n 2. Select Security. \n 3. Under 'Signing in to Google,' select App Passwords. You may need to sign in. 4. At the bottom, choose Select app and choose the app you using and then Select device and choose the device you’re using and then Generate."))

	address, err := getSettings(body)
	if err != nil {
		return nil, err
	}

	body.SetSchema(buildSchema())
	return &Gmail{body}, nil
}

func (inst *Gmail) sendEmail() {
	to := inst.GetInput(node.To).GetValue()
	subject := inst.GetInput(node.Subject).GetValue()
	message := inst.GetInput(node.Message).GetValue()

	settingMap := inst.GetSettings()
	if settingMap == nil {
		return
	}

	e := email.NewEmail()
	e.From = settingMap["fromAddress"].(string)
	e.To = []string{to.(string)}
	e.Subject = subject.(string)
	e.HTML = []byte(message.(string))
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", settingMap["fromAddress"].(string), settingMap["token"].(string), "smtp.gmail.com"))
	if err != nil {
		log.Error(err)
		return
	}
}

// Process
// instructions to generate gmail application token
// 1. Go to your Google Account.
// 2. Select Security.
// 3. Under "Signing in to Google," select App Passwords. You may need to sign in.
// 4. At the bottom, choose Select app and choose the app you using and then Select device and choose the device you’re using and then Generate.
// refer to this page if unclear: https://support.google.com/accounts/answer/185833?visit_id=638089001343155683-2129707965&p=InvalidSecondFactor&rd=1
func (inst *Gmail) Process() {
	_, cov := inst.InputUpdated(node.TriggerInput)
	if cov {
		inst.sendEmail()
	}
}
