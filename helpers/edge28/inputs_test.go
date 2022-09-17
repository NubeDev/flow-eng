package edge28

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

func TestGetUI(t *testing.T) {
	ui, err := testUI()
	if err != nil {
		fmt.Println(err)
		return
	}
	is, err := ProcessInput(ui.Val.UI1.Val, "thermistor_10k_type_2")
	fmt.Println(is, err)

	pprint.PrintJOSN(ui)
}

func TestGetDI(t *testing.T) {

	ui, err := testDI()

	if err != nil {
		fmt.Println(err)
		return
	}
	is, err := ProcessInput(ui.Val.DI1.Val, "digital")
	fmt.Println(is, err)

}
