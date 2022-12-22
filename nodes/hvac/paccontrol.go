package hvac

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type PACControl struct {
	*node.Spec
	fanStatus        bool
	fanStatusOffTime int64
	clgMode          bool
	htgMode          bool
	clgModeEndTime   int64
	htgModeEndTime   int64
	compStage        int
	econoMode        bool
	stageStartTime   int64
	econoConditions  bool
	numComps         int
}

// input
// enable
// sp
// cool offset
// heat offset

func NewPACControl(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, pacControlNode, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs)
	zoneTemp := node.BuildInput(node.ZoneTemp, node.TypeFloat, nil, body.Inputs)
	setPoint := node.BuildInput(node.Setpoint, node.TypeFloat, nil, body.Inputs)
	clgOffset := node.BuildInput(node.ClgOffset, node.TypeFloat, nil, body.Inputs)
	htgOffset := node.BuildInput(node.HtgOffset, node.TypeFloat, nil, body.Inputs)
	stgUpDelay := node.BuildInput(node.StgUpDelay, node.TypeFloat, nil, body.Inputs)
	modeChangeDelay := node.BuildInput(node.ModeChangeDelay, node.TypeFloat, nil, body.Inputs)
	econoAllow := node.BuildInput(node.EconoAllow, node.TypeBool, nil, body.Inputs)
	oaTemp := node.BuildInput(node.OATemp, node.TypeFloat, nil, body.Inputs)
	econoHigh := node.BuildInput(node.EconoHigh, node.TypeFloat, nil, body.Inputs)
	econoLow := node.BuildInput(node.EconoLow, node.TypeFloat, nil, body.Inputs)
	fanStatus := node.BuildInput(node.FanStatus, node.TypeBool, nil, body.Inputs)
	clgLockout := node.BuildInput(node.ClgLockout, node.TypeBool, nil, body.Inputs)
	htgLockout := node.BuildInput(node.HtgLockout, node.TypeBool, nil, body.Inputs)

	clgMode := node.BuildOutput(node.ClgMode, node.TypeBool, nil, body.Outputs)
	htgMode := node.BuildOutput(node.HtgMode, node.TypeBool, nil, body.Outputs)
	compStage := node.BuildOutput(node.CompStage, node.TypeFloat, nil, body.Outputs)
	econoMode := node.BuildOutput(node.EconoMode, node.TypeBool, nil, body.Outputs)
	oaDamper := node.BuildOutput(node.OADamper, node.TypeFloat, nil, body.Outputs)
	revValve := node.BuildOutput(node.ReversingValve, node.TypeBool, nil, body.Outputs)
	comp1 := node.BuildOutput(node.Compressor1, node.TypeBool, nil, body.Outputs)
	comp2 := node.BuildOutput(node.Compressor2, node.TypeBool, nil, body.Outputs)

	inputs := node.BuildInputs(enable, zoneTemp, setPoint, clgOffset, htgOffset, stgUpDelay, modeChangeDelay, econoAllow, oaTemp, econoHigh, econoLow, fanStatus, clgLockout, htgLockout)
	outputs := node.BuildOutputs(clgMode, htgMode, compStage, econoMode, oaDamper, revValve, comp1, comp2)
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetSchema(buildSchema())
	return &PACControl{body, false, 0, false, false, 0, 0, 0, false, 0, false, 2}, nil
}

