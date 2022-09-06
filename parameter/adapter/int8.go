package adapter

import (
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/parameter"
	"github.com/francoispqt/gojay"
	"sync"
)

type Int8 struct {
	parameter.NamedReadWriter
	adapter *adapter.Int8
	mutex   *sync.RWMutex
}

func NewInt8(value int8, parameter parameter.NamedReadWriter) *Int8 {
	typed := adapter.NewInt8(parameter)
	typed.Set(value)
	return &Int8{parameter, typed, &sync.RWMutex{}}
}

func (p *Int8) MarshalJSONObject(enc *gojay.Encoder) {
	p.NamedReadWriter.MarshalJSONObject(enc)
	enc.Int8Key("value", p.Get())
}

func (p *Int8) Set(value int8) {
	p.mutex.Lock()
	p.adapter.Set(value)
	p.mutex.Unlock()
}

func (p *Int8) Get() int8 {
	p.mutex.RLock()
	value := p.adapter.Get()
	p.mutex.RUnlock()
	return value
}
