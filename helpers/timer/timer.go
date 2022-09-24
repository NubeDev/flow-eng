package timer

import (
	"github.com/desertbit/timer"
	"time"
)

type TimedDelay interface {
	WaitFor(duration time.Duration) bool
	Stop() (wasActive bool)
	Reset(d time.Duration) bool
}

type Timer struct {
	timer   timer.Timer
	started bool
}

func NewTimer() *Timer {
	return &Timer{*timer.NewStoppedTimer(), false}
}

// Reset changes the timer to expire after duration d.
// It returns true if the timer had been active,
// false if the timer had expired or been stopped.
// The channel t.C is cleared and calling t.Reset() behaves as creating a
// new Timer.
func (t *Timer) Reset(d time.Duration) bool {
	return t.timer.Reset(d)
}

// Stop prevents the Timer from firing.
// It returns true if the call stops the timer,
// false if the timer has already expired or been stopped.
// Stop does not close the channel, to prevent a read from
// the channel succeeding incorrectly.
func (t *Timer) Stop() (wasActive bool) {
	return t.timer.Stop()
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
