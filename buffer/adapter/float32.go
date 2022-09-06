package adapter

import (
	"github.com/NubeDev/flow-eng/buffer"
	"unsafe"
)

type Float32 struct {
	buff  buffer.TypedReadWriter
	value float32
	temp  []byte
}

func NewFloat32(buff buffer.TypedReadWriter) *Float32 {
	if buff.Type() != buffer.Float32 {
		panic("unsupported buffer type")
	}
	return &Float32{buff, 0, make([]byte, buffer.Float32)}
}

func (t *Float32) Set(value float32) {
	const expectedSize = buffer.Float32

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

func (t *Float32) Get() float32 {
	if _, err := t.buff.Read(t.temp); err != nil {
		panic(err)
	}
	value := *(*float32)(unsafe.Pointer(&t.temp[0]))
	t.value = value
	return t.value
}
