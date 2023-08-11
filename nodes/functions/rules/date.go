package rules

import (
	"time"
)

// Sleep will delay the program for the `duration` passed in (duration is units seconds)
func (inst *RQL) Sleep(duration int) {
	d := time.Duration(duration)
	time.Sleep(d * time.Second)
}

// TimeUTC returns time in UTC
func (inst *RQL) TimeUTC() time.Time {
	return time.Now().UTC()
}

// TimeDate returns time/date formatted as `2006.01.02 15:04:05`
func (inst *RQL) TimeDate() string {
	return time.Now().Format("2006.01.02 15:04:05")
}

// TimeWithMS returns time formatted as `15:04:05.000`
func (inst *RQL) TimeWithMS() string {
	return time.Now().Format("15:04:05.000")
}

// Time returns time formatted as `15:04:05`
func (inst *RQL) Time() string {
	return time.Now().Format("15:04:05")
}

// Date returns date formatted as `2006.01.02`
func (inst *RQL) Date() string {
	return time.Now().Format("2006.01.02")
}

// Year returns current year
func (inst *RQL) Year() string {
	return time.Now().Format("2006")
}

// Day returns current day
func (inst *RQL) Day() string {
	return time.Now().Format("Monday")
}

// TimeDateDay returns time/date formatted as `2006-01-02 15:04:05 Monday`
func (inst *RQL) TimeDateDay() string {
	return time.Now().Format("2006-01-02 15:04:05 Monday")
}
