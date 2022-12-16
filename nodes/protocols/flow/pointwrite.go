package flow

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/enescakir/emoji"
)

type InputData struct {
	InputValue    *float64
	LinkConnected bool
}

type InputChanges struct {
	NewValue         bool
	LinkDisconnected bool
	LinkConnected    bool
}

type FFPointWrite struct {
	*node.Spec
	topic                  string
	netDevicePoint         string
	inputsArray            [17]InputData
	lastSendFail           bool
	lastPointPriorityWrite map[string]*float64
}

func NewFFPointWrite(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowPointWrite, category)
	in1 := node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs)
	in10 := node.BuildInput(node.In10, node.TypeFloat, nil, body.Inputs)
	in15 := node.BuildInput(node.In15, node.TypeFloat, nil, body.Inputs)
	in16 := node.BuildInput(node.In16, node.TypeFloat, nil, body.Inputs)
	inputs := node.BuildInputs(in1, in10, in15, in16)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body.SetAllowSettings()
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body = node.SetNoParent(body)
	return &FFPointWrite{body, "", "", [17]InputData{}, false, map[string]*float64{}}, nil
}

func (inst *FFPointWrite) set() {
	s := inst.GetStore()
	if s == nil {
		return
	}
	parentId := inst.GetParentId()
	nodeUUID := inst.GetID()
	d, ok := s.Get(parentId)
	var mqttData *pointStore
	if !ok {
		s.Set(parentId, &pointStore{
			parentID: parentId,
			payloads: []*pointDetails{&pointDetails{
				nodeUUID:       nodeUUID,
				topic:          inst.topic,
				netDevPntNames: inst.netDevicePoint,
				isWriteable:    true,
			}},
		}, 0)
	} else {
		mqttData = d.(*pointStore)
		payload := &pointDetails{
			nodeUUID:       nodeUUID,
			topic:          inst.topic,
			netDevPntNames: inst.netDevicePoint,
			isWriteable:    true,
		}
		mqttData, _ = addUpdatePayload(nodeUUID, mqttData, payload)
		s.Set(parentId, mqttData, 0)
	}
}

func (inst *FFPointWrite) Process() {
	_, firstLoop := inst.Loop()
	if firstLoop {
		selectedPoint, err := getPointSettings(inst.GetSettings())
		var setTopic bool
		if selectedPoint != nil && err == nil {
			if selectedPoint.Point != "" {
				t := makePointTopic(selectedPoint.Point)
				if t != "" {
					inst.topic = t
					inst.netDevicePoint = selectedPoint.Point
					inst.set()
					setTopic = true
				}
			}
		}
		if !setTopic {
			inst.SetWaringMessage("no point selected")
			inst.SetWaringIcon(string(emoji.OrangeCircle))

		}
	}

	val, null := inst.GetPayloadNull()
	var wroteValue bool
	if null {
		inst.WritePinNull(node.Out)
	} else {
		_, value, _, err := parseCOV(val)
		if err == nil {
			wroteValue = true
			inst.WritePin(node.Out, value)
		}
	}
	if !wroteValue {
		inst.WritePinNull(node.Out)
	}
}

func (inst *FFPointWrite) GetLastPriorityWrite() (priorityArrayWrite map[string]*float64) {
	return inst.lastPointPriorityWrite
}

func (inst *FFPointWrite) EvaluateInputsArray() map[string]*float64 {
	newInputArray := [17]InputData{}

	valueIn1 := inst.ReadPinAsFloatPointer(node.In1)
	linkIn1 := inst.InputHasConnection(node.In1)
	newInputArray[1] = InputData{valueIn1, linkIn1}

	valueIn10 := inst.ReadPinAsFloatPointer(node.In10)
	linkIn10 := inst.InputHasConnection(node.In10)
	newInputArray[10] = InputData{valueIn10, linkIn10}

	valueIn15 := inst.ReadPinAsFloatPointer(node.In15)
	linkIn15 := inst.InputHasConnection(node.In15)
	newInputArray[15] = InputData{valueIn15, linkIn15}

	valueIn16 := inst.ReadPinAsFloatPointer(node.In16)
	linkIn16 := inst.InputHasConnection(node.In16)
	newInputArray[16] = InputData{valueIn16, linkIn16}

	arraysMatch, arrayChanges := compareInputArrays(newInputArray, inst.inputsArray)
	priorityArrayWrite := make(map[string]*float64)

	if !arraysMatch {
		for i, changeData := range arrayChanges {
			if changeData == nil {
				continue
			}
			inputName := fmt.Sprintf("_%d", i)
			if changeData.LinkDisconnected { // has the link been disconnected
				priorityArrayWrite[inputName] = float.New(1)
				priorityArrayWrite[inputName] = nil
			} else if newInputArray[i].LinkConnected && (changeData.NewValue || changeData.LinkConnected) { // has there been a new value or a new link connected
				priorityArrayWrite[inputName] = newInputArray[i].InputValue
			}
		}
	}

	inst.inputsArray = newInputArray

	return priorityArrayWrite
}

func compareInputArrays(newInputs, oldInputs [17]InputData) (arraysMatch bool, changedValues [17]*InputChanges) {
	arraysMatch = true
	for i, inputData := range newInputs {
		if (inputData.InputValue != oldInputs[i].InputValue) || (inputData.LinkConnected != oldInputs[i].LinkConnected) {
			arraysMatch = false
			inputValueChanged := inputData.InputValue != oldInputs[i].InputValue
			inputLinkDisconnected := !inputData.LinkConnected && oldInputs[i].LinkConnected
			inputLinkConnected := inputData.LinkConnected && !oldInputs[i].LinkConnected
			changedValues[i] = &InputChanges{inputValueChanged, inputLinkDisconnected, inputLinkConnected}
		}
	}
	return
}

func (inst *FFPointWrite) GetSchema() *schemas.Schema {
	s := inst.buildSchema()
	return s
}
