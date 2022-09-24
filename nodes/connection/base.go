package connection

import log "github.com/sirupsen/logrus"

const (
	connectNode      = "connect"
	connectionInput  = "input"
	connectionOutput = "output"
	category         = "connection"
)

var db *Store

func getStore() *Store {
	if db == nil {
		log.Error("connection-node: the connection node must be added")
		db = &Store{}
	}
	return db
}
