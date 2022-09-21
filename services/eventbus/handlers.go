package eventbus

import (
	"github.com/NubeIO/nubeio-rubix-lib-models-go/pkg/v1/model"
	"github.com/mustafaturan/bus/v3"
	"strings"
)

// GetTopicPart will split the topics
func GetTopicPart(topic string, index int, contains string) string {
	s := strings.Split(topic, ".")
	for i, e := range s {
		if i == index {
			if strings.Contains(e, contains) { // if topic has pnt (is uuid of point)
				return e
			}
		}
	}
	return ""
}

// IsNetwork check if the payload is of type device
func IsNetwork(topic string, payload bus.Event) (*model.Network, error) {
	if GetTopicPart(topic, 3, "net") != "" {
		p, _ := payload.Data.(*model.Network)
		return p, nil
	}
	return nil, nil
}
