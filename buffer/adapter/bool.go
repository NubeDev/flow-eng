package adapter

import (
	"github.com/NubeDev/flow-eng/buffer"
	"unsafe"
)

type Bool struct {
	buff  buffer.TypedReadWriter
	value bool
	temp  []byte
}

func NewBool(buff buffer.TypedReadWriter) *Bool {
	if buff.Type() != buffer.Bool {
		panic("unsupported buffer type")
	}
	return &Bool{buff, false, make([]byte, buffer.Bool)}
}

func (t *Bool) Set(value bool) {
	const expectedSize = buffer.Bool

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

func (t *Bool) Get() bool {
	if _, err := t.buff.Read(t.temp); err != nil {
		panic(err)
	}
	value := *(*bool)(unsafe.Pointer(&t.temp[0]))
	t.value = value
	return t.value
}
