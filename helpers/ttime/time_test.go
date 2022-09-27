package ttime

import (
	"fmt"
	"testing"
)

func TestRealTime_Now(t *testing.T) {
	rt := &RealTime{}
	fmt.Println(rt.Now())
	fmt.Println(rt.Now(true))

}
