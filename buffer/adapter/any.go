package adapter

import (
	"github.com/NubeDev/flow-eng/buffer"
	"unsafe"
)

type Any struct {
	buff  buffer.TypedReadWriter
	value string
	temp  []byte
}

func NewAny(buff buffer.TypedReadWriter) *Any {
	if buff.Type() != buffer.String {
		panic("unsupported buffer type")
	}
	return &Any{buff, "", make([]byte, buffer.Any)}
}

func (t *Any) Set(value string) {
	const expectedSize = buffer.String

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

func (t *Any) Get() string {
	if _, err := t.buff.Read(t.temp); err != nil {
		panic(err)
	}
	value := *(*string)(unsafe.Pointer(&t.temp[0]))
	t.value = value
	return t.value
}
