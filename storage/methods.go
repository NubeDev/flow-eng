package storage

type Storage interface {
	Close() error
	AddSettings(body *Settings) (*Settings, error)
	UpdateSettings(uuid string, body *Settings) (*Settings, error)
	DeleteSettings() error
	GetGitToken(uuid string, previewToken bool) (string, error)
	GetSetting(uuid string) (*Settings, error)
}
