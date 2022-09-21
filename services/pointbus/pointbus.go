package pointbus

//import (
//	"fmt"
//	"github.com/NubeDev/flow-eng/helpers"
//	"github.com/NubeDev/flow-eng/helpers/cbus"
//	"github.com/NubeDev/flow-eng/helpers/eventbus"
//	mqtt "github.com/eclipse/paho.mqtt.golang"
//	log "github.com/sirupsen/logrus"
//	"strings"
//)
//
//type PointBus struct {
//	connected bool
//}
//
//var bacnetBus cbus.Bus
//
//func New() (*PointBus, error) {
//	bacnetBus = cbus.New(100)
//	return &PointBus{}, nil
//}
//
//type Message struct {
//	UUID string
//	Msg  mqtt.Message
//}
//
//func CheckRubixIO(topic string) (isBacnet bool) { // to try and save spamming random message
//	parts := strings.Split(topic, "/")
//	if len(parts) > 0 {
//		if parts[0] == "rubixcli" {
//			return true
//		}
//	}
//	return isBacnet
//}
//
//func CheckBACnet(topic string) (isBacnet bool) { // to try and save spamming random message
//	parts := strings.Split(topic, "/")
//	if len(parts) > 0 {
//		if parts[0] == "bacnet" {
//			return true
//		}
//	}
//	return isBacnet
//}
//
//var Handler2 mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
//	fmt.Println(11111111)
//	log.Println("NEW MQTT MES", msg.Topic(), " ", string(msg.Payload()))
//
//	eventbus.GetService().RegisterTopic(eventbus.All)
//	err := eventbus.GetService().Emit(eventbus.CTX(), eventbus.All, msg)
//	if err != nil {
//		return
//	}
//
//	//if CheckBACnet(msg.Topic()) {
//	//	mes := &Message{helpers.ShortUUID("bac"), msg}
//	//	bacnetBus.Send(mes)
//	//}
//	//if CheckRubixIO(msg.Topic()) {
//	//	mes := &Message{helpers.ShortUUID("rub"), msg}
//	//	bacnetBus.Send(mes)
//	//}
//}
//
//var Handler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
//	fmt.Println(11111111)
//	log.Println("NEW MQTT MES", msg.Topic(), " ", string(msg.Payload()))
//	if CheckBACnet(msg.Topic()) {
//		mes := &Message{helpers.ShortUUID("bac"), msg}
//		bacnetBus.Send(mes)
//	}
//	if CheckRubixIO(msg.Topic()) {
//		mes := &Message{helpers.ShortUUID("rub"), msg}
//		bacnetBus.Send(mes)
//	}
//}
//
//func (inst *PointBus) PointBus() cbus.Bus {
//	return bacnetBus
//}

//func (inst *Mqtt) PingBroker() error {
//	err := inst.PublishErr("ping from flow-eng", "ping", false)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (inst *Mqtt) PublishErr(value interface{}, topic string, retain bool) error {
//	c := inst.getClient()
//	if c != nil {
//		if topic != "" {
//			v := fmt.Sprintf("%v", value)
//			err := c.Publish(topic, mqttclient.AtMostOnce, retain, v)
//			log.Infof("pointbus-publish val:%v topic:%s", v, topic)
//			if err != nil {
//				log.Error(fmt.Sprintf("pointbus-publish topic:%s err:%s", topic, err.Error()))
//				return err
//			}
//		} else {
//			log.Error(fmt.Sprintf("pointbus-publish topic can not be empty"))
//			return errors.New(fmt.Sprintf("pointbus-publish topic can not be empty"))
//		}
//	} else {
//		return errors.New(fmt.Sprintf("pointbus-client is empty"))
//	}
//
//	return nil
//}
//
//func (inst *Mqtt) Publish(value interface{}, topic string) {
//	c := inst.getClient()
//	if c != nil {
//		if topic != "" {
//			v := fmt.Sprintf("%v", value)
//			err := c.Publish(topic, mqttclient.AtMostOnce, true, v)
//			log.Infof("pointbus-publish val:%v topic:%s", v, topic)
//			if err != nil {
//				log.Error(fmt.Sprintf("pointbus-publish topic:%s err:%s", topic, err.Error()))
//			}
//		} else {
//			log.Error(fmt.Sprintf("pointbus-publish topic can not be empty"))
//		}
//	}
//}
//
//func (inst *Mqtt) Subscribe(topic string) {
//	c := inst.getClient()
//	if topic != "" {
//		err := c.Subscribe(topic, mqttclient.AtMostOnce, handle)
//		if err != nil {
//			log.Errorf(fmt.Sprintf("pointbus-subscribe topic:%s err:%s", topic, err.Error()))
//		}
//	} else {
//		log.Errorf(fmt.Sprintf("pointbus-subscribe topic can not be empty"))
//	}
//}
//
//func (inst *Mqtt) Connected() bool {
//	return inst.connected
//}
//
//func (inst *Mqtt) SetConnect(b bool) {
//	inst.connected = b
//}
//
//func (inst *Mqtt) Connect() {
//	mqttBroker := "tcp://0.0.0.0:1883"
//	_, err := mqttclient.InternalMQTT(mqttBroker)
//	if err != nil {
//		log.Errorf(fmt.Sprintf("pointbus-subscribe-connect err:%s", err.Error()))
//	}
//	client, connected := mqttclient.GetMQTT()
//	inst.connected = connected
//	inst.client = client
//}
