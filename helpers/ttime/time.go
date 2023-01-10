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

func Duration(f time.Duration, format string) time.Duration {
	if format == Ms {
		return f * time.Millisecond
	}
	if format == Sec {
		return f * time.Second
	}
	if format == Min {
		return f * time.Minute
	}
	if format == Hr {
		return f * time.Hour
	}
	return f * time.Second
}
