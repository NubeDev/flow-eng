package nodes

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

func TestEncodePallet(t *testing.T) {
	p := All()
	name := "add"
	for _, spec := range p {
		if spec.Info.Name == name {
			pprint.Print(spec)
		}
	}

	palle2t, err := EncodePalle2t()
	if err != nil {
		return
	}

	pprint.PrintJOSN(palle2t)
}
