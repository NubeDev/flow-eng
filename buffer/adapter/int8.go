package adapter

import (
	"github.com/NubeDev/flow-eng/buffer"
	"unsafe"
)

type Int8 struct {
	buff  buffer.TypedReadWriter
	value int8
	temp  []byte
}

func NewInt8(buff buffer.TypedReadWriter) *Int8 {
	if buff.Type() != buffer.Int8 {
		panic("unsupported buffer type")
	}
	return &Int8{buff, 0, make([]byte, buffer.Int8)}
}

func (t *Int8) Set(value int8) {
	const expectedSize = buffer.Int8

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

func (t *Int8) Get() int8 {
	if _, err := t.buff.Read(t.temp); err != nil {
		panic(err)
	}
	value := *(*int8)(unsafe.Pointer(&t.temp[0]))
	t.value = value
	return t.value
}
