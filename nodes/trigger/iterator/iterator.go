package iterator

import (
	"fmt"
	"time"

	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/trigger"
)

type Iterator struct {
	*node.Spec
	c                  chan int
	s                  map[string]interface{}
	iterationCompleted float64
	running            bool
	instructedPause    bool
}

const (
	Terminate = 0
	Pause     = 1
	Run       = 2
)

const (
	Milliseconds = "milliseconds"
	Seconds      = "seconds"
	Minutes      = "minutes"
	Hours        = "hours"
)

func NewIterator(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, trigger.Iterator, trigger.Category)

	interval := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs, nil)
	iterations := node.BuildInput(node.Iterations, node.TypeFloat, nil, body.Inputs, nil)
	start := node.BuildInput(node.Start, node.TypeBool, nil, body.Inputs, nil)
	stop := node.BuildInput(node.Stop, node.TypeBool, nil, body.Inputs, nil)
	inputs := node.BuildInputs(interval, iterations, start, stop)

	out := node.BuildOutput(node.Outp, node.TypeBool, nil, body.Outputs)
	complete := node.BuildOutput(node.Complete, node.TypeBool, nil, body.Outputs)
	count := node.BuildOutput(node.CountOut, node.TypeFloat, nil, body.Outputs)
	outputs := node.BuildOutputs(out, complete, count)

	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body.SetSchema(buildSchema())
	s := body.GetSettings()
	body.SetHelp("This node generates a sequence of 'false' to 'true' transitions on 'output'.  The number of 'false' to 'true' transitions will be equal to 'count' value (or 'Iterations' setting);  these values are sent over the 'interval' duration (unless interrupted by 'stop' input).  For example, if 'interval' is set to 5 (seconds) and 'iterations' is set to 5, a 'false' to 'true' transition will occur on 'output' every 1 second.  If 'stop' input is 'true' then the next 'true' value will not be sent from 'output' until 'stop' is 'false' again. 'interval' units can be configured from settings. Maximum 'interval' setting is 587 hours.")

	return &Iterator{body, nil, s, 0, false, false}, nil
}

func (inst *Iterator) Process() {
	interval, intervalNull := inst.ReadPinAsFloat(node.Interval)
	iterations, iterationsNull := inst.ReadPinAsFloat(node.Iterations)
	_, startBool := inst.InputUpdated(node.Start)
	stop, _ := inst.ReadPinAsBool(node.Stop)

	// fall back to values set in setting if input is not connected
	if intervalNull && inst.s["interval"] != nil {
		interval = inst.s["interval"].(float64)
	}
	if iterationsNull && inst.s["iterations"] != nil {
		iterations = inst.s["iterations"].(float64)
	}
	var units interface{}
	if inst.s["units"] == nil {
		units = "seconds"
	} else {
		units = inst.s["units"]
	}

	if startBool && !inst.running {
		// output false to complete on iteration start
		inst.WritePinBool(node.Complete, false)
		// calculate period, output default period if any of the inputs is nil
		// use the input period otherwise
		var period float64
		if intervalNull || iterationsNull {
			period = inst.s["interval"].(float64) / inst.s["iterations"].(float64)
		} else {
			period = interval / iterations
		}
		inst.c = make(chan int, 1)
		go iterate(inst, inst.c, period, iterations, units)
		inst.running = true
	}

	// following block only executes when a iteration is running
	// iteration halted
	if inst.running && stop {
		// give the pause signal to iterating routine only once when 'stop' becomes true
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

}

func iterate(inst *Iterator, c chan int, period float64, iterations float64, units interface{}) {
	// set state to 'run' when iteration starts
	state := Run
	var duration time.Duration
	halfPeriod := period / 2
	switch units.(string) {
	case string(Milliseconds):
		duration = time.Duration(halfPeriod * float64(time.Millisecond))
	case string(Seconds):
		duration = time.Duration(halfPeriod * float64(time.Second))
	case string(Minutes):
		duration = time.Duration(halfPeriod * float64(time.Minute))
	case string(Hours):
		duration = time.Duration(halfPeriod * float64(time.Hour))
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
					fmt.Println("iterating...")
				case Pause:
					fmt.Println("paused...")
				case Terminate:
					fmt.Println("terminated...")
					return
				}
			// check state at the beginning of each loop, break if state is 'Pause'
			default:
				if state == Pause {
					break
				}
				// write the current iteration number, starting from 0
				inst.WritePinFloat(node.CountOut, inst.iterationCompleted)
				inst.iterationCompleted++
				// write out the waveform, 'true' for first half period, and 'false' for the second half
				// this arrangement allows the program to stop on false when stop become true
				inst.WritePinBool(node.Outp, true)
				time.Sleep(duration)
				inst.WritePinBool(node.Outp, false)
				time.Sleep(duration)
			}
		}
	}
}
