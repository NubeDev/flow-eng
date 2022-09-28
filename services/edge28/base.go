package edge28lib

import "github.com/NubeDev/flow-eng/services/clients/edgerest"

type Edge28 struct {
	client *edgerest.Client
}

func New(ip string, port ...int) *Edge28 {
	p := 0
	if len(port) > 0 {
		p = port[0]
	}
	client := edgerest.New(ip, p)
	return &Edge28{
		client: client,
	}
}
