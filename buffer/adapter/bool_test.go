package adapter

import (
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBool_Set(t *testing.T) {
	//GIVEN
	b := buffer.NewConst(buffer.Bool)
	adapter := NewBool(b)
	value := true
	valueAsBytes := []byte{0x01}
	//WHEN
	adapter.Set(value)
	rawData := []byte{0}
	b.Read(rawData)
	//THEN
	assert.Equal(t, value, adapter.Get())
	assert.Equal(t, valueAsBytes, rawData)
}
