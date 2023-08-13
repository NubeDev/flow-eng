package ttime

import (
	"fmt"
	"github.com/andanhm/go-prettytime"
	"github.com/rvflash/elapsed"
	"regexp"
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

	Sunday    = "Sunday"
	Monday    = "Monday"
	Tuesday   = "Tuesday"
	Wednesday = "Wednesday"
	Thursday  = "Thursday"
	Friday    = "Friday"
	Saturday  = "Saturday"
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

func GetMinDifference(t1 time.Time, t2 time.Time) float64 {
	return t1.Sub(t2).Minutes()
}

func AdjustTime(original time.Time, adjustment string) (time.Time, error) {
	duration, err := parseDuration(adjustment)
	if err != nil {
		return time.Time{}, err
	}

	return original.Add(duration), nil
}

// ParseDuration converts a systemd relative time string into time.Duration
func parseDuration(raw string) (time.Duration, error) {
	re, err := regexp.Compile(`^\s*-?\s*(\d+\s*[a-z]+)`)
	if err != nil {
		return 0, err
	}

	if !re.MatchString(raw) {
		return 0, fmt.Errorf("ParseDuration: incorrect format for raw input %s", raw)
	}

	reNegative, err := regexp.Compile(`^\s*-.*`)
	if err != nil {
		return 0, err
	}
	isNegative := reNegative.MatchString(raw)

	reGroups, err := regexp.Compile(`\d+\s*[a-z]+`)
	if err != nil {
		return 0, err
	}

	matches := reGroups.FindAllString(raw, -1)

	totalDuration := time.Duration(0)
	reSubGroup, err := regexp.Compile(`^(\d+)\s*([a-z]+)$`)
	if err != nil {
		return 0, err
	}
	for _, match := range matches {
		matchTrimmed := strings.Replace(match, " ", "", -1)
		subGroupMatches := reSubGroup.FindStringSubmatch(matchTrimmed)

		// if we run into a case where there aren't exactly two matches
		// then that means this is an unexpected string and we should error out
		if len(subGroupMatches) != 3 {
			return 0, fmt.Errorf("Unexpected match count for '%s': expected 2 and got %d", matchTrimmed, len(subGroupMatches))
		}

		subGroupMatchValue, err := strconv.Atoi(subGroupMatches[1])
		if err != nil {
			return 0, err
		}

		subGroupMatchUnit, err := unitToDuration(subGroupMatches[2])
		if err != nil {
			return 0, err
		}

		totalDuration += time.Duration(subGroupMatchValue) * subGroupMatchUnit
	}

	if isNegative {
		totalDuration *= -1
	}

	return totalDuration, nil
}

// UnitToDuration converts a systemd unit (e.g. "day") to time.Duration
func unitToDuration(unit string) (time.Duration, error) {
	// microseconds
	if matched, err := regexp.MatchString(`^us(ec)?$`, unit); err == nil && matched {
		return time.Microsecond, nil
	}

	// milliseconds
	if matched, err := regexp.MatchString(`^ms(ec)?$`, unit); err == nil && matched {
		return time.Millisecond, nil
	}

	// seconds
	if matched, err := regexp.MatchString(`^s(ec(onds?)?)?$`, unit); err == nil && matched {
		return time.Second, nil
	}

	// minutes
	if matched, err := regexp.MatchString(`^m(in(utes?)?)?$`, unit); err == nil && matched {
		return time.Minute, nil
	}

	// hours
	if matched, err := regexp.MatchString(`^(hr|h(ours?)?)$`, unit); err == nil && matched {
		return time.Hour, nil
	}

	// days
	if matched, err := regexp.MatchString(`^d(ays?)?$`, unit); err == nil && matched {
		return 24 * time.Hour, nil
	}

	// weeks
	if matched, err := regexp.MatchString(`^w(eeks?)?$`, unit); err == nil && matched {
		return 7 * 24 * time.Hour, nil
	}

	// months
	if matched, err := regexp.MatchString(`^(M|months?)$`, unit); err == nil && matched {
		return time.Duration(30.44 * float64(24) * float64(time.Hour)), nil
	}

	// years
	if matched, err := regexp.MatchString(`^y(ears?)?$`, unit); err == nil && matched {
		return time.Duration(365.25 * float64(24) * float64(time.Hour)), nil
	}

	return 0, fmt.Errorf("Unit %s did not match", unit)
}
