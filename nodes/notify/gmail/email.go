package gmail

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/notify"
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
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

func (inst *Gmail) sendEmail() {
	connection, err := inst.GetDB().GetConnection("")
	pprint.Print(connection)
	e := email.NewEmail()

	// log.Println("The content of settings is: ", inst.Settings)

	e.From = "nubeio <noreply@nube-io.com>"
	e.To = []string{"jfe@nube-io.com"}
	e.Subject = "test"
	e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	err = e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "noreply@nube-io.com", "22222-11eb-111111111", "smtp.gmail.com"))
	fmt.Println(err)
	if err != nil {
		return
	}

}

func (inst *Gmail) Process() {
	_, cov := inst.InputUpdated(node.TriggerInput)
	s, _ := getSettings(inst.GetSettings())
	log.Println("The setting is: ", s)
	if cov {
		fmt.Println("TRIGGER EMAIL")
		inst.sendEmail()
	}

}
