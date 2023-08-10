package pallet

import (
	"fmt"

	"github.com/NubeDev/flow-eng/connections"
	"github.com/NubeDev/flow-eng/helpers/store"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/schemas"
)

type NodeBuilderFunc func(body *node.Spec, opts ...any) (node.Node, error)

type NodeBuilderMap map[string]NodeBuilderFunc
type CategoryBuilder map[string]NodeBuilderMap

var builders CategoryBuilder
var currentSpecPallet []*node.Spec

func RegisterNodeBuilder(category string, nodeName string, builderFunc NodeBuilderFunc) {
	if builders == nil {
		builders = make(CategoryBuilder)
	}
	nodeBuilderMap, ok := builders[category]
	if !ok {
		nodeBuilderMap = make(NodeBuilderMap)
		builders[category] = nodeBuilderMap
	}

	nodeBuilderMap[nodeName] = builderFunc
}

func Builder(body *node.Spec, connIF connections.ConnectionIF, store *store.Store, opts ...interface{}) (node.Node, error) {
	body.AddStore(store)
	body.SetConnections(connIF)
	nodeBuilderMap, ok := builders[body.Info.Category]
	if !ok {
		return nil, noNodeFoundError(body.Info.Category, body.Info.Name)
	}

	nodeFunc, ok := nodeBuilderMap[body.GetName()]
	if !ok {
		return nil, noNodeFoundError(body.Info.Category, body.Info.Name)
	}
	return nodeFunc(body, store, opts)
}

// All get all the node specs, will be used for the UI to list all the nodes
func All() []*node.Spec {
	if currentSpecPallet != nil {
		return currentSpecPallet
	}
	currentSpecPallet = []*node.Spec{}
	for _, nodeMap := range builders {
		for _, nodeFunc := range nodeMap {
			n, err := nodeFunc(nil, nil)
			if n == nil || err != nil {
				continue
			}
			s := node.ConvertToSpec(n)
			currentSpecPallet = append(currentSpecPallet, s)
		}
	}
	return currentSpecPallet
}

func GetSchema(category string, name string) *schemas.Schema {
	nodeBuilderMap, ok := builders[category]
	if !ok {
		return nil
	}

	nodeFunc, ok := nodeBuilderMap[name]
	if !ok {
		return nil
	}
	n, err := nodeFunc(nil, nil)
	if n == nil || err != nil {
		return nil
	}
	s := node.ConvertToSpec(n)
	return s.GetSchema()
}

func noNodeFoundError(category string, name string) error {
	return fmt.Errorf("no nodes found with name:%s/%s", category, name)
}
