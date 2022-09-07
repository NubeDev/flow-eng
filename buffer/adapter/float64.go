package adapter

import (
	"github.com/NubeDev/flow-eng/buffer"
	"unsafe"
)

type Float64 struct {
	buff  buffer.TypedReadWriter
	value *float64
	temp  []byte
}

func NewFloat64(buff buffer.TypedReadWriter) *Float64 {
	if buff.Type() != buffer.Float64 {
		panic("unsupported buffer type")
	}
	return &Float64{buff, nil, make([]byte, buffer.Float64)}
}

func (t *Float64) Set(value *float64) {
	const expectedSize = buffer.Float64

	t.value = value
	bytes := (*[expectedSize]byte)(unsafe.Pointer(&t.value))[:]
	written, err := t.buff.Write(bytes)
	if err != nil {
		panic(err)
	}
	if written != int(expectedSize) {
		panic("data write failed")
	}
}

func (t *Float64) Get() *float64 {
	if _, err := t.buff.Read(t.temp); err != nil {
		panic(err)
	}
	value := *(**float64)(unsafe.Pointer(&t.temp[0]))
	t.value = value
	return t.value
}
