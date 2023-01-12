package streams

import (
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type Flatline struct {
	*node.Spec
	timeout     *time.Timer
	lastVal     float64 // TODO: input value should be allowed to be nil
	alertStatus bool
}

func NewFlatline(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flatline, category)

	in := node.BuildInput(node.In, node.TypeFloat, nil, body.Inputs, nil) // TODO: this input shouldn't have a manual override value
	inputs := node.BuildInputs(in)
	outputs := node.BuildOutputs(node.BuildOutput(node.FlatLine, node.TypeFloat, nil, body.Outputs))
	/*
		// TODO: alert delay value should be set by input value OR by fallback to settings value
		_, setting, _, err := node.NewSetting(body, &node.SettingOptions{Type: node.Number, Title: node.AlertDelayMins, Min: 1, Max: 50000, Value: 1})
		if err != nil {
			return nil, err
		}
		settings, err := node.BuildSettings(setting)
		if err != nil {
			return nil, err
		}
	*/
	// body = node.BuildNode(body, inputs, outputs, settings)
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	return &Flatline{body, nil, 0, false}, nil
}

func (inst *Flatline) Process() {
	in, _ := inst.ReadPinAsFloat(node.In) // TODO: input value should be allowed to be nil
	if in != inst.lastVal {               // COV
		inst.lastVal = in
		inst.WritePin(node.FlatLine, 0)
		inst.alertStatus = false
		// create timeout function with specified delay
		f := func() {
			inst.WritePin(node.FlatLine, 1)
			inst.alertStatus = true
		}
		/*
			delayValue := inst.GetPropValueInt(node.AlertDelayMins, 30)
			alertDelayDuration, _ := time.ParseDuration(fmt.Sprintf("%dm", delayValue)) // TODO: value should come from input, or from settings value as a fallback.
			if alertDelayDuration <= 1*time.Second {
				alertDelayDuration = 1 * time.Minute
			}
			inst.timeout = time.AfterFunc(alertDelayDuration, f)
		*/
		inst.timeout = time.AfterFunc(30*time.Minute, f)
	}
	inst.WritePin(node.FlatLine, inst.alertStatus)
}
