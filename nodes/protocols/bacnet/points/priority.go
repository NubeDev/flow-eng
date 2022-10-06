package points

import (
	"encoding/json"
	"errors"
	"github.com/NubeDev/flow-eng/helpers/float"
	"github.com/NubeDev/flow-eng/helpers/topics"
	"strconv"
	"strings"
)

type BacnetPayload struct {
	Value string `json:"value"`
	UUID  string `json:"uuid"`
}

type Payload struct {
	topic       string
	value       *float64
	priAndValue *priAndValue
	priArray    *PriArray
	objectID    ObjectID
	objectType  ObjectType
}

func NewPayload() *Payload {
	return &Payload{}
}

// outputs are both pub/sub to the server
// inputs don't need to subscribe only publish to the server

// bacnet/program/2508/state -> "started" on start of the bacnet-server
// bacnet/ao/1/pv -> "11.000000"
// bacnet/ao/1/pri -> {Null,33.3,Null,Null,Null,Null,Null,Null,Null,Null,Null,Null,Null,Null,Null,11.000000}

func (inst *Payload) NewMessage(msg interface{}) error {
	m, ok := msg.(*topics.Message)
	if ok {
		msgString := string(m.Msg.Payload())
		var decoded *BacnetPayload
		err := json.Unmarshal(m.Msg.Payload(), &decoded)
		if decoded == nil || err != nil {
			return errors.New("mqtt bacnet message failed to decode")
		}
		if decoded.UUID != "" {
			return errors.New("mqtt bacnet message came from flow (stopping infinite loop)")
		}
		topic := m.Msg.Topic()
		inst.topic = topic
		id, err := objectId(topic)
		if err != nil {
			return err
		}
		if object(topic) == "" {
			return errors.New("bacnet-message: failed to get object-type from")
		}
		inst.objectID = ObjectID(id)
		inst.objectType = object(topic)
		if topics.IsPV(inst.topic) {
			v, err := strconv.ParseFloat(msgString, 64)
			if err != nil {
				inst.value = float.New(v)
			}
		}
		if topics.IsPri(inst.topic) {
			inst.priArray = inst.cleanArray(msgString)
			inst.priAndValue = GetHighest(inst.priArray)
		}
	} else {
		return errors.New("bacnet-message: failed to decode message")
	}
	return nil
}

func NewPriArrayAt15(value float64) *PriArray {
	return &PriArray{
		P15: float.New(value),
	}
}

func NewPriArray(in14, in15 *float64) *PriArray {
	return &PriArray{
		P14: in14,
		P15: in15,
	}
}

type Priority struct {
	P1  *float64 `json:"_1,omitempty"`
	P2  *float64 `json:"_2,omitempty"`
	P3  *float64 `json:"_3,omitempty"`
	P4  *float64 `json:"_4,omitempty"`
	P5  *float64 `json:"_5,omitempty"`
	P6  *float64 `json:"_6,omitempty"`
	P7  *float64 `json:"_7,omitempty"`
	P8  *float64 `json:"_8,omitempty"`
	P9  *float64 `json:"_9,omitempty"`
	P10 *float64 `json:"_10,omitempty"`
	P11 *float64 `json:"_11,omitempty"`
	P12 *float64 `json:"_12,omitempty"`
	P13 *float64 `json:"_13,omitempty"`
	P14 *float64 `json:"_14,omitempty"`
	P15 *float64 `json:"_15,omitempty"`
	P16 *float64 `json:"_16,omitempty"`
}

type PriArray struct {
	P1  *float64 `json:"_1"`
	P2  *float64 `json:"_2"`
	P3  *float64 `json:"_3"`
	P4  *float64 `json:"_4"`
	P5  *float64 `json:"_5"`
	P6  *float64 `json:"_6"`
	P7  *float64 `json:"_7"`
	P8  *float64 `json:"_8"`
	P9  *float64 `json:"_9"`
	P10 *float64 `json:"_10"`
	P11 *float64 `json:"_11"`
	P12 *float64 `json:"_12"`
	P13 *float64 `json:"_13"`
	P14 *float64 `json:"_14"`
	P15 *float64 `json:"_15"`
	P16 *float64 `json:"_16"`
}

func set(part string) *float64 {
	if part == "Null" {
		return nil
	} else {
		f, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil
		}
		return float.New(f)
	}
}

