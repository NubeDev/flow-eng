package nodes

import (
	"fmt"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

func TestGetSchema(t *testing.T) {
	n, err := GetSchema("http-get")
	fmt.Println(err)
	pprint.PrintJOSN(n)
}
