package nodes

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

func Test_test(t *testing.T) {
	nodeList, err := EncodePallet()
	fmt.Println(err)
	pprint.PrintJOSN(nodeList)
}
