package flowctrl

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimer_WaitFor(t *testing.T) {
	//GIVEN
	timer := NewTimer()
	start := time.Now()
	duration := 1 * time.Second
	//WHEN
	for {
		result := timer.WaitFor(duration)
		if result == true {
			break
		}
	}
	end := time.Since(start)
	//THEN
	assert.Equal(t, int(duration.Seconds()), int(end.Seconds()))
}
