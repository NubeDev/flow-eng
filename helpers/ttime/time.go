package ttime

import (
	"errors"
	"github.com/andanhm/go-prettytime"
	"github.com/rvflash/elapsed"
	"strconv"
	"strings"
	"time"
)

const (
	Ms  = "ms"
	Sec = "sec"
	Min = "min"
	Hr  = "hour"
	Day = "day"

	Sun  = "Sun"
	Mon  = "Mon"
	Tue  = "Tue"
	Wed  = "Wed"
	Thur = "Thur"
	Fri  = "Fri"
	Sat  = "Sat"
)

func TimePretty(t time.Time) string {
	return prettytime.Format(t)
}

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

func ParseTime(str string) (hour, min, sec int, err error) {
	chunks := strings.Split(str, ":")
	var hourStr, minStr, secStr string
	switch len(chunks) {
	case 1:
		hourStr = chunks[0]
		minStr = "0"
		secStr = "0"
	case 2:
		hourStr = chunks[0]
		minStr = chunks[1]
		secStr = "0"
	case 3:
		hourStr = chunks[0]
		minStr = chunks[1]
		secStr = chunks[2]
	}
	hour, err = strconv.Atoi(hourStr)
	if err != nil {
		return 0, 0, 0, errors.New("bad time")
	}
	min, err = strconv.Atoi(minStr)
	if err != nil {
		return 0, 0, 0, errors.New("bad time")
	}
	sec, err = strconv.Atoi(secStr)
	if err != nil {
		return 0, 0, 0, errors.New("bad time")
	}

	if hour > 23 || min > 59 || sec > 59 {
		return 0, 0, 0, errors.New("bad time")
	}

	return
}
