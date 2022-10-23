package flow

const (
	category    = "flow"
	flowNetwork = "flow-network"
	flowPoint   = "flow-point"
)

type point struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type pointStore struct {
	parentID string
	points   []*point
	payloads []*pointDetails
}

type pointDetails struct {
	nodeUUID    string
	topic       string
	payload     string
	isWriteable bool
}

func addUpdatePayload(nodeUUID string, p *pointStore, newPayload *pointDetails) (data *pointStore, found bool) {
	for i, payload := range p.payloads {
		if payload.nodeUUID == nodeUUID {
			p.payloads[i] = newPayload
			found = true
		}
	}
	if !found {
		p.payloads = append(p.payloads, newPayload)
	}
	return p, found
}
