package main

import (
	"github.com/NubeDev/flow-eng/_example/nodes"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
)

func main() {
	schema := nodes.All()
	pprint.PrintJOSN(schema)
}
