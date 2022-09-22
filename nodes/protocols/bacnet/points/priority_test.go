package points

import (
	"testing"
)

func mergePriority(existing, p2 *PriArray) *PriArray {
	in14, in15 := GetWriteArrayValues(existing)
	if p2 != nil {
		return &PriArray{
			P1:  p2.P1,
			P2:  p2.P2,
			P3:  p2.P2,
			P4:  p2.P4,
			P5:  p2.P5,
			P6:  p2.P6,
			P7:  p2.P7,
			P8:  p2.P8,
			P9:  p2.P9,
			P10: p2.P10,
			P11: p2.P11,
			P12: p2.P12,
			P13: p2.P13,
			P14: in14, // these are reversed for the flow
			P15: in15, // these are reversed for the flow
			P16: p2.P16,
		}
	}
	return nil

}
func TestNewPriArrayAt15(t *testing.T) {

	//first := 11.0
	//second := 12.0
	//third := 13.0
	//p1 := PriArray{P1: &first, P2: &second}
	//p2 := PriArray{P3: &third}
	//po := mergePriorities(p1, p2)
	//fmt.Println("p1>>>", p1)
	//fmt.Println("p2>>>", p2)
	//fmt.Println("po>>>", po)
	//fmt.Println("P1>>>", *po.P1)
	//fmt.Println("P2>>>", *po.P2)
	//fmt.Println("P3>>>", *po.P3)
	//fmt.Println("P4>>>", po.P4)
	//
	//fmt.Println("!!!!!!!!!!!!!!!!!")
	//p1.P1 = nil
	//p1.P2 = nil
	//po = mergePriorities(p1, p2)
	//pprint.Print(po)

}
