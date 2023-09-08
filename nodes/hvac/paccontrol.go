package hvac

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/array"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	log "github.com/sirupsen/logrus"
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
	enable := node.BuildInput(node.Enable, node.TypeBool, true, body.Inputs, false, false)
	zoneTemp := node.BuildInput(node.ZoneTemp, node.TypeFloat, nil, body.Inputs, false, false)
	setPoint := node.BuildInput(node.Setpoint, node.TypeFloat, 22, body.Inputs, true, false)
	clgOffset := node.BuildInput(node.ClgOffset, node.TypeFloat, 1, body.Inputs, true, false)
	htgOffset := node.BuildInput(node.HtgOffset, node.TypeFloat, 1, body.Inputs, true, false)
	stgUpDelay := node.BuildInput(node.StgUpDelay, node.TypeFloat, 10, body.Inputs, true, false)
	modeChangeDelay := node.BuildInput(node.ModeChangeDelay, node.TypeFloat, 15, body.Inputs, true, false)
	econoAllow := node.BuildInput(node.EconoAllow, node.TypeBool, false, body.Inputs, false, false)
	oaTemp := node.BuildInput(node.OATemp, node.TypeFloat, nil, body.Inputs, false, false)
	econoHigh := node.BuildInput(node.EconoHigh, node.TypeFloat, 20, body.Inputs, true, false)
	econoLow := node.BuildInput(node.EconoLow, node.TypeFloat, 10, body.Inputs, true, false)
	fanStatus := node.BuildInput(node.FanStatus, node.TypeBool, nil, body.Inputs, false, false)
	clgLockout := node.BuildInput(node.ClgLockout, node.TypeBool, false, body.Inputs, false, false)
	htgLockout := node.BuildInput(node.HtgLockout, node.TypeBool, false, body.Inputs, false, false)

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
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	n := &PACControl{body, false, 0, false, false, 0, 0, 0, false, 0, false, 2}
	n.SetSchema(n.buildSchema())
	return n, nil
}

