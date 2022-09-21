package points

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
)

type Store struct {
	Application node.ApplicationName `json:"application"`
	Store       *ObjectStore         `json:"store"`
	Points      []*Point             `json:"points"`
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

func New(app node.ApplicationName, pStore *ObjectStore) *Store {
	bacnetStore := &Store{
		Application: app,
	}
	if pStore == nil {
		pStore = &ObjectStore{}
	}

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

	if app == applications.Edge || app == applications.RubixIO || app == applications.Modbus {
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
		if app == applications.Edge || app == applications.Modbus {
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

func (inst *Store) GetStoreByType(objectType ObjectType) *ObjectStore {
	return inst.Store
}

func (inst *Store) AddPoint(point *Point) (*Point, error) {
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
		p := inst.Store.AV
		err = errNoObj(p, objectType)
		if err != nil {
			return nil, err
		}
		err = inst.checkExisting(point, p.From, p.Count)
		if err != nil {
			return nil, err
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
			return nil, err
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
			return nil, err
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
	to := from + count - 1
	if objectID > to { // is above what is allowed
		return errors.New(fmt.Sprintf("store-add-point: the allwoable max object-id is:%d and the current is:%d", to, objectID))
	}
	return nil
}

func (inst *Store) GetApplication() node.ApplicationName {
	return inst.Application
}
