package flowctrl

import (
	"github.com/desertbit/timer"
	"time"
)

type TimedDelay interface {
	WaitFor(duration time.Duration) bool
}

type Timer struct {
	timer   timer.Timer
	started bool
}

func NewTimer() *Timer {
	return &Timer{*timer.NewStoppedTimer(), false}
}

func (t *Timer) WaitFor(duration time.Duration) bool {
	if !t.started {
		t.timer.Reset(duration)
		t.started = true
	}
	select {
	case <-t.timer.C:
		t.started = false
		return true
	default:
		return false
	}
}
