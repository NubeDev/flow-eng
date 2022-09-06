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

func TestFloat32_Write(t *testing.T) {
	// GIVEN
	param := adapter.NewFloat32(0, parameter.New(buffer.Float32, "param", "description", parameter.UnitNone, false, true))
	expected := float32(123.1234)
	//WHEN
	param.Set(expected)
	read := param.Get()
	//THEN
	assert.Equal(t, expected, read)
}

func TestFloat32_MarshalJSON(t *testing.T) {
	//GIVEN
	expectedJSON := `{"name":"param","description":"description","unit":"%","readOnly":false,"userVisible":true,"raw":[46,63,246,66],"value":123.1234}`
	param := adapter.NewFloat32(123.1234, parameter.New(buffer.Float32, "param", "description", parameter.UnitPercent, false, true))
	//WHEN
	buffer := bytes.Buffer{}
	encoder := gojay.NewEncoder(&buffer)
	err := encoder.EncodeObject(param)
	//THEN
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, string(buffer.Bytes()))
}

func BenchmarkFloat32_MarshalJSON(b *testing.B) {
	b.ReportAllocs()

	param := adapter.NewFloat32(123, parameter.New(buffer.Float32, "param", "description", parameter.UnitPercent, false, true))
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

func BenchmarkFloat32_Write(b *testing.B) {
	b.ReportAllocs()

	param := adapter.NewFloat32(0, parameter.New(buffer.Float32, "param", "description", parameter.UnitNone, false, true))
	expected := float32(123.1234)

	for i := 0; i < b.N; i++ {
		param.Set(expected)
	}
}

func BenchmarkFloat32_ParallelWrite(b *testing.B) {
	b.ReportAllocs()

	param := adapter.NewFloat32(0, parameter.New(buffer.Float32, "param", "description", parameter.UnitNone, false, true))
	expected := float32(123.1234)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			param.Set(expected)
		}
	})
}
