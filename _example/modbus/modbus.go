package main

import (
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	"github.com/NubeDev/flow-eng/node"
	"github.com/NubeDev/flow-eng/nodes/math"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	var err error
	cont, err := math.NewConst(nil)
	if err != nil {
		log.Errorln(err)
		return
	}

	bac, err := bacnet.NewServer(nil)
	if err != nil {
		log.Errorln(err)
		return
	}

	ai1, err := bacnet.NewAI(nil)
	if err != nil {
		log.Errorln(err)
		return
	}

	ao1, err := bacnet.NewAO(nil)
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
