package modbusclient

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

func Test_Modbus(t *testing.T) {

}

func Test_Tests(t *testing.T) {

	list, addr17 := InputAddress(5, 21)

	pprint.Print(list)
	pprint.Print(addr17)

	listOuts, addr17Out := outputAddress(5, 19)

	pprint.Print(listOuts)
	pprint.Print(addr17Out)

}