func (inst *PACControl) Process() {
	inst.numComps = 2 // TODO: replace with settings value
	enable, null := inst.ReadPinAsBool(node.Enable)
	fanStatus := true
	requireFanStatus := true // TODO: replace with settings value
	if requireFanStatus {
		fanStatusInput, _ := inst.ReadPinAsBool(node.FanStatus)
		if fanStatusInput {
			fanStatus = true
			inst.fanStatus = true
		} else {
			if inst.fanStatus {
				fanStatusOffDelay := 30 // TODO: replace with settings value
				fanStatusOffDelayDuration, _ := time.ParseDuration(fmt.Sprintf("%fs", fanStatusOffDelay))
				inst.fanStatusOffTime = time.Now().Add(fanStatusOffDelayDuration).Unix()
			}
			inst.fanStatus = fanStatusInput
			if time.Now().Unix() >= inst.fanStatusOffTime {
				fanStatus = false
			}
		}
	}

	disable := false
	if !enable || (requireFanStatus && !fanStatus) {
		disable = true
	}
	zoneTemp, null := inst.ReadPinAsFloat(node.ZoneTemp)
	if null {
		disable = true
	}

	if !disable {
		setpoint, null := inst.ReadPinAsFloat(node.Setpoint) // TODO: update with settings value if no input value
		if null {
			inst.DisablePAC()
			return
		}
		clgOffset, null := inst.ReadPinAsFloat(node.ClgOffset) // TODO: update with settings value if no input value
		if null {
			inst.DisablePAC()
			return
		}
		htgOffset, null := inst.ReadPinAsFloat(node.HtgOffset) // TODO: update with settings value if no input value
		if null {
			inst.DisablePAC()
			return
		}
		clgSP := setpoint + clgOffset
		htgSP := setpoint - htgOffset
		clgLockout, null := inst.ReadPinAsBool(node.ClgLockout) // TODO: update with settings value if no input value
		if null {
			clgLockout = false
		}
		htgLockout, null := inst.ReadPinAsBool(node.HtgLockout) // TODO: update with settings value if no input value
		if null {
			htgLockout = false
		}
		modeChangeDelay, null := inst.ReadPinAsFloat(node.ModeChangeDelay) // TODO: update with settings value if no input value
		if null {
			modeChangeDelay = 15
		}
		modeChangeDelayDuration, _ := time.ParseDuration(fmt.Sprintf("%fm", modeChangeDelay))

		if !htgLockout && !inst.clgMode && !inst.htgMode && (zoneTemp < htgSP) {
			if inst.clgModeEndTime == 0 || (time.Now().Unix() >= time.Unix(inst.clgModeEndTime, 0).Add(modeChangeDelayDuration).Unix()) {
				inst.htgMode = true
				inst.compStage = 0
			}
		} else if !clgLockout && !inst.clgMode && !inst.htgMode && (zoneTemp > clgSP) {
			if inst.htgModeEndTime == 0 || (time.Now().Unix() >= time.Unix(inst.htgModeEndTime, 0).Add(modeChangeDelayDuration).Unix()) {
				inst.clgMode = true
				inst.compStage = 0
			}
		}

		stageUpDelay, null := inst.ReadPinAsFloat(node.StgUpDelay) // TODO: update with settings value if no input value
		if null {
			stageUpDelay = 10
		}
		stageUpDelayDuration, _ := time.ParseDuration(fmt.Sprintf("%fm", stageUpDelay))

		conditionPastSetpoint := false // TODO: update with settings value

		if inst.htgMode {
			inst.econoConditions = false
			if htgLockout {
				inst.compStage = 0
				inst.htgMode = false
				inst.htgModeEndTime = time.Now().Unix()
				inst.stageStartTime = 0
			} else if zoneTemp < htgSP {
				if inst.compStage == 0 {
					inst.compStage = 1
					inst.stageStartTime = time.Now().Unix()
				} else if inst.compStage > 0 && (time.Now().Unix() >= time.Unix(inst.stageStartTime, 0).Add(stageUpDelayDuration).Unix()) {
					if inst.compStage < inst.numComps {
						inst.compStage++
						inst.stageStartTime = time.Now().Unix()
					}
				}
			} else {
				offSPInterval := htgOffset / float64(inst.numComps)
				if !conditionPastSetpoint {
					for s := inst.compStage; s > 0; s-- {
						if zoneTemp >= (setpoint - (offSPInterval * (float64(s) - 1))) {
							inst.compStage--
							inst.stageStartTime = time.Now().Unix()
						}
					}
				} else if conditionPastSetpoint {
					offSPInterval = (htgOffset * 2) / float64(inst.numComps)
					for s := inst.compStage; s > 0; s-- {
						if zoneTemp >= ((setpoint + htgOffset) - (offSPInterval * (float64(s) - 1))) {
							inst.compStage--
							inst.stageStartTime = time.Now().Unix()
						}
					}
				}
				if inst.compStage == 0 {
					inst.htgMode = false
					inst.htgModeEndTime = time.Now().Unix()
					inst.stageStartTime = 0
				}
			}
		} else if inst.clgMode {
			if clgLockout {
				inst.compStage = 0
				inst.clgMode = false
				inst.clgModeEndTime = time.Now().Unix()
				inst.stageStartTime = 0
			} else if zoneTemp > clgSP {
				if inst.compStage == 0 {
					inst.compStage = 1
					inst.stageStartTime = time.Now().Unix()
				} else if inst.compStage > 0 && (time.Now().Unix() >= time.Unix(inst.stageStartTime, 0).Add(stageUpDelayDuration).Unix()) {
					if inst.compStage < inst.numComps {
						inst.compStage++
						inst.stageStartTime = time.Now().Unix()
					}
				}
			} else {
				offSPInterval := clgOffset / float64(inst.numComps)
				if !conditionPastSetpoint {
					for s := inst.compStage; s > 0; s-- {
						if zoneTemp <= (setpoint + (offSPInterval * (float64(s) - 1))) {
							inst.compStage--
							inst.stageStartTime = time.Now().Unix()
						}
					}
				} else if conditionPastSetpoint {
					offSPInterval = (clgOffset * 2) / float64(inst.numComps)
					for s := inst.compStage; s > 0; s-- {
						if zoneTemp <= ((setpoint - clgOffset) + (offSPInterval * (float64(s) - 1))) {
							inst.compStage--
							inst.stageStartTime = time.Now().Unix()
						}
					}
				}
				if inst.compStage == 0 {
					inst.htgMode = false
					inst.htgModeEndTime = time.Now().Unix()
					inst.stageStartTime = 0
				}
			}
		}
		econoAllow, null := inst.ReadPinAsBool(node.EconoAllow) // TODO: update with settings value if no input value
		if null {
			econoAllow = true
		}
		oaTemp, null := inst.ReadPinAsFloat(node.OATemp)
		if !null && econoAllow && !inst.htgMode {
			if inst.CheckEconoConditions(oaTemp) {
				inst.econoConditions = true
			}
		} else {
			inst.econoConditions = false
		}

		inst.SetOutputs(zoneTemp, setpoint, clgSP)

	} else {
		// For disabled controller
		inst.DisablePAC()
	}

}

