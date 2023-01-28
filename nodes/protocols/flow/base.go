package flow

import (
	"encoding/json"
	"errors"
	"github.com/NubeDev/flow-eng/node"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func setError(body *node.Spec, message string) *node.Spec {
	body.SetStatusError(message)
	return body
}

const (
	category       = "flow"
	flowNetwork    = "flow-network"
	flowPoint      = "flow-point"
	flowSchedule   = "flow-schedule"
	flowPointWrite = "flow-point-write"
)

const (
	pointError = "point deleted or never selected"
)

type covPayload struct {
	Value    float64 `json:"value"`
	ValueRaw float64 `json:"value_raw"`
	Ts       string  `json:"ts"`
	Priority int     `json:"priority"`
}

type PointWriter struct {
	Priority *map[string]*float64 `json:"priority"`
}

// MqttPoint body for getting points from FF over mqtt (can get by name's or uuid, publish on topic rubix/platform/list/points)
type MqttPoint struct {
	NetworkName string       `json:"network_name,omitempty"`
	DeviceName  string       `json:"device_name,omitempty"`
	PointName   string       `json:"point_name,omitempty"`
	PointUUID   string       `json:"point_uuid,omitempty"`
	Priority    *PointWriter `json:"priority,omitempty"`
}

type errorCode string

const (
	errorOk                    errorCode = ""
	errorMQTTClientEmpty       errorCode = "failed to create mqtt client"
	errorFetchPointMQTTConnect errorCode = "failed to connect to flow-framework"
	errorFailedFetchPoint      errorCode = "failed to fetch points list from flow-framework"
)

type point struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type pointStore struct {
	parentID  string
	points    []*point
	payloads  []*pointDetails
	schedules []*Schedule
}

type pointDetails struct {
	nodeUUID       string
	topic          string
	netDevPntNames string
	pointUUID      string
	payload        string
	isWriteable    bool
}

func parseCOV(body any) (payload *covPayload, value float64, priority int, err error) {
	msg, ok := body.(mqtt.Message)
	if !ok {
		return nil, 0, 0, errors.New("failed to parse mqtt cov payload")
	}
	payload = &covPayload{}
	err = json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		return nil, 0, 0, err
	}
	return payload, payload.Value, payload.Priority, nil
}

func getPayloads(children interface{}, ok bool) []*pointDetails {
	if ok {
		mqttData := children.(*pointStore)
		if mqttData != nil {
			return mqttData.payloads
		}
	}
	return nil
}

func addUpdatePayload(nodeUUID string, p *pointStore, newPayload *pointDetails) (data *pointStore, found bool) {
	for i, payload := range p.payloads {
		if payload.nodeUUID == nodeUUID {
			p.payloads[i] = newPayload
			found = true
		}
	}
	if !found {
		p.payloads = append(p.payloads, newPayload)
	}
	return p, found
}

type Schedule struct {
	Uuid            string    `json:"uuid"`
	Name            string    `json:"name"`
	Enable          bool      `json:"enable"`
	ThingClass      string    `json:"thing_class"`
	ThingType       string    `json:"thing_type"`
	IsActive        bool      `json:"is_active"`
	IsGlobal        bool      `json:"is_global"`
	TimeZone        string    `json:"timezone"`
	ActiveWeekly    bool      `json:"active_weekly"`
	ActiveException bool      `json:"active_exception"`
	ActiveEvent     bool      `json:"active_event"`
	CreatedOn       time.Time `json:"created_on"`
	UpdatedOn       time.Time `json:"updated_on"`
}
