package broker

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/schemas"
	"github.com/NubeIO/lib-schema/schema"
	log "github.com/sirupsen/logrus"
)

type nodeSchema struct {
	Conn schemas.EnumString `json:"connections"`
}

const selectConnection = "Please add/select a MQTT connection"

func (inst *Broker) getConnectionsNames() (names []string, uuids []string) {
	names = append(names, selectConnection)
	uuids = append(uuids, selectConnection)
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
func (inst *Broker) buildSchema() *schemas.Schema {
	names, uuids := inst.getConnectionsNames()
	props := &nodeSchema{}
	props.Conn.Title = "connections"
	if len(names) > 0 {
		props.Conn.Default = names[0]
	} else {
		names = append(names, "no connection has been added")
	}
	props.Conn.Options = uuids
	props.Conn.EnumName = names
	schema.Set(props)
	s := &schemas.Schema{
		Schema: schemas.SchemaBody{
			Title:      "settings",
			Properties: props,
		},
		UiSchema: nil,
	}
	return s
}

type nodeSettings struct {
	Conn string `json:"connections"`
}

func getSettings(body map[string]interface{}) (*nodeSettings, error) {
	settings := &nodeSettings{}
	marshal, err := json.Marshal(body)
	if err != nil {
		return settings, err
	}
	err = json.Unmarshal(marshal, &settings)
	return settings, err
}
