package parameter

import (
	"io/ioutil"
	"os"
	"sync"
)

type FileStore struct {
	base  *BaseStore
	path  string
	mode  os.FileMode
	mutex *sync.Mutex
}

func NewFileStore(path string, mode os.FileMode) *FileStore {
	return &FileStore{NewBaseStore(), path, mode, &sync.Mutex{}}
}

func (store *FileStore) Store(container *Container) error {
	if container == nil {
		return nil
	}

	store.mutex.Lock()
	err := store.base.Store(container)
	store.mutex.Unlock()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.path, store.base.buffer.Bytes(), 0655)
}

func (store *FileStore) Restore(container *Container) error {
	if container == nil {
		return nil
	}
	data, err := ioutil.ReadFile(store.path)
	if err != nil {
		return err
	}
	store.mutex.Lock()
	store.base.Write(data)
	store.mutex.Unlock()

	return store.base.Restore(container)
}
