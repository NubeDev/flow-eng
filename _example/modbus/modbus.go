package main

import (
	"fmt"
	flowctrl "github.com/NubeDev/flow-eng"
	"github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {

	bac, err := bacnet.NewServer(nil)
	if err != nil {
		log.Errorln(err)
		return
	}

	p, err := bacnet.NewAI(nil)
	if err != nil {
		log.Errorln(err)
		return
	}

	b, err := bacnet.NewAO(nil)
	if err != nil {
		log.Errorln(err)
		return
	}

	graph := flowctrl.New()

	graph.AddNode(bac)
	graph.AddNodes(p, b)
	runner := flowctrl.NewSerialRunner(graph)
	// pprint.PrintJOSN(graph.GetNodes())

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
