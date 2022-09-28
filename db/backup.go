package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/tidwall/buntdb"
	"sort"
	"strings"
	"time"
)

type Backup struct {
	UUID        string      `json:"uuid"`
	UserComment string      `json:"user_comment"`
	Time        time.Time   `json:"time"`
	Timestamp   string      `json:"timestamp"`
	Data        interface{} `json:"data,omitempty"`
}

func matchBackupUUID(uuid string) bool {
	if len(uuid) == 16 {
		if uuid[0:4] == "flo_" {
			return true
		}
	}
	return false
}

func strip(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

func (inst *db) AddBackup(body *Backup) (*Backup, error) {

	body.Time = ttime.New().Now()
	body.Timestamp = body.Time.Format(time.RFC822)
	body.UUID = helpers.UUID("flo")
	body.UserComment = strip(body.UserComment)
	if body.UserComment == "" {
		body.UserComment = fmt.Sprintf("backup-%s", body.Timestamp)
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

func (inst *db) DeleteBackup(uuid string) error {
	if matchBackupUUID(uuid) {
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
	return errors.New("incorrect backup uuid")
}

func (inst *db) GetBackup(uuid string) (*Backup, error) {
	if matchBackupUUID(uuid) {
		var data *Backup
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
		return nil, errors.New("incorrect backup uuid")
	}
}

func (inst *db) GetLatestBackup() (*Backup, error) {
	backups, err := inst.GetBackups()
	if err != nil {
		return nil, err
	}
	if len(backups) > 0 {
		return &backups[0], nil
	}
	return nil, errors.New("no flow backups have been taken yet")
}

func (inst *db) GetBackups() ([]Backup, error) {
	var resp []Backup

	err := inst.DB.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			var data Backup
			err := json.Unmarshal([]byte(value), &data)
			if err != nil {
				return false
			}
			if matchBackupUUID(data.UUID) {
				resp = append(resp, data)
				//fmt.Printf("key: %s, value: %s\n", key, value)
			}
			return true
		})
		return err
	})
	if err != nil {
		fmt.Printf("Error: %s", err)
		return []Backup{}, err
	}
	sort.Slice(resp, func(i, j int) bool {
		return resp[i].Time.After(resp[j].Time)
	})
	return resp, nil
}
