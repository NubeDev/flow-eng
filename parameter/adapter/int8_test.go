package adapter_test

import (
	"bytes"
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/parameter"
	"github.com/NubeDev/flow-eng/parameter/adapter"
	"github.com/francoispqt/gojay"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInt8_Write(t *testing.T) {
	// GIVEN
	param := adapter.NewInt8(0, parameter.New(buffer.Int8, "param", "description", parameter.UnitNone, false, true))
	expected := int8(123)
	//WHEN
	param.Set(expected)
	read := param.Get()
	//THEN
	assert.Equal(t, expected, read)
}

func TestInt8_MarshalJSON(t *testing.T) {
	//GIVEN
	expectedJSON := `{"name":"param","description":"description","unit":"%","readOnly":false,"userVisible":true,"raw":[123],"value":123}`
	param := adapter.NewInt8(123, parameter.New(buffer.Int8, "param", "description", parameter.UnitPercent, false, true))
	//WHEN
	buffer := bytes.Buffer{}
	encoder := gojay.NewEncoder(&buffer)
	err := encoder.EncodeObject(param)
	//THEN
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, string(buffer.Bytes()))
}

func BenchmarkInt8_MarshalJSON(b *testing.B) {
	b.ReportAllocs()

	param := adapter.NewInt8(123, parameter.New(buffer.Int8, "param", "description", parameter.UnitPercent, false, true))
	buffer := &bytes.Buffer{}
	buffer.Grow(100)
	enc := gojay.NewEncoder(buffer)

	for i := 0; i < b.N; i++ {
		buffer.Reset()
		err := enc.EncodeObject(param)
		if err != nil {
			b.Fail()
		}
	}
}

func BenchmarkInt8_Write(b *testing.B) {
	b.ReportAllocs()

	param := adapter.NewInt8(0, parameter.New(buffer.Int8, "param", "description", parameter.UnitNone, false, true))
	expected := int8(123)

	for i := 0; i < b.N; i++ {
		param.Set(expected)
	}
}

func BenchmarkInt8_ParallelWrite(b *testing.B) {
	b.ReportAllocs()

	param := adapter.NewInt8(0, parameter.New(buffer.Int8, "param", "description", parameter.UnitNone, false, true))
	expected := int8(123)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			param.Set(expected)
		}
	})
}
