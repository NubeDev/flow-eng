package node

type NodeInfo struct {
	Name        string `json:"-"`
	Type        string `json:"-"`
	Description string `json:"-"`
	Version     string `json:"-"`
}

type Node interface {
	Process()
	Cleanup()
	Info() NodeInfo
}
