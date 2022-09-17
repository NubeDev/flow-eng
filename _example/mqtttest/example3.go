package main

import (
	"github.com/NubeDev/flow-eng/helpers/mqttbase"
	pprint "github.com/NubeDev/flow-eng/helpers/print"
	bac "github.com/NubeDev/flow-eng/nodes/protocols/bacnet"
	log "github.com/sirupsen/logrus"
)

func main() {

	m, err := mqttbase.NewMqtt()
	if err != nil {
		return
	}

	m.Connect()

	if m.Connected() {
		m.Publish("1234", "test")
	}

	bacnet, err := bac.NewBacnetBVRead(nil, m)
	if err != nil {
		log.Errorln(err)
		return
	}
	pprint.PrintJOSN(bacnet)

}
