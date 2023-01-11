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

	palle2t, err := EncodePallet()
	if err != nil {
		return
	}

	pprint.PrintJSON(palle2t)
}
