package pid

import (
	"github.com/NubeDev/flow-eng/node"
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
	DIRECT   	PID_DIRECTION = false
	REVERSE 	PID_DIRECTION = true
)

type PID_MODE bool

const (
	MANUAL   	PID_MODE = false
	AUTO 		PID_MODE = true
)

type Pid struct {
	enable			bool			`json:"enable"`
	input			float64			`json:"input"`
	setpoint		float64			`json:"setpoint"`
	output			float64			`json:"output"`
	currMode		PID_MODE		`json:"curr_mode"`
	direction		PID_DIRECTION	`json:"direction"`
	bias			float64			`json:"bias"`
	intervalSecs	float64			`json:"interval_secs"`
	displayP		float64			`json:"display_p"`
	displayI		float64			`json:"display_i"`
	displayD		float64			`json:"display_d"`
	lastTime		float64			`json:"last_time"`
	pFactor			float64			`json:"p_factor"`
	iFactor			float64			`json:"i_factor"`
	dFactor			float64			`json:"d_factor"`
	outputSum		float64			`json:"output_sum"`
	lastInput		float64			`json:"last_input"`
	iTerm			float64			`json:"i_term"`
	outMin			float64			`json:"out_min"`
	outMax			float64			`json:"out_max"`

}

func Init(input, setpoint, p, i, d float64, dir PID_DIRECTION) *Pid {
	if p < 0 {
		p = 0
	}
	if i < 0 {
		i = 0
	}
	if d < 0 {
		d = 0
	}
	interval := float64(10)
	last := float64(time.Now().Unix()) - interval

	pid := &Pid{false, input, setpoint, 0, MANUAL, DIRECT, 0, interval, p, i, d, last, }

	this.lastTime = this.millis() - this.SampleTime
	
	this.myBias = 0
	this.ITerm = 0
	this.myOutput = 0
}

// Constants for backward compatibility
// PID.AUTOMATIC = 1
// PID.MANUAL = 0
// PID.DIRECT = 0
// PID.REVERSE = 1

PID.prototype.setInput = function(current_value) {
const input = Number(current_value)
if (isNaN(input)) return
this.input = input
}

PID.prototype.setPoint = function(current_value) {
const setpoint = Number(current_value)
if (isNaN(setpoint)) return
this.mySetpoint = setpoint
}

PID.prototype.setBias = function(current_value) {
let bias = Number(current_value)
if (isNaN(bias)) return
if (bias > this.outMax) {
bias = this.outMax // POSSIBLY INCORRECT
} else if (bias < this.outMin) {
bias = this.outMin
}
this.myBias = bias
}

PID.prototype.millis = function() {
var d = new Date()
return d.getTime()
}

/**
 * Compute()
 * This, as they say, is where the magic happens.  this function should be called
 * every time "void loop()" executes.  the function will decide for itself whether a new
 * pid Output needs to be computed.  returns true when the output is computed,
 * false when nothing has been done.
 */

PID.prototype.compute = function() {
if (!this.inAuto) {
return false
}

var now = this.millis()
var timeChange = now - this.lastTime
if (timeChange >= this.SampleTime) {
// Compute all the working error variables
var input = this.input

//var error = this.mySetpoint - input
var error = input - this.mySetpoint // above setpoint = positive error
error = error * this.setDirection

this.ITerm += this.ki * error
if (this.ITerm > this.outMax - this.myBias) {
this.ITerm = this.outMax - this.myBias
} else if (this.ITerm < this.outMin - this.myBias) {
this.ITerm = this.outMin - this.myBias
}

var dInput = input - this.lastInput

// Compute PID Output
//var output = ((this.kp * error) + this.ITerm - (this.kd * dInput)) * this.setDirection
var output = this.kp * error + this.ITerm - this.kd * dInput + this.myBias
//var output = ((this.kp * error) + this.ITerm - (this.kd * dInput))

if (output > this.outMax) {
output = this.outMax
} else if (output < this.outMin) {
output = this.outMin
}
this.myOutput = output

// Remember some variables for next time
this.lastInput = input
this.lastTime = now
return true
} else {
return false
}
}

/**
 * SetTunings(...)
 * This function allows the controller's dynamic performance to be adjusted.
 * it's called automatically from the constructor, but tunings can also
 * be adjusted on the fly during normal operation
 */
PID.prototype.setTunings = function(Kp, Ki, Kd) {
const Kp_num = Number(Kp)
const Ki_num = Number(Ki)
const Kd_num = Number(Kd)
if (isNaN(Kp_num) || isNaN(Ki_num) || isNaN(Kd_num)) return
if (Kp < 0 || Ki < 0 || Kd < 0) {
return
}

this.dispKp = Kp
this.dispKi = Ki
this.dispKd = Kd

if(Ki == 0) this.ITerm = 0

this.SampleTimeInSec = this.SampleTime / 1000
this.kp = Kp
this.ki = Ki * this.SampleTimeInSec
this.kd = Kd / this.SampleTimeInSec
}

