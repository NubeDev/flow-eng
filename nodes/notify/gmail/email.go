package gmail

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/notify"
	"github.com/jordan-wright/email"
	"github.com/mitchellh/mapstructure"
	"net/smtp"
	"strings"
)

const (
	gmailNode = "gmail"
)

type Gmail struct {
	*node.Spec
	firstLoop bool
	triggered bool
	loopCount uint64
	address   []string
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
	to := node.BuildInput(node.To, node.TypeString, nil, body.Inputs)
	subject := node.BuildInput(node.Subject, node.TypeString, nil, body.Inputs)
	message := node.BuildInput(node.Message, node.TypeString, nil, body.Inputs)
	trigger := node.BuildInput(node.TriggerInput, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(to, subject, message, trigger)
	outputs := node.BuildOutputs(node.BuildOutput(node.Result, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetSchema(buildSchema())
	return &Gmail{body, false, false, 0, address}, nil
}

func (inst *Gmail) setEmailClient() {
	connection, err := inst.GetDB().GetConnection("")
	if err != nil {
		inst.firstLoop = false // if fail try again
		return
	}
	inst.firstLoop = true
	pprint.Print(connection)
	e := email.NewEmail()
	e.From = "nubeio <noreply@nube-io.com>"
	e.To = []string{"ap@nube-io.com"}
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
	inst.loopCount++
	if !inst.firstLoop {
		//inst.setEmailClient()
	}

	// if trigger == true then set triggered to true
	// if trigger == false && triggered == true then rest triggered to false
	// now we can send gmail again
	trigger, _ := inst.ReadPinAsBool(node.TriggerInput)
	if trigger {
		fmt.Println("TRIGGER EMAIL")
		inst.triggered = true
	}
	if !trigger && inst.triggered {
		fmt.Println("RESET")
		inst.triggered = false
	}

}
