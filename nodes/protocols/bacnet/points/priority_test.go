package points

import (
	"github.com/NubeDev/flow-eng/helpers/float"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"testing"
)

func TestNewPriArrayAt15(t *testing.T) {
	a := NewPriArrayAt15(12)
	aa := [16]*float64{}
	aa[15] = float.New(33)
	pprint.Print(aa)
	pprint.Print(a)
}