type priAndValue struct {
	Number int
	Value  float64
}

func getHighest(num int, val *float64) *priAndValue {
	return &priAndValue{
		Number: num,
		Value:  float.NonNil(val),
	}
}

func (inst *Payload) GetObjectID() (ObjectType, ObjectID) {
	return inst.objectType, inst.objectID
}

func (inst *Payload) GetTopic() string {
	return inst.topic
}

func (inst *Payload) GetPresentValue() *float64 {
	return inst.value
}

func (inst *Payload) GetFullPriority() *PriArray {
	return inst.priArray
}

func (inst *Payload) GetHighestPriority() *priAndValue {
	return inst.priAndValue
}

// GetWriteArrayValues get 1n14, 1n15
func GetWriteArrayValues(payload *PriArray) (in14, in15 *float64) {
	if payload == nil {
		payload = &PriArray{}
	}
	if payload.P14 != nil {
		in14 = payload.P4
	}
	if payload.P15 != nil {
		in15 = payload.P15
	}
	return in14, in15

}

func GetHighest(payload *PriArray) *priAndValue {
	if payload == nil {
		payload = &PriArray{}
	}
	if payload.P1 != nil {
		return getHighest(1, payload.P1)
	}
	if payload.P2 != nil {
		return getHighest(2, payload.P2)
	}
	if payload.P3 != nil {
		return getHighest(3, payload.P3)
	}
	if payload.P4 != nil {
		return getHighest(4, payload.P4)
	}
	if payload.P5 != nil {
		return getHighest(5, payload.P5)
	}
	if payload.P6 != nil {
		return getHighest(6, payload.P6)
	}
	if payload.P7 != nil {
		return getHighest(7, payload.P7)
	}
	if payload.P8 != nil {
		return getHighest(8, payload.P8)
	}
	if payload.P9 != nil {
		return getHighest(9, payload.P9)
	}
	if payload.P10 != nil {
		return getHighest(10, payload.P10)
	}
	if payload.P11 != nil {
		return getHighest(11, payload.P11)
	}
	if payload.P12 != nil {
		return getHighest(12, payload.P12)
	}
	if payload.P13 != nil {
		return getHighest(13, payload.P13)
	}
	if payload.P14 != nil {
		return getHighest(14, payload.P14)
	}
	if payload.P15 != nil {
		return getHighest(15, payload.P15)
	}
	if payload.P16 != nil {
		return getHighest(16, payload.P16)
	}
	return nil

}

func (inst *Payload) cleanArray(payload string) *PriArray {
	payload = strings.ReplaceAll(payload, "{", "")
	payload = strings.ReplaceAll(payload, "}", "")
	parts := strings.Split(payload, ",")
	if len(parts) != 16 {
		return nil
	}
	arr := &PriArray{
		P1:  set(parts[0]),
		P2:  set(parts[1]),
		P3:  set(parts[2]),
		P4:  set(parts[3]),
		P5:  set(parts[4]),
		P6:  set(parts[5]),
		P7:  set(parts[6]),
		P8:  set(parts[7]),
		P9:  set(parts[8]),
		P10: set(parts[9]),
		P11: set(parts[10]),
		P12: set(parts[11]),
		P13: set(parts[12]),
		P14: set(parts[13]),
		P15: set(parts[14]),
		P16: set(parts[15]),
	}
	return arr
}

//object bacnet/ao/1/pv
func object(topic string) ObjectType {
	parts := strings.Split(topic, "/")
	if len(parts) >= 2 {
		if parts[0] == "bacnet" {
			t := parts[1]
			if t == "ai" {
				return AnalogInput
			}
			if t == "ao" {
				return AnalogOutput
			}
			if t == "av" {
				return AnalogVariable
			}
			if t == "bi" {
				return BinaryInput
			}
			if t == "bo" {
				return BinaryOutput
			}
			if t == "bv" {
				return BinaryVariable
			}
		}
	}
	return ""
}

//object bacnet/ao/1/pv
func objectId(topic string) (int, error) {
	parts := strings.Split(topic, "/")
	if len(parts) >= 2 {
		if parts[0] == "bacnet" {
			if object(topic) != "" {
				return strconv.Atoi(parts[2])
			}
		}
	}
	return 0, errors.New("bacnet-message: failed to get bacnet object-id")
}
