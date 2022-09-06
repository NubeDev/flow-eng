package parameter

import (
	"bytes"
	"github.com/francoispqt/gojay"
)

type Store interface {
	Store(container *Container) error
}

type Restore interface {
	Restore(container *Container) error
}

type StoreRestorer interface {
	Store
	Restore
}

type BaseStore struct {
	buffer  *bytes.Buffer
	encoder *gojay.Encoder
	decoder *gojay.Decoder
}

func NewBaseStore() *BaseStore {
	buff := bytes.Buffer{}
	return &BaseStore{&buff, gojay.NewEncoder(&buff), gojay.NewDecoder(&buff)}
}

func (store *BaseStore) Bytes() []byte {
	return store.buffer.Bytes()
}

func (store *BaseStore) Write(data []byte) {
	store.buffer.Reset()
	store.buffer.Write(data)
}

func (store *BaseStore) Store(container *Container) error {
	store.buffer.Reset()
	err := store.encoder.EncodeArray(container)
	if err != nil {
		return err
	}
	return nil
}

func (store *BaseStore) Restore(container *Container) error {
	var loadedContainer Container
	err := store.decoder.DecodeArray(&loadedContainer)
	if err != nil {
		return err
	}
	temp := make([]byte, 8)
	for i := 0; i < loadedContainer.Size(); i++ {
		param := loadedContainer.Parameters[i]
		paramName := param.GetName()
		for j := 0; j < container.Size(); j++ {
			if container.Parameters[j].GetName() != paramName {
				continue
			}
			read, err := param.Read(temp[:param.Type()])
			if err != nil {
				return err
			}
			_, err = container.Parameters[j].Write(temp[:read])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
