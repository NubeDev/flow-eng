package bac

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

func Test_CleanArray(t *testing.T) {
	pri := CleanArray("{Null,33.3,Null,Null,Null,Null,Null,Null,Null,Null,Null,Null,Null,Null,Null,11.000000}")
	pprint.PrintJOSN(pri)

	fmt.Println(GetHighest(pri))
}
