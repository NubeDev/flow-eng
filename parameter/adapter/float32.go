package adapter

import (
	"github.com/NubeDev/flow-eng/buffer/adapter"
	"github.com/NubeDev/flow-eng/parameter"
	"github.com/francoispqt/gojay"
	"sync"
)

type Float32 struct {
	parameter.NamedReadWriter
	adapter *adapter.Float32
	mutex   *sync.RWMutex
}

func NewFloat32(value float32, parameter parameter.NamedReadWriter) *Float32 {
	typed := adapter.NewFloat32(parameter)
	typed.Set(value)
	return &Float32{parameter, typed, &sync.RWMutex{}}
}

func (p *Float32) MarshalJSONObject(enc *gojay.Encoder) {
	p.NamedReadWriter.MarshalJSONObject(enc)
	enc.Float32Key("value", p.Get())
}

func (p *Float32) Set(value float32) {
	p.mutex.Lock()
	p.adapter.Set(value)
	p.mutex.Unlock()
}

func (p *Float32) Get() float32 {
	p.mutex.RLock()
	value := p.adapter.Get()
	p.mutex.RUnlock()
	return value
}
