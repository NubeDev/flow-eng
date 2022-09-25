package helpers

import (
	"fmt"
	"testing"
)

func TestShortUUID(t *testing.T) {
	fmt.Println(ShortUUID())
	fmt.Println(UUID("con"))
}
