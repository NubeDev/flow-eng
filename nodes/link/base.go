package link

import (
	log "github.com/sirupsen/logrus"
	"strings"
)

const (
	linkInputString  = "link-input-string"
	linkOutputString = "link-output-string"
	linkInputNum     = "link-input-number"
	linkOutputNum    = "link-output-number"
	linkInputBool    = "link-input-boolean"
	linkOutputBool   = "link-output-boolean"
	category         = "link"
)

var db *Store

func getStore() *Store {
	if db == nil {
		log.Error("link-node: the link node must be added")
		db = &Store{}
	}
	return db
}

func cleanName(s string) string {
	r := strings.Replace(s, "{parent.name}", "", -1)
	return r
}
