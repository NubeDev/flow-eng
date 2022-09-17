package main

import (
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	bacnet "github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	log "github.com/sirupsen/logrus"
)

func main() {

	bacnet, err := bacnet.NewServer(nil)
	if err != nil {
		log.Errorln(err)
		return
	}
	pprint.PrintJOSN(bacnet)

}
