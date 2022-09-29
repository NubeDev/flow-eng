package email

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/notify"
	"github.com/mitchellh/mapstructure"
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
	body.Inputs = node.BuildInputs(to, subject, message, trigger)
	body.Outputs = node.BuildOutputs(node.BuildOutput(node.Result, node.TypeString, nil, body.Outputs))
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

}

func (inst *Gmail) Process() {
	inst.loopCount++
	if !inst.firstLoop {
		//inst.setEmailClient()
	}

	// if trigger == true then set triggered to true
	// if trigger == false && triggered == true then rest triggered to false
	// now we can send email again
	trigger := inst.ReadPinBool(node.TriggerInput)
	if trigger {
		fmt.Println("TRIGGER EMAIL")
		inst.triggered = true
	}
	if !trigger && inst.triggered {
		fmt.Println("RESET")
		inst.triggered = false
	}

}

func (inst *Gmail) Cleanup() {}
