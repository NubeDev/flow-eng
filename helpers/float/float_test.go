package float

import (
	"fmt"
	"testing"
)

func Test_scaleBetween(t *testing.T) {
	a := Scale(100, 0, 100, 0, 10)
	fmt.Println(a)
}

func TestRoundTo(t *testing.T) {
	o := RoundTo(22.000, 1)
	fmt.Println(o)
}
