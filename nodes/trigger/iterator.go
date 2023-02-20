package trigger

import (
	"encoding/json"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
	"time"

	"github.com/NubeDev/flow-eng/node"
)

type Iterate struct {
	*node.Spec
	c                  chan int
	iterationCompleted float64
	running            bool
	instructedPause    bool
	lastStart          bool
	lastInterval       time.Duration
	lastIterations     int
	currentOutput      bool
}

const (
	Terminate = 0
	Pause     = 1
	Run       = 2
)

func NewIterate(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, Iterator, Category)

	interval := node.BuildInput(node.Interval, node.TypeFloat, 10, body.Inputs, true, false)
	iterations := node.BuildInput(node.Iterations, node.TypeFloat, 5, body.Inputs, true, false)
	start := node.BuildInput(node.Start, node.TypeBool, nil, body.Inputs, false, false)
	stop := node.BuildInput(node.Stop, node.TypeBool, nil, body.Inputs, false, false)
	inputs := node.BuildInputs(interval, iterations, start, stop)

	out := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)
	complete := node.BuildOutput(node.Complete, node.TypeBool, nil, body.Outputs)
	count := node.BuildOutput(node.CountOut, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out, complete, count)

	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetHelp("This node generates a sequence of 'false' to 'true' transitions on 'output'.  The number of 'false' to 'true' transitions will be equal to 'count' value (or 'Iterations' setting);  these values are sent over the 'interval' duration (unless interrupted by 'stop' input).  For example, if 'interval' is set to 5 (seconds) and 'iterations' is set to 5, a 'false' to 'true' transition will occur on 'output' every 1 second.  If 'stop' input is 'true' then the next 'true' value will not be sent from 'output' until 'stop' is 'false' again. 'interval' units can be configured from settings. Maximum 'interval' setting is 587 hours.")

	n := &Iterate{body, nil, 0, false, false, true, -1, -1, false}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *Iterate) Process() {
	intervalDuration, _ := inst.ReadPinAsTimeSettings(node.Interval)
	iterations := inst.ReadPinOrSettingsFloat(node.Iterations)
	iterations = math.Floor(iterations)

	if intervalDuration != inst.lastInterval || iterations != float64(inst.lastIterations) {
		inst.setSubtitle(intervalDuration, int(iterations))
		inst.lastInterval = intervalDuration
		inst.lastIterations = int(iterations)
	}

	start, _ := inst.ReadPinAsBool(node.Start)
	stop, _ := inst.ReadPinAsBool(node.Stop)

	if start && !inst.lastStart && !inst.running && iterations > 0 {
		// output false to complete on iteration start
		inst.WritePinBool(node.Complete, false)
		// calculate period
		period := time.Duration(int64(intervalDuration) / int64(iterations))
		inst.c = make(chan int, 1)
		go iterate(inst, inst.c, period, iterations, stop)
		inst.running = true
	}
	inst.lastStart = start

	// following block only executes when a iteration is running
	// iteration halted
	if inst.running && stop {
		// give the pause signal to iterating routine only once when 'stop' becomes true
		// TODO: why is this extra function wrapper necessary?
		if !inst.instructedPause {
			go func(c chan int) {
				c <- Pause
			}(inst.c)
			// set 'instructedPause' to 'true' to prevent more pause signals send over the channel
			inst.instructedPause = true
		}
		// iteration continues
	} else if inst.running && !stop {
		// give the run signal over the channel only if the iteration has been paused
		if inst.instructedPause {
			go func(c chan int) {
				c <- Run
			}(inst.c)
			// reset 'instructedPause' after giving the 'run' signal
			inst.instructedPause = false
		}
	}
	complete := inst.iterationCompleted == iterations
	inst.WritePinBool(node.Complete, complete)
	inst.WritePinFloat(node.CountOut, inst.iterationCompleted)
	if !inst.running {
		inst.WritePinBool(node.Out, false)
	} else {
		inst.WritePinBool(node.Out, inst.currentOutput)
	}
}

func iterate(inst *Iterate, c chan int, duration time.Duration, iterations float64, startOnPause bool) {
	// set state to 'run' when iteration starts
	state := Run
	if startOnPause {
		state = Pause
		inst.WritePinFloat(node.CountOut, 0)
		inst.WritePinBool(node.Out, false)
		inst.currentOutput = false
	}
	for {
		// check for terminal condition
		if iterations-inst.iterationCompleted == 0 {
			inst.WritePinBool(node.Complete, true)
			inst.iterationCompleted = 0
			inst.running = false
			close(c)
			return
		}
		// start iterating
		for i := 0; i <= int(iterations-inst.iterationCompleted); i++ {
			// start iterating if no message received over the channel
			select {
			case state = <-c:
				switch state {
				case Run:
					log.Info("Iterator iterating...")
				case Pause:
					log.Info("Iterator paused...")
				case Terminate:
					log.Info("Iterator terminated...")
					return
				}
			// check state at the beginning of each loop, break if state is 'Pause'
			default:
				if state == Pause {
					break
				}
				// write the current iteration number, starting from 1
				inst.iterationCompleted++
				inst.WritePinFloat(node.CountOut, inst.iterationCompleted)
				// write out the waveform, 'true' for first half period, and 'false' for the second half
				// this arrangement allows the program to stop on false when stop become true
				inst.WritePinBool(node.Out, true)
				inst.currentOutput = true
				time.Sleep(duration / 2)
				inst.WritePinBool(node.Out, false)
				inst.currentOutput = false
				time.Sleep(duration / 2)
			}
		}
	}
}

func (inst *Iterate) setSubtitle(intervalDuration time.Duration, iterations int) {
	subtitleText := strconv.Itoa(iterations)
	subtitleText += " iterations over "
	subtitleText += intervalDuration.String()
	inst.SetSubTitle(subtitleText)
}

// Custom Node Settings Schema

type IterateSettingsSchema struct {
	Interval          schemas.Number     `json:"interval"`
	IntervalTimeUnits schemas.EnumString `json:"interval_time_units"`
	Iterations        schemas.Integer    `json:"iterations"`
}

type IterateSettings struct {
	Interval          float64 `json:"interval"`
	IntervalTimeUnits string  `json:"interval_time_units"`
	Iterations        int     `json:"iterations"`
}

func (inst *Iterate) buildSchema() *schemas.Schema {
	props := &IterateSettingsSchema{}

	// time selection
	props.Interval.Title = "Interval"
	props.Interval.Default = 10

	props.IntervalTimeUnits.Title = "Interval Units"
	props.IntervalTimeUnits.Default = ttime.Sec
	props.IntervalTimeUnits.Options = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}
	props.IntervalTimeUnits.EnumName = []string{ttime.Ms, ttime.Sec, ttime.Min, ttime.Hr}

	// iterations
	props.Iterations.Title = "Iterations"
	props.Iterations.Default = 5

	schema.Set(props)

	uiSchema := array.Map{
		"interval_time_units": array.Map{
			"ui:widget": "radio",
			"ui:options": array.Map{
				"inline": true,
			},
		},
		"ui:order": array.Slice{"interval", "interval_time_units", "iterations"},
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

func (inst *Iterate) getSettings(body map[string]interface{}) (*IterateSettings, error) {
	settings := &IterateSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
