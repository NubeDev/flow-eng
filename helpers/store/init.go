package store

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type Store struct {
	Store *cache.Cache
}

//Init init store
func Init() *Store {
	newStore := cache.New(cache.NoExpiration, cache.DefaultExpiration)
	store := &Store{Store: newStore}
	return store
}

//All get everything from the store
func (inst *Store) All() map[string]cache.Item {
	return inst.Store.Items()
}

// Get an item from the store. Returns the item or nil, and a bool indicating
// whether the key was found.
func (inst *Store) Get(key string) (interface{}, bool) {
	value, found := inst.Store.Get(key)
	return value, found
}

// Set an item to the store, replacing any existing item. If the duration is 0
// (DefaultExpiration), the store's default expiration time is used.
// If it is -1 (NoExpiration), the item never expires.
func (inst *Store) Set(key string, value interface{}, d time.Duration) {
	inst.Store.Set(key, value, d)
}
