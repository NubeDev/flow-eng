package link

import log "github.com/sirupsen/logrus"

const (
	linkInput  = "link-input"
	linkOutput = "link-output"
	category   = "link"
)

var db *Store

func getStore() *Store {
	if db == nil {
		log.Error("link-node: the link node must be added")
		db = &Store{}
	}
	return db
}
