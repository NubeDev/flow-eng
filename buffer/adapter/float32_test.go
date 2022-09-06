package adapter

import (
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFloat32_Set(t *testing.T) {
	//GIVEN
	b := buffer.NewConst(buffer.Float32)
	adapter := NewFloat32(b)
	value := float32(1234.1234)
	valueAsBytes := []byte{0xf3, 0x43, 0x9a, 0x44}
	//WHEN
	adapter.Set(value)
	rawData := []byte{0, 0, 0, 0}
	b.Read(rawData)
	//THEN
	assert.Equal(t, value, adapter.Get())
	assert.Equal(t, valueAsBytes, rawData)
}
