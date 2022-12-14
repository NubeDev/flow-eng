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
	enable         bool
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

	pid := &Pid{false, input, setpoint, 0, MANUAL, dir, 0, interval, p, i, d, last, 0, 0, 0, 0, input, 0, 0, 100, false}
	pid.Compute()

	return pid
}

func (pid Pid) setInput(newInput float64) error {
	pid.input = newInput
	return nil
}

func (pid Pid) setSetpoint(newSetpoint float64) error {
	pid.setpoint = newSetpoint
	return nil
}

func (pid Pid) setBias(newBias float64) error {
	if newBias > pid.outMax {
		pid.bias = pid.outMax // POSSIBLY INCORRECT
	} else if newBias < pid.outMin {
		pid.bias = pid.outMin
	}
	pid.bias = newBias
	return nil
}

/**
 * Compute()
 * This, as they say, is where the magic happens.  pid function should be called
 * every time "void loop()" executes.  the function will decide for itself whether a new
 * pid Output needs to be computed.  returns true when the output is computed,
 * false when nothing has been done.
 */

func (pid Pid) Compute() (computed bool, err error) {
	if !pid.inAuto {
		return false, errors.New("pid controller is not enabled")
	}

	now := float64(time.Now().UnixMilli())
	timeChange := now - pid.lastTime
	if timeChange >= pid.intervalMillis {
		// Compute all the working error variables
		input := pid.input
		// var error = pid.mySetpoint - input
		errorAmount := input - pid.setpoint // above setpoint = positive error
		directionMultiplier := float64(1)
		if pid.direction {
			directionMultiplier = -1
		}
		errorAmount = errorAmount * directionMultiplier

		pid.iTerm += pid.kp * errorAmount
		if pid.iTerm > pid.outMax-pid.bias {
			pid.iTerm = pid.outMax - pid.bias
		} else if pid.iTerm < pid.outMin-pid.bias {
			pid.iTerm = pid.outMin - pid.bias
		}

		dInput := input - pid.lastInput

		// Compute PID Output
		// var output = ((pid.kp * error) + pid.iTerm - (pid.kd * dInput)) * pid.setDirection
		output := pid.kp*errorAmount + pid.iTerm - pid.kd*dInput + pid.bias
		// var output = ((pid.kp * error) + pid.iTerm - (pid.kd * dInput))

		if output > pid.outMax {
			output = pid.outMax
		} else if output < pid.outMin {
			output = pid.outMin
		}
		pid.output = output

		// Remember some variables for next time
		pid.lastInput = input
		pid.lastTime = now
		return true, nil
	} else {
		return false, nil
	}
}

/**
 * SetTunings(...)
 * This function allows the controller's dynamic performance to be adjusted.
 * it's called automatically from the constructor, but tunings can also
 * be adjusted on the fly during normal operation
 */
func (pid Pid) setTunings(Kp, Ki, Kd float64) error {
	if Kp < 0 || Ki < 0 || Kd < 0 {
		return errors.New("invalid value: all tuning values must be positive")
	}

	pid.displayP = Kp
	pid.displayI = Ki
	pid.displayD = Kd

	if Ki == 0 {
		pid.iTerm = 0
	}

	SampleTimeInSec := pid.intervalMillis / 1000
	pid.kp = Kp
	pid.ki = Ki * SampleTimeInSec
	pid.kd = Kd / SampleTimeInSec

	return nil
}

/**
 * SetSampleTime(...)
 * sets the period, in Milliseconds, at which the calculation is performed
 */
func (pid Pid) setSampleTime(newIntervalMillis float64) error {
	if newIntervalMillis > 0 {
		var ratio = newIntervalMillis / (1.0 * pid.intervalMillis)
		pid.ki *= ratio
		pid.kd /= ratio
		pid.intervalMillis = math.Round(newIntervalMillis)
		return nil
	} else {
		return errors.New("invalid: interval value must be positive")
	}
}

/**
 * SetOutput( )
 * Set output level if in manual mode
 */
func (pid Pid) setOutput(newOutput float64) error {
	if newOutput > pid.outMax {
		newOutput = pid.outMax // POSSIBLY INCORRECT
	} else if newOutput < pid.outMin {
		newOutput = pid.outMin
	}
	pid.output = newOutput
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
func (pid Pid) setOutputLimits(min, max float64) error {
	if min >= max {
		return errors.New("invalid values: min <= max")
	}
	pid.outMin = min
	pid.outMax = max

	if pid.inAuto {
		if pid.output > pid.outMax {
			pid.output = pid.outMax
		} else if pid.output < pid.outMin {
			pid.output = pid.outMin
		}

		if pid.iTerm > pid.outMax-pid.bias {
			pid.iTerm = pid.outMax - pid.bias
		} else if pid.iTerm < pid.outMin-pid.bias {
			pid.iTerm = pid.outMin - pid.bias
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
func (pid Pid) setMode(newMode PID_MODE) error {
	/*  Removed in favor of manually triggered 'Reset'(using Initialize()).
	if (newAuto == !pid.inAuto) {
	  //we just went from manual to auto
	  pid.initialize()
	}
	*/
	if newMode == MANUAL {
		pid.inAuto = false
	} else if newMode == AUTO {
		pid.inAuto = true
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
func (pid Pid) setControllerDirection(newDirection PID_DIRECTION) error {
	pid.direction = newDirection
	return nil
}

/**
 * Initialize()
 * does all the things that need to happen to ensure a bumpless transfer
 * from manual to automatic mode.
 */
func (pid Pid) initialize() error {
	// pid.iTerm = pid.myOutput
	pid.iTerm = 0
	pid.output = pid.bias
	pid.lastInput = pid.input
	/*
		  if (pid.iTerm > pid.outMax) {
			pid.iTerm = pid.outMax
		  } else if (pid.iTerm < pid.outMin) {
			pid.iTerm = pid.outMin
		  }
	*/
	return nil
}

/**
 * Status Functions
 * Just because you set the Kp=-1 doesn't mean it actually happened.  these
 * functions query the internal state of the PID.  they're here for display
 * purposes.  pid are the functions the PID Front-end uses for example
 */
func (pid Pid) getKp() float64 {
	return pid.displayP
}

func (pid Pid) getKi() float64 {
	return pid.displayI
}

func (pid Pid) getKd() float64 {
	return pid.displayD
}

func (pid Pid) getMode() PID_MODE {
	return pid.currMode
}

func (pid Pid) getDirection() PID_DIRECTION {
	return pid.direction
}

func (pid Pid) getOutput() float64 {
	return pid.output
}

func (pid Pid) getInput() float64 {
	return pid.input
}

func (pid Pid) getSetPoint() float64 {
	return pid.setpoint
}

func (pid Pid) getBias() float64 {
	return pid.bias
}
