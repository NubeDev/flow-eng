package topics

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
)

type Message struct {
	UUID string
	Msg  mqtt.Message
}

func IsPV(topic string) (isBacnet bool) {
	parts := strings.Split(topic, "/")
	if len(parts) > 3 {
		if parts[0] == "bacnet" {
			if parts[3] == "pv" {
				return true
			}
		}
	}
	return isBacnet
}

func IsPri(topic string) (isBacnet bool) {
	parts := strings.Split(topic, "/")
	if len(parts) > 3 {
		if parts[0] == "bacnet" {
			if parts[3] == "pri" {
				return true
			}
		}
	}
	return isBacnet
}

func CheckRubixIO(topic string) (isBacnet bool) { // to try and save spamming random message
	parts := strings.Split(topic, "/")
	if len(parts) > 0 {
		if parts[0] == "rubixio" {
			return true
		}
	}
	return isBacnet
}

func CheckBACnet(topic string) (isBacnet bool) { // to try and save spamming random message
	parts := strings.Split(topic, "/")
	if len(parts) > 0 {
		if parts[0] == "bacnet" {
			return true
		}
	}
	return isBacnet
}
