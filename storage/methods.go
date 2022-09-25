package storage

type Storage interface {
	Close() error
	AddSettings(body *Settings) (*Settings, error)
	UpdateSettings(uuid string, body *Settings) (*Settings, error)
	DeleteSettings() error
	GetGitToken(uuid string, previewToken bool) (string, error)
	GetSetting(uuid string) (*Settings, error)
	AddConnection(body *Connection) (*Connection, error)
	UpdateConnection(uuid string, body *Connection) (*Connection, error)
	GetConnections() ([]Connection, error)
	GetConnection(uuid string) (*Connection, error)
	DeleteConnection(uuid string) error
}
