package bstore

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
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

type BacnetStore struct {
	Application node.ApplicationName `json:"application"`
	Store       *PointStore          `json:"store"`
	Points      []*Point             `json:"points"`
}

type pointAllowance struct {
	Object ObjectType
	From   int
	Count  int
}

type Point struct {
	UUID         string               `json:"uuid"`
	Application  node.ApplicationName `json:"application"`
	ObjectType   ObjectType           `json:"objectType"`
	ObjectID     ObjectID             `json:"ObjectID"`
	presentValue *float64
	priAndValue  *priAndValue
	priArray     *priArray
	//IoType      IoType               `json:"ioType"`
	//IoNumber    IoNumber             `json:"ioNumber"`
	//ReadValue   float64              `json:"readValue"`
	//WriteValue  float64              `json:"writeValue"`
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

type PointStore struct {
	AI *AIStore `json:"ai"`
	AO *AOStore `json:"ao"`
	AV *AVStore `json:"av"`
	BI *BIStore `json:"bi"`
	BO *BOStore `json:"bo"`
	BV *BVStore `json:"bv"`
}

var (
	edgeUICount = 8
	edgeUOCount = 8
	edgeDICount = 8
	edgeDOCount = 8 // 6DOs and r1, r2
)

func New(app node.ApplicationName, pStore *PointStore) *BacnetStore {
	bacnetStore := &BacnetStore{
		Application: app,
	}
	if pStore == nil {
		pStore = &PointStore{}
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
	store := &PointStore{
		AI: ai,
		AV: nil,
	}
	bacnetStore.Store = store
	return bacnetStore
}

func (inst *BacnetStore) GetStore() *PointStore {
	return inst.Store
}

func (inst *BacnetStore) GetStoreByType(objectType ObjectType) *PointStore {
	return inst.Store
}

func (inst *BacnetStore) AddPoint(point *Point) (*Point, error) {
	var err error
	if point == nil {
		return nil, errors.New(fmt.Sprintf("store-add-point: point can not be empty"))
	}
	point.UUID = helpers.ShortUUID()
	objectType := point.ObjectType
	if point.ObjectType == "" {
		if point == nil {
			return nil, errors.New(fmt.Sprintf("store-add-point: point objectType can not be empty"))
		}
	}
	var checked bool

	if objectType == AnalogInput {
		checked = true
		p := inst.Store.AI
		err = errNoObj(p, objectType)
		if err != nil {
			return nil, err
		}
		err = inst.checkExisting(point, p.From, p.Count)
		if err != nil {
			return nil, err
		}
	}
	if objectType == AnalogOutput {
		checked = true
		p := inst.Store.AO
		err = errNoObj(p, objectType)
		if err != nil {
			return nil, err
		}
		err = inst.checkExisting(point, p.From, p.Count)
		if err != nil {
			return nil, err
		}
	}

	if objectType == AnalogVariable {
		checked = true
		p := inst.Store.AO
		err = errNoObj(p, objectType)
		if err != nil {
			return nil, err
		}
		err = inst.checkExisting(point, p.From, p.Count)
		if err != nil {
			return nil, err
		}
	}

	if !checked {
		return nil, errors.New(fmt.Sprintf("store-add-point: not type found for object type:%s", objectType))
	}
	inst.Points = append(inst.Points, point)
	return point, nil
}

func errNoObj(pnt interface{}, objectType ObjectType) error {
	if pnt == nil {
		return errors.New(fmt.Sprintf("store-add-point: the server does not support object type:%s", objectType))
	}
	return nil
}

func (inst *BacnetStore) checkExisting(point *Point, from, to int) error {
	err := inst.allowableCount(int(point.ObjectID), from, to)
	if err != nil {
		return err
	}
	// check if there is a free address
	err = inst.CheckExistingPointErr(point)
	if err != nil {
		return err
	}
	return nil
}

func (inst *BacnetStore) allowableCount(objectID, from, count int) error {
	to := from + count - 1
	if objectID > to { // is above what is allowed
		return errors.New(fmt.Sprintf("store-add-point: the allwoable max object-id is:%d and the current is:%d", to, objectID))
	}
	return nil
}

func (inst *BacnetStore) GetApplication() node.ApplicationName {
	return inst.Application
}

func (inst *BacnetStore) GetPoints() []*Point {
	return inst.Points
}

func (inst *BacnetStore) GetPointByObject(t ObjectType, id ObjectID) *Point {
	for _, point := range inst.GetPoints() {
		if point.ObjectType == t {
			if point.ObjectID == id {
				return point
			}
		}
	}
	return nil
}

func (inst *BacnetStore) GetPoint(uuid string) *Point {
	for _, point := range inst.GetPoints() {
		if point.UUID == uuid {
			return point
		}
	}
	return nil
}

//presentValue *float64
//priAndValue  *priAndValue
//priArray     *priArray

func (inst *BacnetStore) ReadPresentValue(uuid string) (*float64, bool) {
	p := inst.GetPoint(uuid)
	if p != nil {
		return p.presentValue, true
	}
	return nil, false
}

func (inst *BacnetStore) WritePointValue(uuid string, value *float64) bool {
	p := inst.GetPoint(uuid)
	if p != nil {
		p.presentValue = value
		return true
	}
	return false
}

func (inst *BacnetStore) GetByType(objectType ObjectType) (out []*Point, count int) {
	out = []*Point{}
	for _, pnt := range inst.GetPoints() {
		if pnt.ObjectType == objectType {
			out = append(out, pnt)
		}
	}
	return out, len(out)
}

func (inst *BacnetStore) CheckExistingPointErr(point *Point) error {
	if inst.CheckExistingPoint(point) {
		return errors.New(fmt.Sprintf("store-add-point: point is existing object-type:%s:%d", point.ObjectType, point.ObjectID))
	}
	return nil
}

func (inst *BacnetStore) CheckExistingPoint(point *Point) bool {
	for _, pnt := range inst.GetPoints() {
		if pnt.ObjectType == point.ObjectType {
			if pnt.ObjectID == point.ObjectID {
				return true
			}
		}
	}
	return false
}
