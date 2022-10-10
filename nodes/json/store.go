package nodejson

import (
	"encoding/json"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/cmap"
	"github.com/NubeDev/flow-eng/node"
	"github.com/tidwall/gjson"
	"strconv"
	"time"
)

type Store struct {
	*node.Spec
	//store array.ArrStore
	store cmap.MapUtil
}

func NewStore(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, dataStore, category)
	add := node.BuildInput(node.In, node.TypeString, nil, body.Inputs)
	maxSize := node.BuildInput(node.MaxSize, node.TypeFloat, nil, body.Inputs)
	clear := node.BuildInput(node.Delete, node.TypeBool, nil, body.Inputs)
	inputs := node.BuildInputs(add, maxSize, clear)
	outputs := node.BuildOutputs(node.BuildOutput(node.Out, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, nil)
	var m = map[string]interface{}{}
	return &Store{body, cmap.New(m)}, nil
}

func (inst *Store) Process() {
	updated, _ := inst.InputUpdated(node.In)
	add, _ := inst.ReadPinAsString(node.In)
	//maxSize := inst.ReadPinAsInt(node.MaxSize)
	//clear := inst.ReadPinAsInt(node.Delete)
	if updated {
		if add != "null" {
			if add != "<nil>" {
				if add != "" {
					fmt.Println("add", add)
					s := strconv.FormatInt(time.Now().Unix(), 10)
					inst.store.Add(s, add)
				}
			}
		}
	}
	out, _ := json.Marshal(inst.store.GetAll())
	value := gjson.ParseBytes(out)
	inst.WritePin(node.Out, value.String())
}
