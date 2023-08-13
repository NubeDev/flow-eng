package ttime

import (
	"fmt"
	"testing"
	"time"
)

func TestRealTime_Since(t *testing.T) {
	time_ := time.Now().Add(-1234 * time.Hour)
	fmt.Println(TimeSince(time_))

	time_ = time.Now().Add(-1 * time.Minute)
	fmt.Println(TimeSince(time_))

	time_ = time.Now().Add(-1 * time.Second)
	fmt.Println(TimeSince(time_))

}

func TestRealTime_Now(t *testing.T) {
	rt := &RealTime{}
	fmt.Println(rt.Now())
	fmt.Println(rt.Now(true))

}

func TestAdjustTime(t *testing.T) {

	now := time.Now()

	fmt.Println("now:", now.Format(time.Stamp))

	count := 10
	then := now.Add(time.Duration(count) * time.Minute)
	// if we had fix number of units to subtract, we can use following line instead fo above 2 lines. It does type convertion automatically.
	// then := now.Add(-10 * time.Minute)
	fmt.Println("10 minutes ago:", then.Format(time.Stamp))

}
