package buffer

import (
	"sync"
)

type Sync struct {
	Const
	mutex sync.RWMutex
}

func NewSync(_type Type) *Sync {
	return &Sync{*NewConst(_type), sync.RWMutex{}}
}

func (b *Sync) Write(data []byte) (int, error) {
	b.mutex.Lock()
	written, err := b.Const.Write(data)
	b.mutex.Unlock()
	return written, err
}

func (b *Sync) Read(data []byte) (int, error) {
	b.mutex.RLock()
	read, err := b.Const.Read(data)
	b.mutex.RUnlock()
	return read, err
}

func (b *Sync) Copy(other *Sync) (int, error) {
	if b._type != other._type {
		return 0, ErrTypesMismatch
	}

	b.mutex.RLock()
	other.mutex.Lock()
	copied := fastcopy(other.data, b.data)
	b.mutex.RUnlock()
	other.mutex.Unlock()
	return copied, nil
}
