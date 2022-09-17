package storage

type Connection struct {
	UUID                          string
	Enabled                       *bool  `json:"enabled,omitempty"`
	Application                   string `json:"application"` // bacnet nodes or mqttbase nodes
	Name                          string `json:"name,omitempty"`
	Host                          string `json:"host,omitempty"`
	Port                          int    `json:"port,omitempty"`
	Authentication                *bool  `json:"authentication,omitempty"`
	HTTPS                         *bool  `json:"https,omitempty"`
	Username                      string `json:"username,omitempty"`
	Password                      string `json:"password,omitempty"`
	Token                         string `json:"token,omitempty"`
	Keepalive                     int    `json:"keepalive,omitempty"`
	Qos                           int    `json:"qos,omitempty"`
	Retain                        *bool  `json:"retain,omitempty"`
	AttemptReconnectOnUnavailable *bool  `json:"attemptReconnectOnUnavailable,omitempty"`
	AttemptReconnectSecs          int    `json:"attemptReconnectSecs,omitempty"`
	Timeout                       int    `json:"timeout,omitempty"`
}
