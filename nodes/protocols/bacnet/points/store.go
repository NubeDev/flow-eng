package points

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/names"
	log "github.com/sirupsen/logrus"
)

type Store struct {
	Application       names.ApplicationName `json:"application"`
	Store             *ObjectStore          `json:"store"`
	Points            []*Point              `json:"points"`
	ModbusDeviceCount int                   `json:"modbusDeviceCount"`
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
	if app == names.Modbus {
		return calcModbusRubix(deviceCount)
	}
	return 0, 0

}

func calcModbusRubix(deviceCount int) (rubixUIStart, rubixUOStart ObjectID) {
	calculatedUICount = io16UICount * deviceCount
	calculatedUOCount = io16UOCount * deviceCount
	log.Infof(" calculated bacnet point MODBUS -> calculatedUICount:%d, calculatedUOCount:%d, calculatedDICount:%d, calculatedDOCount:%d,", calculatedUICount, calculatedUOCount, calculatedDICount, calculatedDOCount)
	return 0, 0

}

func New(app names.ApplicationName, pStore *ObjectStore, deviceCount, avAllowance, bvAllowance int) *Store {
	bacnetStore := &Store{
		Application: app,
	}
	if pStore == nil {
		pStore = &ObjectStore{}
	}
	bacnetStore.ModbusDeviceCount = deviceCount
	bacnetStore.Application = app

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
	log.Infof("bacnet-store AI:%d AO:%d AV:%d BI:%d BO:%d BV:%d", ai.Count, ao.Count, av.Count, bi.Count, bo.Count, bv.Count)
	return bacnetStore
}

func (inst *Store) GetStore() *ObjectStore {
	return inst.Store
}

func (inst *Store) AddPoint(point *Point, ignoreError bool, deviceAddr int) (*Point, error) {
	var err error
	if point == nil {
		return nil, errors.New(fmt.Sprintf("store-add-point: point can not be empty"))
	}
	point.UUID = helpers.ShortUUID()
	objectType := point.ObjectType
	if point.ObjectType == "" {
		return nil, errors.New(fmt.Sprintf("store-add-point: point objectType can not be empty"))
	}
	// rubixUIStart, rubixUOStart := CalcPointCount(inst.ModbusDeviceCount, inst.Application)
	if point.ObjectType == AnalogInput {
		point.ModbusDevAddr = deviceAddr
		point.Application = names.Modbus
	}
	if point.ObjectType == AnalogOutput {
		addr, _ := ModbusBuildOutput(point.IoType, point.ObjectID)
		point.ModbusDevAddr = deviceAddr
		if point.IoType == IoTypeDigital {
			point.ModbusRegister = addr.Volt
		}
		if point.IoType == IoTypeVolts {
			point.ModbusRegister = addr.Volt
		}
		point.Application = names.Modbus

	}

	var checked bool
	if objectType == AnalogInput {
		checked = true
	}
	if objectType == AnalogOutput {
		checked = true
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
		errMsg := fmt.Sprintf("store-add-point: not type found for object type: %s", objectType)
		log.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	log.Infof("bacnet-add-point type-%s:%d application-type %s", point.ObjectType, point.ObjectID, point.Application)
	inst.Points = append(inst.Points, point)

	return point, nil
}

func errNoObj(pnt interface{}, objectType ObjectType) error {
	if pnt == nil {
		errMsg := fmt.Sprintf("store-add-point: point cant not be empty")
		log.Error(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func (inst *Store) checkExisting(point *Point, from, to int) error {
	err := inst.allowableCount(int(point.ObjectID), from, to)
	if err != nil {
		return err
	}
	// check if there is a free address
	// err = inst.CheckExistingPointErr(point)
	if err != nil {
		return err
	}
	return nil
}

func (inst *Store) allowableCount(objectID, from, count int) error {
	if objectID > count { // is above what is allowed
		errMsg := fmt.Sprintf("store-add-point: the allowable max object-id is: %d and the current is: %d", count, objectID)
		log.Error(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func (inst *Store) GetApplication() names.ApplicationName {
	return inst.Application
}
