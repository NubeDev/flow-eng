package main

import (
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/constant"
	"github.com/NubeDev/flow-eng/nodes/protocols/applications"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet/points"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	var err error
	cont, err := constant.NewConstNum(nil)
	if err != nil {
		log.Errorln(err)
		return
	}
	store := points.New(applications.RubixIO, nil, 1, 200, 200)
	bac, err := bacnet.NewServer(nil, store)
	if err != nil {
		log.Errorln(err)
		return
	}

	ai1, err := bacnet.NewAI(nil, store)
	if err != nil {
		log.Errorln(err)
		return
	}

	ao1, err := bacnet.NewAO(nil, store)
	if err != nil {
		log.Errorln(err)
		return
	}
	err = ao1.OverrideInputValue(node.In14, 1)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = ao1.OverrideInputValue(node.In15, 0)
	if err != nil {
		fmt.Println(err)
		return
	}
	graph := flowctrl.New()

	graph.AddNode(bac)
	graph.AddNodes(ai1, ao1, cont)
	runner := flowctrl.NewSerialRunner(graph)
	pprint.PrintJOSN(graph.GetNodes())

	log.Println("Flow started")
	for {
		err := runner.Process()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		time.Sleep(3000 * time.Millisecond)
	}

}
