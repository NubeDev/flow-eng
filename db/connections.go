package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/names"
	"github.com/tidwall/buntdb"
)

type Connection struct {
	UUID                          string                `json:"uuid"`
	Enabled                       *bool                 `json:"enabled,omitempty"`
	Application                   names.ApplicationName `json:"application"` // bacnet
	Name                          string                `json:"name,omitempty"`
	Host                          string                `json:"host,omitempty"`
	Port                          int                   `json:"port,omitempty"`
	Authentication                *bool                 `json:"authentication,omitempty"`
	HTTPS                         *bool                 `json:"https,omitempty"`
	Username                      string                `json:"username,omitempty"`
	Password                      string                `json:"password,omitempty"`
	Email                         string                `json:"gmail,omitempty"`
	Token                         string                `json:"token,omitempty"`
	Keepalive                     int                   `json:"keepalive,omitempty"`
	Qos                           int                   `json:"qos,omitempty"`
	Retain                        *bool                 `json:"retain,omitempty"`
	AttemptReconnectOnUnavailable *bool                 `json:"attemptReconnectOnUnavailable,omitempty"`
	AttemptReconnectSecs          int                   `json:"attemptReconnectSecs,omitempty"`
	Timeout                       int                   `json:"timeout,omitempty"`
}

func matchConnection(t names.ApplicationName) bool {
	switch t {
	case names.FlowFramework:
		return true
	case names.MQTT:
		return true
	case names.Email:
		return true
	}
	return false
}

func matchConnectionUUID(uuid string) bool {
	if len(uuid) == 16 {
		if uuid[0:4] == "con_" {
			return true
		}
	}
	return false
}

func (inst *db) AddConnection(body *Connection) (*Connection, error) {
	body.UUID = helpers.UUID("con")
	if !matchConnection(body.Application) {
		return nil, errors.New("application name does not match, try mqtt")
	}

	if body.Name == "" {
		return nil, errors.New("name can not be empty")
	}
	c, err := inst.GetConnections()
	for _, connection := range c {
		if connection.Name == body.Name {
			return nil, errors.New("existing name can not be empty")
		}
	}
	data, err := json.Marshal(body)
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

func (inst *db) GetConnectionByName(name string) (*Connection, error) {
	conn, err := inst.GetConnections()
	if err != nil {
		return nil, err
	}
	for _, conn := range conn {
		if conn.Name == name {
			return &conn, nil
		}
	}
	return nil, nil
}

func (inst *db) GetConnectionsByType(t names.ApplicationName) ([]Connection, error) {
	var resp []Connection
	conn, err := inst.GetConnections()
	if err != nil {
		return conn, err
	}
	for _, conn := range conn {
		if conn.Application == t {
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
