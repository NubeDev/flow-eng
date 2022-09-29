package uuid_test

import (
	"github.com/NubeDev/flow-eng/helpers/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUUID(t *testing.T) {
	expected := uuid.Value(1)
	generated := uuid.New()

	assert.Equal(t, expected, generated)
}

func BenchmarkNewUUID(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = uuid.New()
	}
}
