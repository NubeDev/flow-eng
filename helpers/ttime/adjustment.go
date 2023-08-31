package ttime

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetBeginOfYear(t time.Time) time.Time {
	y, _, _ := t.Date()
	return time.Date(y, 1, 1, 0, 0, 0, 0, t.Location())
}

func GetBeginOfMonth(t time.Time) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, t.Location())
}

func GetBeginOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func GetBeginOfHour(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, t.Hour(), 0, 0, 0, t.Location())
}

func GetLastDayOfYear(t time.Time) time.Time {
	y, _, _ := t.Date()
	t = time.Date(y+1, 1, 1, 0, 0, 0, 0, t.Location())
	return t.AddDate(0, 0, -1)
}

func GetLastDayOfMonth(t time.Time) time.Time {
	return GetBeginOfMonth(t).AddDate(0, 1, -1)
}

func GetTomorrow() time.Time {
	return GetBeginOfDay(time.Now()).AddDate(0, 0, 1)
}

func GetYesterday() time.Time {
	return GetBeginOfDay(time.Now()).AddDate(0, 0, -1)
}

func GeDayAfterTomorrow() time.Time {
	return GetTomorrow().AddDate(0, 0, 1)
}

func GeDayBeforeYesterday() time.Time {
	return GetYesterday().AddDate(0, 0, -1)
}

func AddSeconds(t time.Time, second int) time.Time {
	return t.Add(time.Duration(second) * time.Second)
}

func AddMinutes(t time.Time, minutes int) time.Time {
	return t.Add(time.Duration(minutes) * time.Minute)
}

func AddHour(t time.Time, hours int) time.Time {
	return t.Add(time.Duration(hours) * time.Hour)
}

func NextDay(t time.Time) time.Time {
	return t.AddDate(0, 0, 1)
}

func ParseExcelNumber(days int) time.Time {
	start := time.Date(1899, 12, 30, 0, 0, 0, 0, time.Now().Location())
	return start.Add(time.Duration(days*24) * time.Hour)
}

func GetYearsDifference(t1 time.Time, t2 time.Time) float64 {
	return t1.Sub(t2).Hours() / 24 / 365
}

func GetHoursDifference(t1 time.Time, t2 time.Time) float64 {
	return t1.Sub(t2).Hours()
}

func GetMinDifference(t1 time.Time, t2 time.Time) float64 {
	return t1.Sub(t2).Minutes()
}

func GetSecDifference(t1 time.Time, t2 time.Time) float64 {
	return t1.Sub(t2).Seconds()
}

func GetDaysBetween(t1 time.Time, t2 time.Time) int {
	duration := int64(0)
	timestamp1 := GetBeginOfDay(t1).Unix()
	timestamp2 := GetBeginOfDay(t2).Unix()

	if timestamp1 > timestamp2 {
		duration = timestamp1 - timestamp2
	} else {
		duration = timestamp2 - timestamp1
	}

	return int((duration / 86400) | 0)
}

/*
AdjustTime takes a systemd time adjustment string and uses it to modify a time.Time
	- adjustment examples
	- 1day
	- 1 day
	- 1day
	- 15 min
	- 1 hour
	- 2 hour
	- 2years
*/
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
