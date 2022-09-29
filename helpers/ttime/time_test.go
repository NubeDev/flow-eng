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
