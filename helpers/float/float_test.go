package float

import (
	"fmt"
	"testing"
)

func Test_scaleBetween(t *testing.T) {
	a := Scale(2, 1, 10, 1, 1000)
	fmt.Println(a)
}
