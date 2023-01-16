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

func NewIterator(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, trigger.Iterator, trigger.Category)

	interval := node.BuildInput(node.Interval, node.TypeFloat, nil, body.Inputs)
	iterations := node.BuildInput(node.Iterations, node.TypeFloat, nil, body.Inputs)
	start := node.BuildInput(node.Start, node.TypeBool, nil, body.Inputs)
	stop := node.BuildInput(node.Stop, node.TypeBool, nil, body.Inputs)
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
	// fmt.Println("stop is: ", stop)

	//fall back to values set in setting if input is not connected
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
		// calculate period
		period := interval / iterations
		fmt.Println("the period is: ", period)
		inst.c = make(chan int, 1)
		go iterate(inst, inst.c, period, iterations, units)
		inst.running = true
	}

	// iteration halted
	if inst.running && stop {
		if !inst.instructedPause {
			fmt.Println("entered pause outside!!!!!!!!!!!!!!")
			go func(c chan int) {
				c <- Pause
				fmt.Println("entered pause inside!!!!!!!!!!!!!!")
			}(inst.c)
			inst.instructedPause = true
		}
	} else if inst.running && !stop {
		if inst.instructedPause {
			fmt.Println("stopped pause outside!!!!!!!!!!!!!!")
			go func(c chan int) {
				c <- Run
				fmt.Println("stopped pause inside!!!!!!!!!!!!!!")
			}(inst.c)
			inst.instructedPause = false
		}
	}

}

func iterate(inst *Iterator, c chan int, period float64, iterations float64, units interface{}) {
	state := Run
	var duration time.Duration
	switch units.(string) {
	case string(trigger.Milliseconds):
		duration = time.Duration(period/2) * time.Millisecond
	case string(trigger.Seconds):
		duration = time.Duration(period/2) * time.Second
	case string(trigger.Minutes):
		duration = time.Duration(period/2) * time.Minute
	case string(trigger.Hours):
		duration = time.Duration(period/2) * time.Hour
	}

	for {
		// check for terminal condition
		if iterations-inst.iterationCompleted == 0 {
			inst.WritePinBool(node.Complete, true)
			inst.iterationCompleted = 0
			inst.running = false
			return
		}
		// start iterating
		for i := 0; i <= int(iterations-inst.iterationCompleted); i++ {
			select {
			case state = <-c:
				fmt.Println("state is: ", state)
				switch state {
				case Run:
					fmt.Println("iterating...")
				case Pause:
					fmt.Println("paused...")
				case Terminate:
					fmt.Println("terminated...")
					return
				}
			default:
				if state == Pause {
					break
				}
				inst.WritePinFloat(node.CountOut, inst.iterationCompleted)
				inst.iterationCompleted++
				inst.WritePinBool(node.Outp, false)
				time.Sleep(duration)
				inst.WritePinBool(node.Outp, true)
				time.Sleep(duration)
			}
		}
	}
}
