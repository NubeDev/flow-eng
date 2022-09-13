package main

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/nodes"
)

func main() {
	schema := nodes.All()
	pprint.PrintJOSN(schema)
}
