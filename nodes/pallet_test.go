package nodes

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

func TestEncodePallet(t *testing.T) {
	p, _ := EncodePallet()

	pprint.PrintJOSN(p)
}
