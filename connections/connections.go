package connections

import (
	"github.com/NubeDev/flow-eng/helpers/names"
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

type ConnectionIF interface {
	GetConnection(uuid string) (*Connection, error)
	GetConnectionByName(name string) (*Connection, error)
	GetConnections() ([]Connection, error)
}

var globalConnIF ConnectionIF

func SetGlobalConnectionIF(connIF ConnectionIF) {
	globalConnIF = connIF
}

func GetGlobalConnectionIF() ConnectionIF {
	return globalConnIF
}
