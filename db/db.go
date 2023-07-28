package db

import (
	"github.com/tidwall/buntdb"
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
		return &db{DB: newDb}
	}
	size := 10 * 1024 * 1024 // 10mb
	c := buntdb.Config{
		SyncPolicy:           buntdb.EverySecond,
		AutoShrinkPercentage: 30,
		AutoShrinkMinSize:    size,
	}
	err = newDb.SetConfig(c)
	if err != nil {
		return &db{DB: newDb}
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
