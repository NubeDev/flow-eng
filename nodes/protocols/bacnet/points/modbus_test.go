package points

import (
	"github.com/NubeDev/flow-eng/helpers/pprint"
	"testing"
)

func Test_outputAddress(t *testing.T) {
	r, _ := outputAddress(4, 0)
	pprint.PrintJSON(r)

}

func TestModbusBuildOutput(t *testing.T) {
	got, got1 := ModbusBuildOutput(IoTypeVolts, 7)
	pprint.Print(got)
	pprint.Print(got1)

}