func (inst *PACControl) Process() {
	settings, err := inst.getSettings(inst.GetSettings())
	if err != nil {
		log.Errorf("PAC Control Node err: failed to get settings err:%s", err.Error())
		return
	}

	inst.numComps = 2 // TODO: replace with settings value
	enable, _ := inst.ReadPinAsBool(node.Enable)
	fanStatus := true
	requireFanStatus := settings.RequireFan
	if requireFanStatus {
		fanStatusInput, _ := inst.ReadPinAsBool(node.FanStatus)
		if fanStatusInput {
			fanStatus = true
			inst.fanStatus = true
		} else {
			if inst.fanStatus {

				fanStatusOffDelay := settings.FanStatusOffDelay
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
		setpoint := inst.ReadPinOrSettingsFloat(node.Setpoint)
		clgOffset := inst.ReadPinOrSettingsFloat(node.ClgOffset)
		htgOffset := inst.ReadPinOrSettingsFloat(node.HtgOffset)
		clgSP := setpoint + clgOffset
		htgSP := setpoint - htgOffset
		clgLockout := inst.ReadPinOrSettingsBool(node.ClgLockout)
		htgLockout := inst.ReadPinOrSettingsBool(node.HtgLockout)
		modeChangeDelay := inst.ReadPinOrSettingsFloat(node.ModeChangeDelay)
		if modeChangeDelay < 0 {
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

		stageUpDelay := inst.ReadPinOrSettingsFloat(node.StgUpDelay)
		if stageUpDelay <= 0 {
			stageUpDelay = 10
		}
		stageUpDelayDuration, _ := time.ParseDuration(fmt.Sprintf("%fm", stageUpDelay))

		conditionPastSetpoint := settings.SetpointMode

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
					inst.clgMode = false
					inst.clgModeEndTime = time.Now().Unix()
					inst.stageStartTime = 0
				}
			}
		}
		econoAllow := inst.ReadPinOrSettingsBool(node.EconoAllow)
		oaTemp, null := inst.ReadPinAsFloat(node.OATemp)
		if !null && econoAllow && !inst.htgMode {
			if inst.CheckEconoConditions(oaTemp, settings) {
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

func (inst *PACControl) CheckEconoConditions(oaTemp float64, settings *PACControlSettings) bool {
	econoHigh := inst.ReadPinOrSettingsFloat(node.EconoHigh)
	if econoHigh <= 0 {
		econoHigh = 22
	}
	econoLow := inst.ReadPinOrSettingsFloat(node.EconoLow)
	if econoLow < -10 {
		econoLow = 10
	}
	econoDB := settings.EconoDeadband
	if econoDB <= 0 {
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

// Custom Node Settings Schema

type PACControlSettingsSchema struct {
	Setpoint     schemas.Number  `json:"setpoint"`
	SetpointMode schemas.Boolean `json:"setpoint_mode"`
	ClgOffset    schemas.Number  `json:"clg-offset"`
	HtgOffset    schemas.Number  `json:"htg-offset"`
	StageDelay   schemas.Number  `json:"stage-up-delay"`
	ModeDelay    schemas.Number  `json:"mode-change-delay"`
	// TODO: add variable number of compressors
	EconoDeadband     schemas.Number  `json:"econo_deadband"`
	EconoHighLimit    schemas.Number  `json:"econo-high"`
	EconoLowLimit     schemas.Number  `json:"econo-low"`
	RequireFan        schemas.Boolean `json:"require_fan"`
	FanStatusOffDelay schemas.Number  `json:"fan_status_off_delay"`
	ClgLockout        schemas.Boolean `json:"clg-lockout"`
	HtgLockout        schemas.Boolean `json:"htg-lockout"`
}

type PACControlSettings struct {
	Setpoint     float64 `json:"setpoint"`
	SetpointMode bool    `json:"setpoint_mode"`
	ClgOffset    float64 `json:"clg-offset"`
	HtgOffset    float64 `json:"htg-offset"`
	StageDelay   float64 `json:"stage-up-delay"`
	ModeDelay    float64 `json:"mode-change-delay"`
	// TODO: add variable number of compressors
	EconoDeadband     float64 `json:"econo_deadband"`
	EconoHighLimit    float64 `json:"econo-high"`
	EconoLowLimit     float64 `json:"econo-low"`
	RequireFan        bool    `json:"require_fan"`
	FanStatusOffDelay float64 `json:"fan_status_off_delay"`
	ClgLockout        bool    `json:"clg-lockout"`
	HtgLockout        bool    `json:"htg-lockout"`
}

func (inst *PACControl) buildSchema() *schemas.Schema {
	props := &PACControlSettingsSchema{}

	// setpoint
	props.Setpoint.Title = "Setpoint"
	props.Setpoint.Default = 22

	// setpoint_mode  ie. Condition-Past-Setpoint
	props.SetpointMode.Title = "Setpoint Logic"
	props.SetpointMode.Default = false
	props.SetpointMode.EnumNames = []string{"Condition Past Setpoint", "Condition To Setpoint"}

	// offsets
	props.ClgOffset.Title = "Cooling Offset"
	props.ClgOffset.Default = 1
	props.HtgOffset.Title = "Heating Offset"
	props.HtgOffset.Default = 1

	// delays
	props.StageDelay.Title = "Stage Up Delay (minutes)"
	props.StageDelay.Default = 10
	props.ModeDelay.Title = "Mode Change Delay (minutes)"
	props.ModeDelay.Default = 15

	// economy mode
	props.EconoDeadband.Title = "Economy Deadband"
	props.EconoDeadband.Default = 1
	props.EconoHighLimit.Title = "Economy OA High Limit"
	props.EconoHighLimit.Default = 22
	props.EconoLowLimit.Title = "Economy OA Low Limit"
	props.EconoLowLimit.Default = 10

	// fan required
	props.RequireFan.Title = "Require Fan Status"
	props.RequireFan.Default = true
	props.FanStatusOffDelay.Title = "Fan Status Off Delay (seconds)"
	props.FanStatusOffDelay.Default = 15

	// lockouts
	props.ClgLockout.Title = "Cooling Lockout"
	props.ClgLockout.Default = false
	props.HtgLockout.Title = "Heating Lockout"
	props.HtgLockout.Default = false

	schema.Set(props)

	uiSchema := array.Map{
		"setpoint_mode": array.Map{
			"ui:widget": "select",
		},
		"ui:order": array.Slice{"setpoint", "setpoint_mode", "clg-offset", "htg-offset", "stage-up-delay", "mode-change-delay", "econo_deadband", "econo-high", "econo-low", "require_fan", "fan_status_off_delay", "clg-lockout", "htg-lockout"},
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

func (inst *PACControl) getSettings(body map[string]interface{}) (*PACControlSettings, error) {
	settings := &PACControlSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
