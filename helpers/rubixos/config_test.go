package rubixos

import (
	"fmt"

	"testing"
)

func TestGetConfig(t *testing.T) {
	got, err := GetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(got.UnitsImperial)
}

func TestGetConfigBACnet(t *testing.T) {
	got, err := GetBACnetConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(got)
}
