package node

const (
	nodeA = "nodeA"
	nodeB = "nodeB"
)

func GetNodeSpec(name string, body *Node) *Node {
	switch name {
	case nodeA:
		return SpecNodeA(body)
	case nodeB:
		return SpecNodeA(body)
	}
	return nil
}

//type inputSpec struct {
//	portName string
//	portType string
//}
//
//type Spec struct {
//	Name        string
//	InputCount  int
//	OutputCount int
//	Category    string
//	Description string
//	Version     string
//}
//
//var Schema = struct {
//	NodeA *Spec
//	NodeB *Spec
//}{
//	NodeA: &Spec{
//		Name:        nodeA,
//		InputCount:  2,
//		OutputCount: 1,
//		Category:    "",
//		Description: "",
//		Version:     "",
//	},
//	NodeB: &Spec{
//		Name:        nodeB,
//		InputCount:  2,
//		OutputCount: 1,
//		Category:    "",
//		Description: "",
//		Version:     "",
//	},
//}
