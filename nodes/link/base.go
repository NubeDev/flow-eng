package link

import log "github.com/sirupsen/logrus"

const (
	linkInput     = "link-input-string"
	linkOutput    = "link-output-string"
	linkInputNum  = "link-input-number"
	linkOutputNum = "link-output-number"
	category      = "link"
)

var db *Store

func getStore() *Store {
	if db == nil {
		log.Error("link-node: the link node must be added")
		db = &Store{}
	}
	return db
}
