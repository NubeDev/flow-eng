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
	// store array.ArrStore
	store cmap.MapUtil
}

func NewStore(body *node.Spec) (node.Node, error) {
	body = node.Defaults(body, dataStore, category)
	add := node.BuildInput(node.Inp, node.TypeString, nil, body.Inputs, nil)
	maxSize := node.BuildInput(node.MaxSize, node.TypeFloat, nil, body.Inputs, nil)
	clear := node.BuildInput(node.Delete, node.TypeBool, nil, body.Inputs, nil)
	inputs := node.BuildInputs(add, maxSize, clear)
	outputs := node.BuildOutputs(node.BuildOutput(node.Outp, node.TypeString, nil, body.Outputs))
	body = node.BuildNode(body, inputs, outputs, body.Settings)
	var m = map[string]interface{}{}
	return &Store{body, cmap.New(m)}, nil
}

func (inst *Store) Process() {
	updated, _ := inst.InputUpdated(node.Inp)
	add, _ := inst.ReadPinAsString(node.Inp)
	// maxSize := inst.ReadPinAsInt(node.MaxSize)
	// clear := inst.ReadPinAsInt(node.Delete)
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
	inst.WritePin(node.Outp, value.String())
}
