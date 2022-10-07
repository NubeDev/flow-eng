package modbuscli

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

func Test_Modbus(t *testing.T) {

}

func Test_Tests(t *testing.T) {

	list, addr17 := InputAddress(1, 1)

	for _, addr := range list {
		fmt.Println(addr)
	}

	pprint.Print(list)
	pprint.Print(addr17)

	listOuts, addr17Out := outputAddress(5, 19)

	pprint.Print(listOuts)
	pprint.Print(addr17Out)

}
