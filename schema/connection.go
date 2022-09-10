package schema

type Connection struct {
	UUID                          string `json:"uuid"`
	Name                          string `json:"name,omitempty"`
	Enabled                       *bool  `json:"enabled,omitempty"`
	Host                          string `json:"host,omitempty"`
	Port                          int    `json:"port,omitempty"`
	Authentication                *bool  `json:"authentication,omitempty"`
	HTTPS                         *bool  `json:"https,omitempty"`
	Username                      string `json:"username,omitempty"`
	Password                      string `json:"password,omitempty"`
	Token                         string `json:"token,omitempty"`
	Keepalive                     int    `json:"keepalive,omitempty"`
	Qos                           int    `json:"qos,omitempty"`
	Retain                        bool   `json:"retain,omitempty"`
	AttemptReconnectOnUnavailable bool   `json:"attemptReconnectOnUnavailable,omitempty"`
	AttemptReconnectSecs          int    `json:"attemptReconnectSecs,omitempty"`
	Timeout                       int    `json:"timeout,omitempty"`
}
