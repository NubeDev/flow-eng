package node

import "reflect"

func Ports(node *Node, direction Direction) []Port {
	ports := make([]Port, 0, 1)
	for _, input := range node.GetInputs() {
		if input.Direction() != direction {
			continue
		}
		ports = append(ports, input)
	}
	for _, output := range node.GetOutputs() {
		if output.Direction() != direction {
			continue
		}
		ports = append(ports, output)
	}
	return ports
}

func Connectors(ports []Port) []*Connector {
	connectors := make([]*Connector, 0, 1)
	for i := 0; i < len(ports); i++ {
		port := ports[i]
		portConnectors := port.Connectors()
		if len(portConnectors) == 0 {
			continue
		}
		connectors = append(connectors, portConnectors...)
	}
	return connectors
}

func PortsOld(node Node, direction Direction) []Port {
	ports := make([]Port, 0, 1)
	portType := reflect.TypeOf((*Port)(nil)).Elem()
	nodeType := reflect.TypeOf(node).Elem()
	fieldsCount := nodeType.NumField()
	for i := 0; i < fieldsCount; i++ {
		field := nodeType.Field(i)
		if !field.Type.Implements(portType) {
			continue
		}
		port := reflect.ValueOf(node).Elem().Field(i).Interface().(Port)
		if port.Direction() != direction {
			continue
		}
		ports = append(ports, port)
	}
	return ports
}
