package adapter

import (
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInt8_Set(t *testing.T) {
	//GIVEN
	b := buffer.NewConst(buffer.Int8)
	adapter := NewInt8(b)
	value := int8(123)
	valueAsBytes := []byte{0x7b}
	//WHEN
	adapter.Set(value)
	rawData := []byte{0}
	b.Read(rawData)
	//THEN
	assert.Equal(t, value, adapter.Get())
	assert.Equal(t, valueAsBytes, rawData)
}

func BenchmarkInt8_Write(b *testing.B) {
	buff := buffer.NewConst(buffer.Int8)
	intValue := NewInt8(buff)
	toWrite := int8(123)

	for i := 0; i < b.N; i++ {
		intValue.Set(toWrite)
	}
}

func BenchmarkInt8_Read(b *testing.B) {
	buff := buffer.NewConst(buffer.Int8)
	intValue := NewInt8(buff)
	toWrite := int8(123)
	intValue.Set(toWrite)

	for i := 0; i < b.N; i++ {
		_ = intValue.Get()
	}
}
