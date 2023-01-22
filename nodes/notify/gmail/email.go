package gmail

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/notify"
	"github.com/jordan-wright/email"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"net/smtp"
	"strings"
)

const (
	gmailNode = "gmail"
)

type Gmail struct {
	*node.Spec
	address []string
}

type nodeSettings struct {
	Address string `json:"address"`
}

func getSettings(body *node.Spec) ([]string, error) {
	settings := &nodeSettings{}
	err := mapstructure.Decode(body.Settings, &settings)
	if err != nil {
		return nil, err
	}
	if settings != nil {
		address := strings.Split(settings.Address, ",")
		return address, nil
	}
	return nil, nil

}

func NewGmail(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, gmailNode, notify.Category)
	address, err := getSettings(body)
	if err != nil {
		return nil, err
	}
	to := node.BuildInput(node.To, node.TypeString, nil, body.Inputs, nil)
	subject := node.BuildInput(node.Subject, node.TypeString, nil, body.Inputs, nil)
	message := node.BuildInput(node.Message, node.TypeString, nil, body.Inputs, nil)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs, nil)
	inputs := node.BuildInputs(to, subject, message, trigger)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	return &Gmail{body, address}, nil
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
	// e.Text = []byte(ed["message"])
	e.HTML = []byte(message.(string))
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", settingMap["fromAddress"].(string), settingMap["token"].(string), "smtp.gmail.com"))
	if err != nil {
		log.Error(err)
		return
	}

}

func (inst *Gmail) Process() {
	_, cov := inst.InputUpdated(node.TriggerInput)
	if cov {
		inst.sendEmail()
	}

}
