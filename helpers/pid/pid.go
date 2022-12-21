package pid

import (
	"errors"
	"math"
	"time"
)

/**
 * pid-controller -  A node advanced PID controller based on the Arduino PID library
 * github@wilberforce.co.nz Rhys Williams
 * Based on:
 * Arduino PID Library - Version 1.0.1
 * by Brett Beauregard <br3ttb@gmail.com> brettbeauregard.com
 *
 * This Library is licensed under a GPL-3.0 License
 */

type PID_DIRECTION bool

const (
	DIRECT  PID_DIRECTION = false
	REVERSE PID_DIRECTION = true
)

type PID_MODE bool

const (
	MANUAL PID_MODE = false
	AUTO   PID_MODE = true
)

type Pid struct {
	input          float64
	setpoint       float64
	output         float64
	currMode       PID_MODE
	direction      PID_DIRECTION
	bias           float64
	intervalMillis float64
	displayP       float64
	displayI       float64
	displayD       float64
	lastTime       float64
	kp             float64
	ki             float64
	kd             float64
	outputSum      float64
	lastInput      float64
	iTerm          float64
	outMin         float64
	outMax         float64
	inAuto         bool
}

func NewPid(input, setpoint, p, i, d, intervalSecs float64, dir PID_DIRECTION) *Pid {
	if p < 0 {
		p = 0
	}
	if i < 0 {
		i = 0
	}
	if d < 0 {
		d = 0
	}
	interval := intervalSecs * 1000
	last := float64(time.Now().UnixMilli()) - (interval)

	newPid := &Pid{input, setpoint, 0, MANUAL, dir, 0, interval, p, i, d, last, 0, 0, 0, 0, input, 0, 0, 100, false}
	newPid.Compute()

	return newPid
}

/**
 * Compute()
 * This, as they say, is where the magic happens.  pid function should be called
 * every time "void loop()" executes.  the function will decide for itself whether a new
 * pid Output needs to be computed.  returns true when the output is computed,
 * false when nothing has been done.
 */

func (p *Pid) Compute() (computed bool, err error) {
	if !p.inAuto {
		return false, errors.New("pid controller is not enabled")
	}

	now := float64(time.Now().UnixMilli())
	timeChange := now - p.lastTime
	if timeChange >= p.intervalMillis {
		// Compute all the working error variables
		input := p.input
		// var error = p.mySetpoint - input
		errorAmount := input - p.setpoint // above setpoint = positive error
		directionMultiplier := float64(1)
		if p.direction {
			directionMultiplier = -1
		}
		errorAmount = errorAmount * directionMultiplier

		p.iTerm += p.ki * errorAmount
		if p.iTerm > p.outMax-p.bias {
			p.iTerm = p.outMax - p.bias
		} else if p.iTerm < p.outMin-p.bias {
			p.iTerm = p.outMin - p.bias
		}

		dInput := input - p.lastInput

		// Compute PID Output
		// var output = ((p.kp * error) + p.iTerm - (p.kd * dInput)) * p.setDirection
		output := p.kp*errorAmount + p.iTerm - p.kd*dInput + p.bias
		// var output = ((p.kp * error) + p.iTerm - (p.kd * dInput))

		if output > p.outMax {
			output = p.outMax
		} else if output < p.outMin {
			output = p.outMin
		}
		p.output = output

		// Remember some variables for next time
		p.lastInput = input
		p.lastTime = now
		return true, nil
	} else {
		return false, nil
	}
}

func (p *Pid) SetInput(newInput float64) error {
	p.input = newInput
	return nil
}

func (p *Pid) SetSetpoint(newSetpoint float64) error {
	p.setpoint = newSetpoint
	return nil
}

func (p *Pid) SetBias(newBias float64) error {
	if newBias > p.outMax {
		p.bias = p.outMax // POSSIBLY INCORRECT
	} else if newBias < p.outMin {
		p.bias = p.outMin
	}
	p.bias = newBias
	return nil
}

/**
 * SetTunings(...)
 * This function allows the controller's dynamic performance to be adjusted.
 * it's called automatically from the constructor, but tunings can also
 * be adjusted on the fly during normal operation
 */
func (p *Pid) SetTunings(Kp, Ki, Kd float64) error {
	if Kp < 0 || Ki < 0 || Kd < 0 {
		return errors.New("invalid value: all tuning values must be positive")
	}

	p.displayP = Kp
	p.displayI = Ki
	p.displayD = Kd

	if Ki == 0 {
		p.iTerm = 0
	}

	SampleTimeInSec := p.intervalMillis / 1000
	p.kp = Kp
	p.ki = Ki * SampleTimeInSec
	p.kd = Kd / SampleTimeInSec

	return nil
}

