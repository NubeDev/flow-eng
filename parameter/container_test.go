package parameter_test

import (
	"bytes"
	"github.com/NubeDev/flow-eng/buffer"
	"github.com/NubeDev/flow-eng/parameter"
	"github.com/NubeDev/flow-eng/parameter/adapter"
	"github.com/francoispqt/gojay"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainer_MarshalJSON(t *testing.T) {
	//GIVEN
	expectedJSON := `[{"name":"paramA","description":"description","unit":"","readOnly":false,"userVisible":true,"raw":[123],"value":123},{"name":"paramB","description":"description","unit":"%","readOnly":false,"userVisible":true,"raw":[0],"value":0},{"name":"paramC","description":"description","unit":"Hz","readOnly":false,"userVisible":true,"raw":[0],"value":0}]`
	paramA := adapter.NewInt8(123, parameter.New(buffer.Int8, "paramA", "description", parameter.UnitNone, false, true))
	paramB := adapter.NewInt8(0, parameter.New(buffer.Int8, "paramB", "description", parameter.UnitPercent, false, true))
	paramC := adapter.NewInt8(0, parameter.New(buffer.Int8, "paramC", "description", parameter.UnitHertz, false, true))
	container := parameter.NewContainer(paramA, paramB, paramC)
	//WHEN
	buff := bytes.Buffer{}
	encoder := gojay.NewEncoder(&buff)
	err := encoder.EncodeArray(container)
	//THEN
	assert.Nil(t, err)
	assert.Equal(t, expectedJSON, string(buff.Bytes()))
}

func TestContainer_UnmarshalJSON(t *testing.T) {
	//GIVEN
	expectedParamCount := 3
	sourceJSON := `[{"name":"paramA","description":"description","unit":"","readOnly":false,"userVisible":true,"raw":[123],"value":123},{"name":"paramB","description":"description","unit":"%","readOnly":false,"userVisible":true,"raw":[0],"value":0},{"name":"paramC","description":"description","unit":"Hz","readOnly":false,"userVisible":true,"raw":[0],"value":0}]`
	buff := bytes.NewBuffer([]byte(sourceJSON))
	decoder := gojay.NewDecoder(buff)
	container := parameter.NewContainer()
	//WHEN
	err := decoder.DecodeArray(container)
	data := []byte{0}
	container.Parameters[0].Read(data)
	//THEN
	assert.Nil(t, err)
	assert.Equal(t, expectedParamCount, len(container.Parameters))
	assert.Equal(t, []byte{123}, data)
}

func BenchmarkContainer_MarshalJSON(b *testing.B) {
	b.ReportAllocs()

	paramA := adapter.NewInt8(0, parameter.New(buffer.Int8, "paramA", "description", parameter.UnitNone, false, true))
	paramB := adapter.NewInt8(0, parameter.New(buffer.Int8, "paramB", "description", parameter.UnitPercent, false, true))
	paramC := adapter.NewInt8(0, parameter.New(buffer.Int8, "paramC", "description", parameter.UnitHertz, false, true))
	paramD := adapter.NewInt8(0, parameter.New(buffer.Int8, "paramD", "description", parameter.UnitHertz, false, true))
	paramE := adapter.NewInt8(0, parameter.New(buffer.Int8, "paramE", "description", parameter.UnitHertz, false, true))
	paramF := adapter.NewInt8(0, parameter.New(buffer.Int8, "paramF", "description", parameter.UnitHertz, false, true))
	paramG := adapter.NewInt8(0, parameter.New(buffer.Int8, "paramG", "description", parameter.UnitHertz, false, true))
	paramH := adapter.NewInt8(0, parameter.New(buffer.Int8, "paramH", "description", parameter.UnitHertz, false, true))
	paramI := adapter.NewInt8(0, parameter.New(buffer.Int8, "paramI", "description", parameter.UnitHertz, false, true))
	paramJ := adapter.NewInt8(0, parameter.New(buffer.Int8, "paramJ", "description", parameter.UnitHertz, false, true))
	container := parameter.NewContainer(paramA, paramB, paramC, paramD, paramE, paramF, paramG, paramH, paramI, paramJ)

	buff := &bytes.Buffer{}
	buff.Grow(200)
	enc := gojay.NewEncoder(buff)

	for i := 0; i < b.N; i++ {
		buff.Reset()
		err := enc.EncodeArray(container)
		if err != nil {
			b.Fail()
		}
	}
}