func (inst *PACControl) CheckEconoConditions(oaTemp float64) bool {
	econoHigh, null := inst.ReadPinAsFloat(node.EconoHigh) // TODO: update with settings value if no input value
	if null {
		econoHigh = 22
	}
	econoLow, null := inst.ReadPinAsFloat(node.EconoLow) // TODO: update with settings value if no input value
	if null {
		econoHigh = 20
	}
	econoDB := float64(1) // TODO: update with settings value
	if null {
		econoDB = float64(1)
	}
	if econoHigh <= (econoLow+econoDB) || econoHigh == 0 || econoLow == 0 || econoDB == 0 {
		inst.econoConditions = false
		return false
	}
	if oaTemp < econoLow || oaTemp > econoHigh {
		inst.econoConditions = false
		return false
	}
	if oaTemp <= (econoHigh-econoDB) && oaTemp >= (econoLow+econoDB) {
		inst.econoConditions = true
	}
	return inst.econoConditions
}

func (inst *PACControl) DisablePAC() {
	inst.clgMode = false
	inst.htgMode = false
	inst.compStage = 0
	inst.econoMode = false
	inst.econoConditions = false
	inst.stageStartTime = 0
	inst.WritePinFalse(node.ClgMode)
	inst.WritePinFalse(node.HtgMode)
	inst.WritePinFloat(node.CompStage, 0)
	inst.WritePinFalse(node.EconoMode)
	inst.WritePinFloat(node.OADamper, 0)
	inst.WritePinFalse(node.ReversingValve)
	inst.WritePinFalse(node.Compressor1)
	inst.WritePinFalse(node.Compressor2)
}

func (inst *PACControl) SetOutputs(zoneTemp, setpoint, clgSP float64) {
	inst.WritePinBool(node.ClgMode, inst.clgMode)
	inst.WritePinBool(node.HtgMode, inst.htgMode)
	inst.WritePinFloat(node.CompStage, float64(inst.compStage))
	inst.WritePinBool(node.EconoMode, inst.econoConditions)
	if inst.econoConditions && inst.clgMode {
		oaDamperCmd := float.Scale(zoneTemp, setpoint, clgSP, 0, 100)
		inst.WritePinFloat(node.OADamper, oaDamperCmd)
	} else {
		inst.WritePinFloat(node.OADamper, 0)
	}
	inst.WritePinBool(node.ReversingValve, inst.htgMode)
	for x := 0; x < inst.numComps; x++ {
		if x == 0 {
			if inst.compStage >= (x + 1) {
				inst.WritePinTrue(node.Compressor1)
			} else {
				inst.WritePinFalse(node.Compressor1)
			}
		} else if x == 1 {
			if inst.compStage >= (x + 1) {
				inst.WritePinTrue(node.Compressor2)
			} else {
				inst.WritePinFalse(node.Compressor2)
			}
		}
	}
}

func (inst *PACControl) DisableEconoMode() {
	inst.WritePinFalse(node.EconoMode)
	inst.WritePinFloat(node.OADamper, 0)
}
