package bacnet

import (
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
)

/*
| Store      |         |       |         |       |
| ---------- | ------- | ----- | ------- | ----- |
| application| AI from | AI to | AV from | AV to |
| bacnet     | 0       | 0     | 1       | 100   |
| edge       | 1       | 8     | 0       | 0     |
| rubix-io   | 1       | 8     | 0       | 0     |
*/

type ObjectID int
type ObjectType string
type IoNumber int  // 1, 2
type IoPort string // UI, UO eg: IoType:IoNumber -> UI1, UI2
type IoType string // digital

const AnalogInput ObjectType = "analogInput"
const AnalogOutput ObjectType = "analogOutput"
const AnalogVariable ObjectType = "analogVariable"

const BinaryInput ObjectType = "binaryInput"
const BinaryOutput ObjectType = "binaryInput"
const BinaryVariable ObjectType = "binaryVariable"

type pointAllowance struct {
	Object ObjectType
	From   int
	Count  int
}

type Point struct {
	Application node.ApplicationName `json:"application"`
	ObjectType  ObjectType           `json:"objectType"`
	ObjectID    ObjectID             `json:"ObjectID"`
	IoType      IoType               `json:"ioType"`
	IoNumber    IoNumber             `json:"ioNumber"`
}

type AIStore struct {
	pointAllowance
}

type AOStore struct {
	pointAllowance
}

type AVStore struct {
	pointAllowance
}

type BIStore struct {
	pointAllowance
}

type BOStore struct {
	pointAllowance
}

type BVStore struct {
	pointAllowance
}

type pointStore struct {
	AI *AIStore `json:"ai"`
	AO *AOStore `json:"ao"`
	AV *AVStore `json:"av"`
	BI *BIStore `json:"bi"`
	BO *BOStore `json:"bo"`
	BV *BVStore `json:"bv"`
}

const (
	edgeUICount = 8
	edgeUOCount = 8
	edgeDICount = 8
	edgeDOCount = 8 // 6DOs and r1, r2
)

func NewStore(app node.ApplicationName, pStore *pointStore) *pointStore {
	if pStore == nil {
		pStore = &pointStore{}
	}
	ai := pStore.AI
	ao := pStore.AO
	av := pStore.AV
	bi := pStore.BI
	bo := pStore.BO
	bv := pStore.BV

	if app == applications.BACnet { // bacnet-server
		if av == nil {
			av = &AVStore{
				pointAllowance: pointAllowance{
					Object: AnalogVariable,
					From:   1,
					Count:  200,
				},
			}
		}
		if bv == nil {
			bv = &BVStore{
				pointAllowance: pointAllowance{
					Object: BinaryVariable,
					From:   1,
					Count:  200,
				},
			}
		}
	}

	if app == applications.Edge || app == applications.RubixIO {
		if ai == nil {
			ai = &AIStore{
				pointAllowance: pointAllowance{
					Object: AnalogInput,
					From:   1,
					Count:  edgeUICount,
				},
			}
		}
		if ao == nil {
			ao = &AOStore{
				pointAllowance: pointAllowance{
					Object: AnalogOutput,
					From:   1,
					Count:  edgeUOCount,
				},
			}
		}
		if app == applications.Edge {
			if bo == nil {
				bo = &BOStore{
					pointAllowance: pointAllowance{
						Object: BinaryOutput,
						From:   1,
						Count:  edgeDOCount,
					},
				}
			}
			if bi == nil {
				bi = &BIStore{
					pointAllowance: pointAllowance{
						Object: BinaryInput,
						From:   1,
						Count:  edgeDICount,
					},
				}
			}
		}
	}
	ps := &pointStore{
		AI: ai,
		AV: nil,
	}
	return ps
}

func CheckExistingPoint(points []*Point, point *Point) bool {
	for _, pnt := range points {
		if pnt.ObjectType == point.ObjectType {
			if pnt.ObjectID == point.ObjectID {
				return true
			}
		}
	}
	return false
}
