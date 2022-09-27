package points

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/names"
	log "github.com/sirupsen/logrus"
)

type Store struct {
	Application names.ApplicationName `json:"application"`
	Store       *ObjectStore          `json:"store"`
	Points      []*Point              `json:"points"`
}

type pointAllowance struct {
	Object ObjectType
	From   int
	Count  int
}

type ToBacnet struct {
	CovEvent    bool
	ToBacnet    float64
	ToBacnetPri int
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

type ObjectStore struct {
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

// if modbus and rubix-io modbus will still start and 1 and then the last modbus addr,
// is where the rubix addr will start (but if the user add a new modbus device the rubix-io address's will be push back)
var (
	rubixUICount = 8
	rubixUOCount = 6 // 6UOs and DO1, DO2
	rubixDOCount = 2
)

var (
	io16UICount = 8
	io16UOCount = 8
)

var (
	calculatedUICount = 0
	calculatedUOCount = 0
	calculatedDICount = 0
	calculatedDOCount = 0
)

func CalcPointCount(deviceCount int, app names.ApplicationName) (rubixUIStart, rubixUOStart ObjectID) {
	if deviceCount == 0 {
		deviceCount = 1
	}
	if app == names.Edge {
		return calcModbusRubix(deviceCount, false, false, true)
	}
	if app == names.Modbus {
		return calcModbusRubix(deviceCount, true, false, false)
	}
	if app == names.RubixIOAndModbus {
		return calcModbusRubix(deviceCount, true, true, false)
	}
	if app == names.RubixIO {
		return calcModbusRubix(deviceCount, false, true, false)
	}
	return 0, 0

}

func calcModbusRubix(deviceCount int, isModbus, isRubix, isEdge bool) (rubixUIStart, rubixUOStart ObjectID) {
	if isEdge { // edge
		calculatedUICount = edgeUICount
		calculatedUOCount = edgeUOCount
		calculatedDICount = edgeDICount
		calculatedDOCount = edgeDOCount
		log.Infof(" calculated bacnet point EDGE-28 -> calculatedUICount:%d, calculatedUOCount:%d, calculatedDICount:%d, calculatedDOCount:%d,", calculatedUICount, calculatedUOCount, calculatedDICount, calculatedDOCount)
		return 0, 0
	}
	if isRubix && !isModbus { // just rubix
		calculatedUICount = rubixUICount
		calculatedUOCount = rubixUOCount
		calculatedDOCount = rubixDOCount
		log.Infof(" calculated bacnet point RUBIX-IO -> calculatedUICount:%d, calculatedUOCount:%d, calculatedDICount:%d, calculatedDOCount:%d,", calculatedUICount, calculatedUOCount, calculatedDICount, calculatedDOCount)
		return 1, 1
	}

	if !isRubix && isModbus { // just modbus
		calculatedUICount = io16UICount * deviceCount
		calculatedUOCount = io16UOCount * deviceCount
		log.Infof(" calculated bacnet point MODBUS -> calculatedUICount:%d, calculatedUOCount:%d, calculatedDICount:%d, calculatedDOCount:%d,", calculatedUICount, calculatedUOCount, calculatedDICount, calculatedDOCount)
		return 0, 0
	}

	if isRubix && isModbus { // rubix & modbus
		calculatedUICount = io16UICount*deviceCount + rubixUICount
		calculatedUOCount = io16UOCount*deviceCount + rubixUOCount
		calculatedDOCount = rubixDOCount
		log.Infof(" calculated bacnet point MODBUS & RUBIX-IO -> calculatedUICount:%d, calculatedUOCount:%d, calculatedDICount:%d, calculatedDOCount:%d,", calculatedUICount, calculatedUOCount, calculatedDICount, calculatedDOCount)
		return ObjectID(calculatedUICount - 7), ObjectID(calculatedUOCount - 7)
	}

	return 0, 0
}

func New(app names.ApplicationName, pStore *ObjectStore, deviceCount, avAllowance, bvAllowance int) *Store {
	bacnetStore := &Store{
		Application: app,
	}
	if pStore == nil {
		pStore = &ObjectStore{}
	}

	CalcPointCount(deviceCount, app)

	ai := pStore.AI
	ao := pStore.AO
	av := pStore.AV
	bi := pStore.BI
	bo := pStore.BO
	bv := pStore.BV

	if av == nil {
		av = &AVStore{
			pointAllowance: pointAllowance{
				Object: AnalogVariable,
				From:   1,
				Count:  avAllowance,
			},
		}
	}
	if bv == nil {
		bv = &BVStore{
			pointAllowance: pointAllowance{
				Object: BinaryVariable,
				From:   1,
				Count:  bvAllowance,
			},
		}
	}
	if ai == nil {
		ai = &AIStore{
			pointAllowance: pointAllowance{
				Object: AnalogInput,
				From:   1,
				Count:  calculatedUICount,
			},
		}
		if ao == nil {
			ao = &AOStore{
				pointAllowance: pointAllowance{
					Object: AnalogOutput,
					From:   1,
					Count:  calculatedUOCount,
				},
			}
		}
		if bo == nil {
			bo = &BOStore{
				pointAllowance: pointAllowance{
					Object: BinaryOutput,
					From:   1,
					Count:  calculatedDOCount,
				},
			}
		}
		if bi == nil {
			bi = &BIStore{
				pointAllowance: pointAllowance{
					Object: BinaryInput,
					From:   1,
					Count:  calculatedDICount,
				},
			}
		}
	}
	store := &ObjectStore{
		AI: ai,
		AO: ao,
		AV: av,
		BI: bi,
		BO: bo,
		BV: bv,
	}
	bacnetStore.Store = store
	return bacnetStore
}

func (inst *Store) GetStore() *ObjectStore {
	return inst.Store
}

func (inst *Store) AddPoint(point *Point, ignoreError bool) (*Point, error) {
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
	log.Infof("bacnet-add-point type-%s:%d", point.ObjectType, point.ObjectID)

	if objectType == AnalogInput {
		checked = true
		p := inst.Store.AI
		err = errNoObj(p, objectType)
		if err != nil {
			return nil, err
		}
		err = inst.checkExisting(point, p.From, p.Count)
		if err != nil {
			if !ignoreError {
				return nil, err
			}
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
			if !ignoreError {
				return nil, err
			}
		}
	}

	if objectType == AnalogVariable {
		checked = true
		p := inst.Store.AV
		err = errNoObj(p, objectType)
		if err != nil {
			return nil, err
		}
		err = inst.checkExisting(point, p.From, p.Count)
		if err != nil {
			if !ignoreError {
				return nil, err
			}
		}
	}

	if objectType == BinaryInput {
		checked = true
		p := inst.Store.BI
		err = errNoObj(p, objectType)
		if err != nil {
			return nil, err
		}
		err = inst.checkExisting(point, p.From, p.Count)
		if err != nil {
			if !ignoreError {
				return nil, err
			}
		}
	}

	if objectType == BinaryOutput {
		checked = true
		p := inst.Store.BO
		err = errNoObj(p, objectType)
		if err != nil {
			return nil, err
		}
		err = inst.checkExisting(point, p.From, p.Count)
		if err != nil {
			if !ignoreError {
				return nil, err
			}
		}
	}

	if objectType == BinaryVariable {
		checked = true
		p := inst.Store.BV
		err = errNoObj(p, objectType)
		if err != nil {
			return nil, err
		}
		err = inst.checkExisting(point, p.From, p.Count)
		if err != nil {
			if !ignoreError {
				return nil, err
			}
		}
	}

	if !checked {
		return nil, errors.New(fmt.Sprintf("store-add-point: not type found for object type: %s", objectType))
	}
	inst.Points = append(inst.Points, point)
	return point, nil
}

func errNoObj(pnt interface{}, objectType ObjectType) error {
	if pnt == nil {
		return errors.New(fmt.Sprintf("store-add-point: the server does not support object type: %s", objectType))
	}
	return nil
}

func (inst *Store) checkExisting(point *Point, from, to int) error {

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

func (inst *Store) allowableCount(objectID, from, count int) error {
	if objectID > count { // is above what is allowed
		return errors.New(fmt.Sprintf("store-add-point: the allowable max object-id is: %d and the current is: %d", count, objectID))
	}
	return nil
}

func (inst *Store) GetApplication() names.ApplicationName {
	return inst.Application
}
