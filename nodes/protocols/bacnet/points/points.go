package points

import (
	"github.com/NubeDev/flow-eng/helpers/names"
)

type Point struct {
	UUID                        string `json:"uuid"`
	Name                        string
	Application                 names.ApplicationName `json:"application"`
	ObjectType                  ObjectType            `json:"objectType"`
	ObjectID                    ObjectID
	IoType                      IoType
	IsIO                        bool // if it's an io-pin for a real device
	IsWriteable                 bool
	Enable                      bool
	ValueFromRead               float64
	PresentValue                float64
	Offset                      float64
	ScaleEnable                 bool
	ScaleInMin                  float64
	ScaleInMax                  float64
	ScaleOutMin                 float64
	ScaleOutMax                 float64
	Factor                      float64
	WriteValue                  *PriArray
	WriteValueFromBACnet        *PriArray
	PendingWriteValueFromBACnet bool
	PendingWriteCount           uint64
	PendingMQTTPublish          bool
	InError                     bool
	Message                     string
	ModbusDevAddr               int
	ModbusRegister              int
}

type ModbusPoints struct {
	DeviceOne   []*Point
	DeviceTwo   []*Point
	DeviceThree []*Point
	DeviceFour  []*Point
}

func (inst *Store) mergePriority(p2 *PriArray, in14, in15 *float64) *PriArray {
	if p2 == nil {
		p2 = &PriArray{
			P14: in14, // these are reversed for the flow
			P15: in15, // these are reversed for the flow
		}
		return p2
	}
	out := &PriArray{
		P1:  p2.P1,
		P2:  p2.P2,
		P3:  p2.P3,
		P4:  p2.P4,
		P5:  p2.P5,
		P6:  p2.P6,
		P7:  p2.P7,
		P8:  p2.P8,
		P9:  p2.P9,
		P10: p2.P10,
		P11: p2.P11,
		P12: p2.P12,
		P13: p2.P13,
		P14: in14, // these are reversed for the flow
		P15: in15, // these are reversed for the flow
		P16: p2.P16,
	}
	return out
}
