package notify

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/enescakir/emoji"
	"github.com/jordan-wright/email"
	log "github.com/sirupsen/logrus"
	"net/smtp"
	"strings"
)

type Gmail struct {
	*node.Spec
	lastFromAddress string
}

func NewGmail(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, gmail, category)

	to := node.BuildInput(node.To, node.TypeString, nil, body.Inputs, true, false)
	subject := node.BuildInput(node.Subject, node.TypeString, nil, body.Inputs, true, false)
	message := node.BuildInput(node.Message, node.TypeString, nil, body.Inputs, true, false)
	trigger := node.BuildInput(node.TriggerInput, node.TypeBool, nil, body.Inputs, false, true)
	inputs := node.BuildInputs(to, subject, message, trigger)

	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &Gmail{body, ""}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *Gmail) sendEmail() {
	to := inst.ReadPinOrSettingsString(node.To)
	// fmt.Println(fmt.Sprintf("sendEmail() to: %+v", to))
	toArray := strings.Split(to, ",")
	for i, _ := range toArray {
		toArray[i] = strings.Trim(toArray[i], " ")
	}
	// fmt.Println(fmt.Sprintf("sendEmail() toArray: %+v", toArray))
	subject := inst.ReadPinOrSettingsString(node.Subject)
	message := inst.ReadPinOrSettingsString(node.Message)
	settingMap := inst.GetSettings()
	if settingMap == nil {
		return
	}
	e := email.NewEmail()
	e.From = settingMap["sender-address"].(string)
	if e.From != inst.lastFromAddress {
		inst.SetSubTitle(e.From)
		inst.lastFromAddress = e.From
	}
	e.To = toArray
	e.Subject = subject
	// e.Text = []byte(ed["message"])
	e.HTML = []byte(message)
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", settingMap["sender-address"].(string), settingMap["password"].(string), "smtp.gmail.com"))
	if err != nil {
		log.Error(err)
		inst.SetWaringMessage(err.Error())
		inst.SetWaringIcon(string(emoji.RedCircle))
		return
	} else {
		inst.SetWaringMessage("")
		inst.SetWaringIcon(string(emoji.GreenCircle))
	}
}

func (inst *Gmail) Process() {
	_, cov := inst.InputUpdated(node.TriggerInput)
	if cov {
		inst.sendEmail()
	}
	inst.WritePin(node.Out, "")
}

// Custom Node Settings Schema

type GmailSettingsSchema struct {
	SenderAddress    schemas.String `json:"sender-address"`
	Password         schemas.String `json:"password"`
	RecipientAddress schemas.String `json:"to"`
	Subject          schemas.String `json:"subject"`
	Body             schemas.String `json:"message"`
}

type GmailSettings struct {
	SenderAddress    string `json:"sender-address"`
	Password         string `json:"password"`
	RecipientAddress string `json:"to"`
	Subject          string `json:"subject"`
	Body             string `json:"message"`
}

func (inst *Gmail) buildSchema() *schemas.Schema {
	props := &GmailSettingsSchema{}

	// sender email address
	props.SenderAddress.Title = "Sender Email Address"
	props.SenderAddress.Default = ""

	// sender password
	props.Password.Title = "Sender Email Password"
	props.Password.Default = ""

	// recipient email address
	props.RecipientAddress.Title = "Recipient Email Address"
	props.RecipientAddress.Default = ""

	// email subject
	props.Subject.Title = "Email Subject"
	props.Subject.Default = ""

	// email body
	props.Body.Title = "Email Body"
	props.Body.Default = ""

	schema.Set(props)

	uiSchema := array.Map{
		"ui:order": array.Slice{"sender-address", "password", "to", "subject", "message"},
		"password": array.Map{
			"ui:widget": "password",
		},
		"sender-address": array.Map{
			"ui:widget": "email",
		},
		"to": array.Map{
			"ui:widget": "email",
		},
	}
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "Node Settings",
			Properties: props,
		},
		UiSchema: uiSchema,
	}
	return s
}

func (inst *Gmail) getSettings(body map[string]interface{}) (*GmailSettings, error) {
	settings := &GmailSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
