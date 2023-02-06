package flow

import (
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/enescakir/emoji"
	"time"
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
	lastUpdate             time.Time
}

func NewFFPointWrite(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, flowPointWrite, category)
	in1 := node.BuildInput(node.In1, node.TypeFloat, nil, body.Inputs, false)
	in10 := node.BuildInput(node.In10, node.TypeFloat, nil, body.Inputs, false)
	in15 := node.BuildInput(node.In15, node.TypeFloat, nil, body.Inputs, false)
	in16 := node.BuildInput(node.In16, node.TypeFloat, nil, body.Inputs, false)
	inputs := node.BuildInputs(in1, in10, in15, in16)
	value := node.BuildOutput(node.Out, node.TypeFloat, nil, body.Outputs)
	currentPriority := node.BuildOutput(node.CurrentPriority, node.TypeFloat, nil, body.Outputs)
	lastUpdated := node.BuildOutput(node.LastUpdated, node.TypeString, nil, body.Outputs)
	outputs := node.BuildOutputs(value, currentPriority, lastUpdated)
	body.SetAllowSettings()
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	body = node.SetNoParent(body)
	return &FFPointWrite{body, "", "", [17]InputData{}, false, map[string]*float64{}, time.Now()}, nil
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

func (inst *FFPointWrite) checkStillExists() {
	s := inst.GetStore()
	if s == nil {
		return
	}
	parentId := inst.GetParentId()
	topic := fmt.Sprintf("pointsList_%s", parentId)
	children, ok := s.Get(topic)
	if ok {
		existingPoints := children.([]*point)
		var pointExists bool
		for _, existingPoint := range existingPoints {
			t := makePointTopic(existingPoint.Name)
			if t == inst.topic {
				pointExists = true
				inst.SetSubTitle(existingPoint.Name)
				inst.SetWaringMessage("")
				inst.SetWaringIcon(string(emoji.GreenCircle))
			}
		}
		if !pointExists {
			inst.SetWaringMessage(pointError)
			inst.SetWaringIcon(string(emoji.OrangeCircle))
			inst.SetSubTitle("")
		}
	}
}

func (inst *FFPointWrite) setTopic() {
	selectedPoint, err := getPointSettings(inst.GetSettings())
	if selectedPoint != nil && err == nil {
		if selectedPoint.Point != "" {
			t := makePointTopic(selectedPoint.Point)
			if t != "" {
				inst.topic = t
				inst.netDevicePoint = selectedPoint.Point
				inst.SetSubTitle(selectedPoint.Point)
				inst.SetWaringMessage("")
				inst.SetWaringIcon(string(emoji.GreenCircle))
				inst.set()
			} else {
				inst.SetWaringMessage(pointError)
				inst.SetWaringIcon(string(emoji.OrangeCircle))
				inst.SetSubTitle("")
			}
		}
	}
}

func (inst *FFPointWrite) Process() {
	loopCount, firstLoop := inst.Loop()
	if firstLoop {
		inst.setTopic()
	}
	if loopCount%50 == 0 {
		inst.setTopic()
		inst.checkStillExists()
	}
	val := inst.GetPayload()
	if val == nil {
		inst.WritePinNull(node.Out)
	} else {
		_, value, currentPri, err := parseCOV(val.Any)
		if err == nil {
			inst.WritePinFloat(node.Out, value)
			inst.WritePinFloat(node.CurrentPriority, float64(currentPri))
			inst.lastUpdate = val.LastUpdate
		}
	}
	inst.WritePin(node.LastUpdated, ttime.TimeSince(inst.lastUpdate))

}

func (inst *FFPointWrite) GetLastPriorityWrite() (priorityArrayWrite map[string]*float64) {
	return inst.lastPointPriorityWrite
}

func (inst *FFPointWrite) EvaluateInputsArray(forceResend bool) map[string]*float64 {
	newInputArray := [17]InputData{}

	valueIn1 := inst.ReadPinAsFloatPointer(node.In1)
	linkIn1 := inst.InputHasConnectionOrValue(node.In1)
	newInputArray[1] = InputData{valueIn1, linkIn1}

	valueIn10 := inst.ReadPinAsFloatPointer(node.In10)
	linkIn10 := inst.InputHasConnectionOrValue(node.In10)
	newInputArray[10] = InputData{valueIn10, linkIn10}

	valueIn15 := inst.ReadPinAsFloatPointer(node.In15)
	linkIn15 := inst.InputHasConnectionOrValue(node.In15)
	newInputArray[15] = InputData{valueIn15, linkIn15}

	valueIn16 := inst.ReadPinAsFloatPointer(node.In16)
	linkIn16 := inst.InputHasConnectionOrValue(node.In16)
	newInputArray[16] = InputData{valueIn16, linkIn16}

	arraysMatch, arrayChanges := compareInputArrays(newInputArray, inst.inputsArray)
	for _, val := range arrayChanges {
		if val != nil {
			// log.Infof(fmt.Sprintf("FF Network EvaluateInputsArray() arrayChanges %d: %+v", f, arrayChanges[f]))
		}
	}

	priorityArrayWrite := make(map[string]*float64)

	if !arraysMatch || forceResend {
		for i, changeData := range arrayChanges {
			if i == 0 || (!forceResend && changeData == nil) {
				continue
			}
			inputName := fmt.Sprintf("_%d", i)
			if !forceResend && changeData.LinkDisconnected { // has the link been disconnected
				priorityArrayWrite[inputName] = float.New(1)
				priorityArrayWrite[inputName] = nil
			} else if newInputArray[i].LinkConnected && (forceResend || changeData.NewValue || changeData.LinkConnected) { // has there been a new value or a new link connected
				if forceResend && newInputArray[i].InputValue == nil {
					continue
				}
				priorityArrayWrite[inputName] = float.New(1)
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
		newValIsNil := inputData.InputValue == nil
		oldValIsNil := oldInputs[i].InputValue == nil
		newVal := float.NonNil(inputData.InputValue)
		oldVal := float.NonNil(oldInputs[i].InputValue)

		if (newValIsNil != oldValIsNil) || (!newValIsNil && !oldValIsNil && newVal != oldVal) || (inputData.LinkConnected != oldInputs[i].LinkConnected) {
			arraysMatch = false
			inputValueChanged := newValIsNil != oldValIsNil || !newValIsNil && !oldValIsNil && newVal != oldVal
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
