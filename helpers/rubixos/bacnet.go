package rubixos

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var configBACnet *ConfigBACnetServer

func GetBACnetConfig() (*ConfigBACnetServer, error) {
	if configBACnet != nil {
		return configBACnet, nil
	} else {
		yamlFile, err := ioutil.ReadFile("/data/driver-bacnet/config/config.yml")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("yamlFile.Get err   #%v ", err))
		}
		err = yaml.Unmarshal(yamlFile, &configBACnet)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("yamlFile.Get err   #%v ", err))
		}
		return configBACnet, nil
	}
}

type BacnetClient struct { // this is for bacnet-master
	Debug    bool     `json:"debug" yaml:"debug"`
	Enable   bool     `json:"enable" yaml:"enable"`
	Commands []string `json:"commands" yaml:"commands"`
	Tokens   []string `json:"tokens" yaml:"tokens"`
}

type Mqtt struct {
	BrokerIp          string `json:"broker_ip"  yaml:"broker_ip"`
	BrokerPort        int    `json:"broker_port"  yaml:"broker_port"`
	Debug             bool   `json:"debug" yaml:"debug"`
	Enable            bool   `json:"enable" yaml:"enable"`
	WriteViaSubscribe bool   `json:"write_via_subscribe" yaml:"write_via_subscribe"`
	RetryEnable       bool   `json:"retry_enable" yaml:"retry_enable"`
	RetryLimit        int    `json:"retry_limit" yaml:"retry_limit"`
	RetryInterval     int    `json:"retry_interval" yaml:"retry_interval"`
}

type ConfigBACnetServer struct {
	ServerName   string       `json:"server_name" yaml:"server_name"`
	DeviceId     int          `json:"device_id" yaml:"device_id"`
	Port         int          `json:"port" yaml:"port"`
	Iface        string       `json:"iface" yaml:"iface"`
	BiMax        int          `json:"bi_max" yaml:"bi_max"`
	BoMax        int          `json:"bo_max" yaml:"bo_max"`
	BvMax        int          `json:"bv_max" yaml:"bv_max"`
	AiMax        int          `json:"ai_max" yaml:"ai_max"`
	AoMax        int          `json:"ao_max" yaml:"ao_max"`
	AvMax        int          `json:"av_max" yaml:"av_max"`
	Objects      []string     `json:"objects" yaml:"objects"`
	Properties   []string     `json:"properties" yaml:"properties"`
	Mqtt         Mqtt         `json:"mqtt" yaml:"mqtt"`
	BacnetClient BacnetClient `json:"bacnet_client" yaml:"bacnet_client"`
}
