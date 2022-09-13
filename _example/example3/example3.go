package main

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/nodes"
)

func main() {
	schema := nodes.All()

	//for _, n := range schema {
	//
	//	//aa, ok := n.(node.BaseNode)
	//
	//	fmt.Println(aa, ok)
	//}

	pprint.PrintJOSN(schema)
}