/**
 * SetSampleTime(...)
 * sets the period, in Milliseconds, at which the calculation is performed
 */
PID.prototype.setSampleTime = function(NewSampleTime) {
let NewSampleTime_num = Number(NewSampleTime)
if (isNaN(NewSampleTime_num)) return

if (NewSampleTime > 0) {
var ratio = NewSampleTime / (1.0 * this.SampleTime)
this.ki *= ratio
this.kd /= ratio
this.SampleTime = Math.round(NewSampleTime)
}
}

/**
 * SetOutput( )
 * Set output level if in manual mode
 */
PID.prototype.setOutput = function(val) {
let val_num = Number(val)
if (isNaN(val_num)) return

if (val > this.outMax) {
val = this.outMax // POSSIBLY INCORRECT
} else if (val < this.outMin) {
val = this.outMin
}
this.myOutput = val
}

/**
 * SetOutputLimits(...)
 * This function will be used far more often than SetInputLimits.  while
 * the input to the controller will generally be in the 0-1023 range (which is
 * the default already,)  the output will be a little different.  maybe they'll
 * be doing a time window and will need 0-8000 or something.  or maybe they'll
 * want to clamp it from 0-125.  who knows.  at any rate, that can all be done here.
 */
PID.prototype.setOutputLimits = function(Min, Max) {
let Min_num = Number(Min)
let Max_num = Number(Max)
if (isNaN(Min_num) || isNaN(Max_num)) return

if (Min >= Max) {
return
}
this.outMin = Min
this.outMax = Max

if (this.inAuto) {
if (this.myOutput > this.outMax) {
this.myOutput = this.outMax
} else if (this.myOutput < this.outMin) {
this.myOutput = this.outMin
}

if (this.ITerm > this.outMax - this.myBias) {
this.ITerm = this.outMax - this.myBias
} else if (this.ITerm < this.outMin - this.myBias) {
this.ITerm = this.outMin - this.myBias
}
}
}

/**
 * SetMode(...)
 * Allows the controller Mode to be set to manual (0) or Automatic (non-zero)
 * when the transition from manual to auto occurs, the controller is
 * automatically initialized
 */
PID.prototype.setMode = function(Mode) {
var newAuto
if (
Mode == 0 ||
Mode.toString().toLowerCase() == 'automatic' ||
Mode.toString().toLowerCase() == 'auto'
) {
newAuto = 1
} else if (Mode == 1 || Mode.toString().toLowerCase() == 'manual') {
newAuto = 0
} else {
throw new Error('Incorrect Mode Chosen')
}
/*  Removed in favor of manually triggered 'Reset'(using Initialize()).
if (newAuto == !this.inAuto) {
  //we just went from manual to auto
  this.initialize()
}
*/
this.inAuto = newAuto
}

/**
 * Initialize()
 * does all the things that need to happen to ensure a bumpless transfer
 * from manual to automatic mode.
 */
PID.prototype.initialize = function() {
//this.ITerm = this.myOutput
this.ITerm = 0
this.myOutput = this.myBias
this.lastInput = this.input
/*
  if (this.ITerm > this.outMax) {
    this.ITerm = this.outMax
  } else if (this.ITerm < this.outMin) {
    this.ITerm = this.outMin
  }
*/
}

/**
 * SetControllerDirection(...)
 * The PID will either be connected to a DIRECT acting process (+Output leads
 * to +Input) or a REVERSE acting process(+Output leads to -Input.)  we need to
 * know which one, because otherwise we may increase the output when we should
 * be decreasing.  This is called from the constructor.
 */
PID.prototype.setControllerDirection = function(ControllerDirection) {
if (ControllerDirection == 0 || ControllerDirection.toString().toLowerCase() == 'direct') {
this.setDirection = 1
} else if (
ControllerDirection == 1 ||
ControllerDirection.toString().toLowerCase() == 'reverse'
) {
this.setDirection = -1
} else {
throw new Error('Incorrect Controller Direction Chosen')
}
}

/**
 * Status Functions
 * Just because you set the Kp=-1 doesn't mean it actually happened.  these
 * functions query the internal state of the PID.  they're here for display
 * purposes.  this are the functions the PID Front-end uses for example
 */
PID.prototype.getKp = function() {
return this.dispKp
}

PID.prototype.getKd = function() {
return this.dispKd
}

PID.prototype.getKi = function() {
return this.dispKi
}

PID.prototype.getMode = function() {
return this.inAuto ? 'Auto' : 'Manual'
}

PID.prototype.getDirection = function() {
return this.controllerDirection
}

PID.prototype.getOutput = function() {
return this.myOutput
}

PID.prototype.getInput = function() {
return this.input
}

PID.prototype.getSetPoint = function() {
return this.mySetpoint
}

PID.prototype.getBias = function() {
return this.myBias
}

module.exports = PID

