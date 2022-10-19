package flow

import (
	"fmt"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
)

type nodeSchema struct {
	Conn schemas.EnumString `json:"connections"`
}

func (inst *Network) getConnectionsNames() (names []string, uuids []string) {
	db := inst.GetDB()
	if db != nil {
		connections, err := inst.GetDB().GetConnections()
		if err != nil {
			log.Errorf("flow-networks get connections err %s", err.Error())
			return nil, nil
		}
		for _, connection := range connections {
			name := fmt.Sprintf("name:%s ip:%s port:%d", connection.Name, connection.Host, connection.Port)
			names = append(names, name)
			uuids = append(uuids, connection.UUID)
		}
		return names, uuids
	}
	log.Errorf("flow-networks failed to get db instance")
	return nil, nil

}

func (inst *Network) buildSchema() *schemas.Schema {
	names, uuids := inst.getConnectionsNames()
	props := &nodeSchema{}
	props.Conn.Title = "connections"
	if len(names) > 0 {
		props.Conn.Default = names[0]
	}

	props.Conn.Options = uuids
	props.Conn.EnumName = names
	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "connections",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}

type nodeSettings struct {
	Conn string `json:"connections"`
}

func getSettings(body map[string]interface{}) (string, error) {
	settings := &nodeSettings{}
	err := mapstructure.Decode(body, settings)
	if err != nil {
		return "", err
	}
	if settings != nil {
		return settings.Conn, nil
	}
	return "", nil
}
