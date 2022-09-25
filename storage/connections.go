package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/tidwall/buntdb"
)

func matchConnectionUUID(uuid string) bool {
	if len(uuid) == 16 {
		if uuid[0:4] == "con_" {
			return true
		}
	}
	return false
}

type ConnType string

const ConnMQTT ConnType = "mqtt"
const ConnFlowFramework ConnType = "flow-framework"

type Connection struct {
	UUID                          string
	Enabled                       *bool    `json:"enabled,omitempty"`
	Application                   ConnType `json:"application"` // bacnet
	Name                          string   `json:"name,omitempty"`
	Host                          string   `json:"host,omitempty"`
	Port                          int      `json:"port,omitempty"`
	Authentication                *bool    `json:"authentication,omitempty"`
	HTTPS                         *bool    `json:"https,omitempty"`
	Username                      string   `json:"username,omitempty"`
	Password                      string   `json:"password,omitempty"`
	Token                         string   `json:"token,omitempty"`
	Keepalive                     int      `json:"keepalive,omitempty"`
	Qos                           int      `json:"qos,omitempty"`
	Retain                        *bool    `json:"retain,omitempty"`
	AttemptReconnectOnUnavailable *bool    `json:"attemptReconnectOnUnavailable,omitempty"`
	AttemptReconnectSecs          int      `json:"attemptReconnectSecs,omitempty"`
	Timeout                       int      `json:"timeout,omitempty"`
	db                            Storage
}

func (inst *db) AddConnection(body *Connection) (*Connection, error) {
	data, err := json.Marshal(body)
	body.UUID = helpers.UUID("con")
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	err = inst.DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(body.UUID, string(data), nil)
		return err
	})
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	return body, nil
}

func (inst *db) UpdateConnection(uuid string, body *Connection) (*Connection, error) {
	j, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	err = inst.DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(uuid, string(j), nil)
		return err
	})
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	return body, nil
}

func (inst *db) GetConnectionsByType(t ConnType) ([]Connection, error) {
	var resp []Connection
	conn, err := inst.GetConnections()
	if err != nil {
		return conn, err
	}
	for _, connection := range conn {
		if connection.Application == t {
			var data Connection
			resp = append(resp, data)
		}
	}
	return resp, nil
}

func (inst *db) GetConnections() ([]Connection, error) {
	var resp []Connection
	err := inst.DB.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			var data Connection
			err := json.Unmarshal([]byte(value), &data)
			if err != nil {
				return false
			}
			if matchConnectionUUID(data.UUID) {
				resp = append(resp, data) // put into arry
			}
			return true
		})
		return err
	})
	if err != nil {
		fmt.Printf("Error: %s", err)
		return []Connection{}, err
	}
	return resp, nil
}

func (inst *db) GetConnection(uuid string) (*Connection, error) {
	if matchConnectionUUID(uuid) {
		var data *Connection
		err := inst.DB.View(func(tx *buntdb.Tx) error {
			val, err := tx.Get(uuid)
			if err != nil {
				return err
			}
			err = json.Unmarshal([]byte(val), &data)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error: %s", err)
			return nil, err
		}
		return data, nil
	} else {
		return nil, errors.New("incorrect connection uuid")
	}
}

func (inst *db) DeleteConnection(uuid string) error {
	if matchConnectionUUID(uuid) {
		err := inst.DB.Update(func(tx *buntdb.Tx) error {
			_, err := tx.Delete(uuid)
			return err
		})
		if err != nil {
			fmt.Printf("Error delete: %s", err)
			return err
		}
		return nil
	}
	return errors.New("incorrect connection uuid")
}
