package ttime

import (
	"github.com/rvflash/elapsed"
	"time"
)

const (
	Ms  = "ms"
	Sec = "sec"
	Min = "min"
	Hr  = "hour"
	Day = "day"
)

// TimeSince returns in a human readable format the elapsed time
// eg 12 hours, 12 days
func TimeSince(t time.Time) string {
	return elapsed.Time(t)
}

// Time represents ttime.
type Time interface {
	Now(notInUTC ...bool) time.Time
}

// RealTime is a concrete implementation of Time interface.
type RealTime struct{}

// New initializes and returns a new Time instance.
func New() Time {
	return &RealTime{}
}

// Now returns a timestamp of the current datetime in UTC.
func (rt *RealTime) Now(notUTC ...bool) time.Time {
	if len(notUTC) > 0 {
		return time.Now().UTC()
	}
	return time.Now()
}

func Duration(amount float64, units string) time.Duration {
	if units == Ms {
		return time.Duration(amount * float64(time.Millisecond))
	}
	if units == Sec {
		return time.Duration(amount * float64(time.Second))
	}
	if units == Min {
		return time.Duration(amount * float64(time.Minute))
	}
	if units == Hr {
		return time.Duration(amount * float64(time.Hour))
	}
	return time.Duration(amount * float64(time.Second))
}
