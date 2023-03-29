package db

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers"
	"github.com/tidwall/buntdb"
)

type AutoRefresh struct {
	Enable bool `json:"enable"`
	Rate   int  `json:"rate"`
}

type Settings struct {
	UUID              string `json:"uuid"`
	GitToken          string `json:"git_token"`
	AutoRefreshEnable bool   `json:"auto_refresh_enable"`
	AutoRefreshRate   int    `json:"auto_refresh_rate"`
}

func (inst *db) AddSettings(body *Settings) (*Settings, error) {
	settings, err := inst.GetSettings()
	if err != nil {
		return nil, err
	}
	if len(settings) > 0 {
		if settings != nil {
			return nil, errors.New(fmt.Sprintf("settings can only be added once uuid: %s", settings[0].UUID))
		}
		return nil, errors.New("settings can only be added once")
	}
	body.UUID = helpers.UUID("set")
	if body.GitToken != "" {
		body.GitToken = encodeToken(body.GitToken)
	}
	if body.AutoRefreshRate == 0 {
		body.AutoRefreshRate = 200
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

func (inst *db) UpdateSettings(uuid string, body *Settings) (*Settings, error) {
	uuid_ := uuid
	if body.GitToken != "" {
		body.GitToken = encodeToken(body.GitToken)
	}
	body.UUID = uuid
	j, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	err = inst.DB.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(uuid_, string(j), nil)
		return err
	})
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}
	return inst.GetSetting(uuid)
}

func (inst *db) DeleteSettings() error {
	settings, err := inst.GetSettings()
	if err != nil {
		return err
	}
	for _, setting := range settings {
		inst.DeleteSetting(setting.UUID)
	}
	return nil
}

func (inst *db) DeleteSetting(uuid string) error {
	if matchSettingsUUID(uuid) {
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
	return errors.New("incorrect setting uuid")
}

func (inst *db) GetGitToken(uuid string, previewToken bool) (string, error) {
	uuid_ := uuid
	var data *Settings
	err := inst.DB.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(uuid_)
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
		return "", err
	}
	if data.GitToken != "" {
		// conversions.GitToken = conversions.GitToken
		data.GitToken = decodeToken(data.GitToken)
	}
	if previewToken {
		if len(data.GitToken) > 5 {
			return fmt.Sprintf("%s.....", data.GitToken[0:5]), nil
		} else {
			return fmt.Sprintf("token....."), nil
		}
	}
	return data.GitToken, nil
}

func (inst *db) GetSettings() ([]Settings, error) {
	var resp []Settings
	err := inst.DB.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			var data Settings
			err := json.Unmarshal([]byte(value), &data)
			if err != nil {
				return false
			}
			if matchSettingsUUID(key) {
				resp = append(resp, data)
			}
			return true
		})
		return err
	})
	if err != nil {
		fmt.Printf("Error: %s", err)
		return []Settings{}, err
	}
	return resp, nil
}

func (inst *db) GetFirstSettings() (*Settings, error) {
	s, err := inst.GetSettings()
	if err != nil {
		return nil, err
	}
	if len(s) > 0 {
		return &s[0], nil
	}
	return nil, nil
}

func (inst *db) GetSetting(uuid string) (*Settings, error) {
	settings, err := inst.GetSettings()
	if err != nil {
		return nil, err
	}
	if len(settings) == 0 {
		return nil, errors.New("no settings have been added")
	}
	uuid_ := uuid
	if matchSettingsUUID(uuid_) {
		var data *Settings
		err := inst.DB.View(func(tx *buntdb.Tx) error {
			val, err := tx.Get(uuid_)
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
		return nil, errors.New("incorrect settings uuid")
	}
}

func encodeToken(token string) string {
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func decodeToken(token string) string {
	data, _ := base64.StdEncoding.DecodeString(token)
	return string(data)
}
