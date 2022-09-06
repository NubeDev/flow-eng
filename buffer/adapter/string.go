package adapter

import (
	"github.com/NubeDev/flow-eng/buffer"
	"unsafe"
)

type String struct {
	buff  buffer.TypedReadWriter
	value string
	temp  []byte
}

func NewString(buff buffer.TypedReadWriter) *String {
	if buff.Type() != buffer.String {
		panic("unsupported buffer type")
	}
	return &String{buff, "", make([]byte, buffer.String)}
}

func (t *String) Set(value string) {
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

func (t *String) Get() string {
	if _, err := t.buff.Read(t.temp); err != nil {
		panic(err)
	}
	value := *(*string)(unsafe.Pointer(&t.temp[0]))
	t.value = value
	return t.value
}
