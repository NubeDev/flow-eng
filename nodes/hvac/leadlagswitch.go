package hvac

import (
	"github.com/NubeDev/flow-eng/node"
)

type LeadLagSwitch struct {
	*node.Spec
	lastRotate bool
	swapLead   bool
}

func NewLeadLagSwitch(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, leadLagSwitch, category)
	in := node.BuildInput(node.RotateLead, node.TypeBool, nil, body.Inputs, nil)
	stage1 := node.BuildInput(node.Stage1, node.TypeBool, nil, body.Inputs, nil)
	stage2 := node.BuildInput(node.Stage2, node.TypeBool, nil, body.Inputs, nil)
	inputs := node.BuildInputs(in, stage1, stage2)

	leadUnit := node.BuildOutput(node.LeadUnit, node.TypeString, nil, body.Outputs)
	leadUnitBool := node.BuildOutput(node.LeadUnitBool, node.TypeString, nil, body.Outputs)
	enableA := node.BuildOutput(node.EnableA, node.TypeBool, nil, body.Outputs)
	enableB := node.BuildOutput(node.EnableB, node.TypeBool, nil, body.Outputs)
	outputs := node.BuildOutputs(leadUnit, leadUnitBool, enableA, enableB)
	body = node.BuildNode(body, inputs, outputs, body.Settings)

	return &LeadLagSwitch{body, true, false}, nil
}

func (inst *LeadLagSwitch) Process() {
	rotate, _ := inst.ReadPinAsBool(node.RotateLead)
	if rotate && !inst.lastRotate {
		inst.swapLead = !inst.swapLead
	}
	inst.lastRotate = rotate

	stage1, _ := inst.ReadPinAsBool(node.Stage1)
	stage2, _ := inst.ReadPinAsBool(node.Stage2)
	if !inst.swapLead {
		inst.WritePin(node.LeadUnit, "A")
		inst.WritePin(node.LeadUnitBool, false)
		inst.WritePin(node.EnableA, stage1)
		inst.WritePin(node.EnableB, stage2)
	} else {
		inst.WritePin(node.LeadUnit, "B")
		inst.WritePin(node.LeadUnitBool, true)
		inst.WritePin(node.EnableA, stage2)
		inst.WritePin(node.EnableB, stage1)
	}
}
