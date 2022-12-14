package hvac

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"time"
)

type PIDNode struct {
	*node.Spec
	enable           bool
	fanStatusOffTime int64
	clgMode          bool
	htgMode          bool
	clgModeEnd       int64
	htgModeEnd       int64
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

func NewPIDNode(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, pacControlNode, category)
	enable := node.BuildInput(node.Enable, node.TypeBool, nil, body.Inputs)
	processValue := node.BuildInput(node.ProcessValue, node.TypeFloat, nil, body.Inputs)
	setPoint := node.BuildInput(node.Setpoint, node.TypeFloat, nil, body.Inputs)
	minOut := node.BuildInput(node.MinOut, node.TypeFloat, nil, body.Inputs)
	maxOut := node.BuildInput(node.MaxOut, node.TypeFloat, nil, body.Inputs)
	inP := node.BuildInput(node.InP, node.TypeFloat, nil, body.Inputs)
	inI := node.BuildInput(node.InI, node.TypeFloat, nil, body.Inputs)
	inD := node.BuildInput(node.InD, node.TypeFloat, nil, body.Inputs)
	direction := node.BuildInput(node.PIDDirection, node.TypeBool, nil, body.Inputs)
	intervalSecs := node.BuildInput(node.IntervalSecs, node.TypeFloat, nil, body.Inputs)
	bias := node.BuildInput(node.Bias, node.TypeFloat, nil, body.Inputs)
	manual := node.BuildInput(node.Manual, node.TypeFloat, nil, body.Inputs)
	reset := node.BuildInput(node.Reset, node.TypeBool, nil, body.Inputs)

	output := node.BuildOutput(node.Out, node.TypeBool, nil, body.Outputs)

	inputs := node.BuildInputs(enable, processValue, setPoint, minOut, maxOut, inP, inI, inD, direction, intervalSecs, bias, manual, reset)
	outputs := node.BuildOutputs(output)
	body = node.BuildNode(body, inputs, outputs, nil)
	body.SetSchema(buildSchema())
	return &PIDNode{body, false, }, nil
}

func (inst *PIDNode) Process() {
	inst.numComps = 2 // TODO: replace with settings value
	enable, null := inst.ReadPinAsBool(node.Enable)
	if null {
		enable = false
	}
	fanStatus := true
	requireFanStatus := true // TODO: replace with settings value
	if requireFanStatus {
		fanStatusInput, null := inst.ReadPinAsBool(node.FanStatus)
		if null {
			fanStatusInput = false
		}
		if fanStatusInput {
			fanStatus = true
			inst.fanStatus = true
		} else {
			if inst.fanStatus {
				fanStatusOffDelay := 30 // TODO: replace with settings value
				inst.fanStatusOffTime = time.Now().Add(time.Second * time.Duration(fanStatusOffDelay)).Unix()
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
		if !htgLockout && !inst.clgMode && !inst.htgMode && (zoneTemp < htgSP) {
			if inst.clgModeEnd == 0 || (time.Now().Unix() >= time.Unix(inst.clgModeEnd, 0).Add(time.Duration(modeChangeDelay)*time.Minute).Unix()) {
				inst.htgMode = true
				inst.compStage = 0
			}
		} else if !clgLockout && !inst.clgMode && !inst.htgMode && (zoneTemp > clgSP) {
			if inst.htgModeEnd == 0 || (time.Now().Unix() >= time.Unix(inst.htgModeEnd, 0).Add(time.Duration(modeChangeDelay)*time.Minute).Unix()) {
				inst.clgMode = true
				inst.compStage = 0
			}
		}

		stageUpDelay, null := inst.ReadPinAsFloat(node.StgUpDelay) // TODO: update with settings value if no input value
		if null {
			stageUpDelay = 10
		}
		conditionPastSetpoint := false // TODO: update with settings value

		if inst.htgMode {
			inst.econoConditions = false
			if htgLockout {
				inst.compStage = 0
				inst.htgMode = false
				inst.htgModeEnd = time.Now().Unix()
				inst.stageStartTime = 0
			} else if zoneTemp < htgSP {
				if inst.compStage == 0 {
					inst.compStage = 1
					inst.stageStartTime = time.Now().Unix()
				} else if inst.compStage > 0 && (time.Now().Unix() >= time.Unix(inst.stageStartTime, 0).Add(time.Duration(stageUpDelay)*time.Minute).Unix()) {
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
					inst.htgModeEnd = time.Now().Unix()
					inst.stageStartTime = 0
				}
			}
		} else if inst.clgMode {
			if clgLockout {
				inst.compStage = 0
				inst.clgMode = false
				inst.clgModeEnd = time.Now().Unix()
				inst.stageStartTime = 0
			} else if zoneTemp > clgSP {
				if inst.compStage == 0 {
					inst.compStage = 1
					inst.stageStartTime = time.Now().Unix()
				} else if inst.compStage > 0 && (time.Now().Unix() >= time.Unix(inst.stageStartTime, 0).Add(time.Duration(stageUpDelay)*time.Minute).Unix()) {
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
					inst.htgModeEnd = time.Now().Unix()
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
	} else {
		// For disabled controller
		inst.DisablePAC()
	}

}

func (inst *PIDNode) CheckEconoConditions(oaTemp float64) bool {
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

func (inst *PIDNode) DisablePAC() {
	inst.clgMode = false
	inst.WritePinFalse(node.ClgMode)
	inst.htgMode = false
	inst.WritePinFalse(node.HtgMode)
	inst.compStage = 0
	inst.WritePinFloat(node.CompStage, 0)
	inst.econoMode = false
	inst.WritePinFalse(node.EconoMode)
	inst.WritePinFloat(node.OADamper, 0)
	inst.WritePinFalse(node.ReversingValve)
	inst.WritePinFalse(node.Compressor1)
	inst.WritePinFalse(node.Compressor2)
}

func (inst *PIDNode) SetOutputs(zoneTemp, setpoint, clgSP float64) {
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

func (inst *PIDNode) DisableEconoMode() {
	inst.WritePinFalse(node.EconoMode)
	inst.WritePinFloat(node.OADamper, 0)
}
