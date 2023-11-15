package ros

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/NubeDev/flow-eng/node"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func setError(body *node.Spec, message string) *node.Spec {
	body.SetStatusError(message)
	return body
}

var retryCount uint64 = 100

const (
	Category      = "rubix"
	rosNetwork    = "rubix-network"
	rosPoint      = "rubix-point"
	rosSchedule   = "rubix-schedule"
	rosPointWrite = "rubix-point-write"
)

const (
	pointError = "point deleted or never selected"
)

type covPayload struct {
	Value    *float64 `json:"value"`
	ValueRaw *float64 `json:"value_raw"`
	Ts       string   `json:"ts"`
	Priority *int     `json:"priority"`
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
	errorFetchPointMQTTConnect errorCode = "failed to connect to rubix-os"
	errorFailedFetchPoint      errorCode = "failed to fetch points list from rubix-os"
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

func parseCOV(body any) (payload *covPayload, value *float64, priority *int, err error) {
	// fmt.Println(fmt.Sprintf("FLOW POINT parseCOV() body: %+v", body))
	msg, ok := body.(mqtt.Message)
	if !ok {
		return nil, nil, nil, errors.New("failed to parse mqtt cov payload")
	}
	payload = &covPayload{}
	// fmt.Println(fmt.Sprintf("FLOW POINT parseCOV() msg.Payload(): %+v", string(msg.Payload())))
	err = json.Unmarshal(msg.Payload(), &payload)
	if err != nil {
		return nil, nil, nil, err
	}
	// fmt.Println(fmt.Sprintf("FLOW POINT parseCOV() payload: %+v", payload))
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
			break
		}
	}
	if !found {
		p.payloads = append(p.payloads, newPayload)
	}
	return p, found
}

type Schedule struct {
	Uuid              string    `json:"uuid"`
	Name              string    `json:"name"`
	Enable            bool      `json:"enable"`
	ThingClass        string    `json:"thing_class"`
	ThingType         string    `json:"thing_type"`
	TimeZone          string    `json:"timezone"`
	IsActive          bool      `json:"is_active"`
	ActiveWeekly      bool      `json:"active_weekly"`
	ActiveException   bool      `json:"active_exception"`
	ActiveEvent       bool      `json:"active_event"`
	Payload           float64   `json:"payload"`
	DefaultPayload    float64   `json:"default_payload"`
	PeriodStart       int64     `json:"period_start"`        // unix (seconds) timestamp
	PeriodStop        int64     `json:"period_stop"`         // unix (seconds) timestamp
	NextStart         int64     `json:"next_start"`          // unix (seconds) timestamp
	NextStop          int64     `json:"next_stop"`           // unix (seconds) timestamp
	PeriodStartString string    `json:"period_start_string"` // human readable timestamp
	PeriodStopString  string    `json:"period_stop_string"`  // human readable timestamp
	NextStartString   string    `json:"next_start_string"`   // human readable timestamp
	NextStopString    string    `json:"next_stop_string"`    // human readable timestamp
	CreatedOn         time.Time `json:"created_on"`
	UpdatedOn         time.Time `json:"updated_on"`
}