/**
 * SetSampleTime(...)
 * sets the period, in Milliseconds, at which the calculation is performed
 */
func (p *Pid) SetSampleTime(newIntervalMillis float64) error {
	if newIntervalMillis > 0 {
		var ratio = newIntervalMillis / (p.intervalMillis)
		p.ki *= ratio
		p.kd /= ratio
		p.intervalMillis = math.Round(newIntervalMillis)
		return nil
	} else {
		return errors.New("invalid: interval value must be positive")
	}
}

/**
 * SetOutput( )
 * Set output level if in manual mode
 */
func (p *Pid) SetOutput(newOutput float64) error {
	if newOutput > p.outMax {
		newOutput = p.outMax // POSSIBLY INCORRECT
	} else if newOutput < p.outMin {
		newOutput = p.outMin
	}
	p.output = newOutput
	return nil
}

/**
 * SetOutputLimits(...)
 * This function will be used far more often than SetInputLimits.  while
 * the input to the controller will generally be in the 0-1023 range (which is
 * the default already,)  the output will be a little different.  maybe they'll
 * be doing a time window and will need 0-8000 or something.  or maybe they'll
 * want to clamp it from 0-125.  who knows.  at any rate, that can all be done here.
 */
func (p *Pid) SetOutputLimits(min, max float64) error {
	if min >= max {
		return errors.New("invalid values: min <= max")
	}
	p.outMin = min
	p.outMax = max

	if p.inAuto {
		if p.output > p.outMax {
			p.output = p.outMax
		} else if p.output < p.outMin {
			p.output = p.outMin
		}

		if p.iTerm > p.outMax-p.bias {
			p.iTerm = p.outMax - p.bias
		} else if p.iTerm < p.outMin-p.bias {
			p.iTerm = p.outMin - p.bias
		}
	}
	return nil
}

/**
 * SetMode(...)
 * Allows the controller Mode to be set to manual (0) or Automatic (non-zero)
 * when the transition from manual to auto occurs, the controller is
 * automatically initialized
 */
func (p *Pid) SetMode(newMode PID_MODE) error {
	/*  Removed in favor of manually triggered 'Reset'(using Initialize()).
	if (newAuto == !p.inAuto) {
	  //we just went from manual to auto
	  p.initialize()
	}
	*/
	if newMode == MANUAL {
		p.inAuto = false
		p.currMode = MANUAL
	} else if newMode == AUTO {
		p.inAuto = true
		p.currMode = AUTO
	} else {
		return errors.New("invalid value: mode setting is not a valid value")
	}
	return nil
}

/**
 * SetControllerDirection(...)
 * The PID will either be connected to a DIRECT acting process (+Output leads
 * to +Input) or a REVERSE acting process(+Output leads to -Input.)  we need to
 * know which one, because otherwise we may increase the output when we should
 * be decreasing.  This is called from the constructor.
 */
func (p *Pid) SetControllerDirection(newDirection PID_DIRECTION) error {
	p.direction = newDirection
	return nil
}

/**
 * Initialize()
 * does all the things that need to happen to ensure a bumpless transfer
 * from manual to automatic mode.
 */
func (p *Pid) Initialize() error {
	// p.iTerm = p.myOutput
	p.iTerm = 0
	p.output = p.bias
	p.lastInput = p.input
	/*
		  if (p.iTerm > p.outMax) {
			p.iTerm = p.outMax
		  } else if (p.iTerm < p.outMin) {
			p.iTerm = p.outMin
		  }
	*/
	p.Compute()
	return nil
}

/**
 * Status Functions
 * Just because you set the Kp=-1 doesn't mean it actually happened.  these
 * functions query the internal state of the PID.  they're here for display
 * purposes.  pid are the functions the PID Front-end uses for example
 */
func (p *Pid) GetKp() float64 {
	return p.displayP
}

func (p *Pid) getKi() float64 {
	return p.displayI
}

func (p *Pid) GetKd() float64 {
	return p.displayD
}

func (p *Pid) GetMode() PID_MODE {
	return p.currMode
}

func (p *Pid) GetDirection() PID_DIRECTION {
	return p.direction
}

func (p *Pid) GetOutput() float64 {
	return p.output
}

func (p *Pid) GetInput() float64 {
	return p.input
}

func (p *Pid) GetSetPoint() float64 {
	return p.setpoint
}

func (p *Pid) GetBias() float64 {
	return p.bias
}
