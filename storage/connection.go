package storage

import "errors"

type Adapter interface {
	Add(body *Connection) (*Connection, error)
	Get(uuid string) (*Connection, error)
}

func NewAdapter(db Storage) Adapter {
	return &Connection{db: db}
}

func (inst *Connection) Add(body *Connection) (*Connection, error) {
	if body == nil {
		return nil, errors.New("connection body can not be empty")
	}
	if !matchConnection(body.Application) {
		return nil, errors.New("application type type did not match try, mqtt")
	}
	return inst.db.AddConnection(body)
}

func (inst *Connection) Get(uuid string) (*Connection, error) {
	return inst.db.GetConnection(uuid)
}

func matchConnection(t ConnType) bool {
	switch t {
	case ConnMQTT:
		return true
	case ConnFlowFramework:
		return true
	}
	return false
}
