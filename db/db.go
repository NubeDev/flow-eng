package db

import (
	"github.com/tidwall/buntdb"
	"log"
)

type db struct {
	DB *buntdb.DB
}

func New(dbFile string) DB {
	if dbFile == "" {
		dbFile = "./flow.db"
	}
	newDb, err := buntdb.Open(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	return &db{DB: newDb}
}

func (inst *db) Close() error {
	inst.DB.Close()
	return nil
}

func matchSettingsUUID(uuid string) bool {
	if len(uuid) == 16 {
		if uuid[0:4] == "set_" {
			return true
		}
	}
	return false
}
