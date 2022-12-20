package bacnetio

import (
	"fmt"
	"testing"
)

func Test_bacnetAddress(t *testing.T) {
	bacnetAddress(4, "AO", "UO")
	a := bacnetAddress(4, "AO", "UO")
	for _, s := range a {
		fmt.Println(s)
	}
}
